package sifen

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"

	"github.com/rodascaar/sifen-go-py/internal/signature"
	"github.com/rodascaar/sifen-go-py/sifen/models"
	"github.com/rodascaar/sifen-go-py/sifen/request"
	"github.com/rodascaar/sifen-go-py/sifen/response"
	"github.com/rodascaar/sifen-go-py/sifen/types"
)

// ============================================================================
// Constantes de Lotes
// ============================================================================
const (
	MaxDocumentosLote     = 50
	MaxSizeLoteKB         = 10000 // 10,000 KB para envío de lote
	MaxSizeConsultaLoteKB = 1000  // 1,000 KB para consulta de resultado
)

// ============================================================================
// Estructura XML del Lote según ProtProcesLoteDE_v150.xsd
// ============================================================================

// RLoteDE representa la raíz del lote de documentos electrónicos
type RLoteDE struct {
	XMLName xml.Name `xml:"rLoteDE"`
	DVerFor int16    `xml:"dVerFor"` // Versión del formato (150)
	// Lista de documentos electrónicos firmados
	RDEList []RDEWrapper `xml:"rDE"`
}

// RDEWrapper envuelve cada DE firmado dentro del lote
type RDEWrapper struct {
	XMLName xml.Name `xml:"rDE"`
	// Contenido del DE firmado (se inserta como raw XML)
	InnerXML string `xml:",innerxml"`
}

// ============================================================================
// Parámetros de Lote
// ============================================================================

// LoteParams contiene los parámetros para crear un lote
type LoteParams struct {
	Documentos    []*models.DocumentoElectronico
	TipoDocumento types.TTiDE // Todos deben ser del mismo tipo
}

// LoteResult contiene el resultado del envío de lote
type LoteResult struct {
	NumeroLote      string // Número/Ticket del lote para consulta posterior
	TiempoEstimado  int32  // Tiempo estimado de procesamiento en minutos
	TamanoEnviadoKB int    // Tamaño del paquete enviado
	NumDocumentos   int    // Cantidad de documentos en el lote
}

// ============================================================================
// Creación y Envío de Lotes
// ============================================================================

// CrearLoteDE construye, comprime y codifica un lote de documentos
// Retorna el contenido listo para enviar (Base64 del .zip)
func (c *SifenClient) CrearLoteDE(params LoteParams) (string, error) {
	// 1. Validaciones
	if len(params.Documentos) == 0 {
		return "", fmt.Errorf("se requiere al menos un documento")
	}
	if len(params.Documentos) > MaxDocumentosLote {
		return "", fmt.Errorf("máximo %d documentos por lote, recibidos: %d",
			MaxDocumentosLote, len(params.Documentos))
	}

	// Validar que todos los documentos sean del mismo tipo
	for i, doc := range params.Documentos {
		if doc.DE.GTimb.ITiDE != params.TipoDocumento {
			return "", fmt.Errorf("documento %d tiene tipo %d, esperado %d (todos deben ser del mismo tipo)",
				i, doc.DE.GTimb.ITiDE, params.TipoDocumento)
		}
	}

	// 2. Firmar cada documento individualmente
	var rdeList []RDEWrapper
	for _, de := range params.Documentos {
		deBytes, err := xml.Marshal(de)
		if err != nil {
			return "", fmt.Errorf("error al serializar DE %s: %w", de.DE.Id, err)
		}

		// Firmar si está configurado
		signedBytes := deBytes
		if c.config.UsarCertificadoCliente {
			cert := c.soapClient.GetCertificate()
			if cert.PrivateKey != nil {
				signer := signature.NewSigner(cert)
				signedBytes, err = signer.Sign(deBytes, de.DE.Id)
				if err != nil {
					return "", fmt.Errorf("error al firmar DE %s: %w", de.DE.Id, err)
				}
			}
		}

		rdeList = append(rdeList, RDEWrapper{
			InnerXML: string(signedBytes),
		})
	}

	// 3. Construir estructura del lote según ProtProcesLoteDE_v150.xsd
	lote := RLoteDE{
		DVerFor: 150,
		RDEList: rdeList,
	}

	loteXML, err := xml.Marshal(lote)
	if err != nil {
		return "", fmt.Errorf("error al serializar lote: %w", err)
	}

	// Agregar declaración XML
	loteXMLFull := append([]byte(`<?xml version="1.0" encoding="UTF-8"?>`), loteXML...)

	// 4. Comprimir en formato .zip
	zipBuffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(zipBuffer)

	// Crear archivo dentro del zip
	fileWriter, err := zipWriter.Create("lote.xml")
	if err != nil {
		return "", fmt.Errorf("error al crear archivo en zip: %w", err)
	}

	_, err = fileWriter.Write(loteXMLFull)
	if err != nil {
		return "", fmt.Errorf("error al escribir en zip: %w", err)
	}

	err = zipWriter.Close()
	if err != nil {
		return "", fmt.Errorf("error al cerrar zip: %w", err)
	}

	// Validar tamaño
	sizeKB := zipBuffer.Len() / 1024
	if sizeKB > MaxSizeLoteKB {
		return "", fmt.Errorf("tamaño del lote %d KB excede máximo permitido %d KB",
			sizeKB, MaxSizeLoteKB)
	}

	// 5. Codificar en Base64
	base64Content := base64.StdEncoding.EncodeToString(zipBuffer.Bytes())

	return base64Content, nil
}

// EnviarLoteDE envía un lote de documentos para procesamiento asíncrono
func (c *SifenClient) EnviarLoteDE(params LoteParams) (*LoteResult, error) {
	// 1. Crear el lote comprimido y codificado
	base64Content, err := c.CrearLoteDE(params)
	if err != nil {
		return nil, err
	}

	// 2. Construir request según especificación SIFEN
	req := request.REnviLoteDe{
		DId:     c.nextID(),
		XDEList: []request.XDE{{RawRDE: []byte(base64Content)}},
	}

	url := c.getURL(c.config.PathRecibeLote)

	rawResp, err := c.soapClient.Send(url, req)
	if err != nil {
		return nil, err
	}

	var env response.EnvelopeRefResponse
	if err := xml.Unmarshal(rawResp, &env); err != nil {
		return nil, fmt.Errorf("error al deserializar respuesta: %w", err)
	}

	if env.Body.RRetEnviLoteDe == nil {
		return nil, fmt.Errorf("respuesta vacía o inválida")
	}

	resp := env.Body.RRetEnviLoteDe

	return &LoteResult{
		NumeroLote:      resp.DProtConsLot,
		TiempoEstimado:  resp.DTmpLot,
		TamanoEnviadoKB: len(base64Content) / 1024,
		NumDocumentos:   len(params.Documentos),
	}, nil
}

// ConsultarResultadoLote consulta el estado de un lote enviado previamente
func (c *SifenClient) ConsultarResultadoLote(numeroLote string) (*response.RespuestaConsultaLoteDE, error) {
	if numeroLote == "" {
		return nil, fmt.Errorf("número de lote es requerido")
	}

	req := request.REnviConsLoteDe{
		DId:           c.nextID(),
		DProtConsLote: numeroLote,
	}

	url := c.getURL(c.config.PathConsultaLote)

	rawResp, err := c.soapClient.Send(url, req)
	if err != nil {
		return nil, err
	}

	var env response.EnvelopeRefResponse
	if err := xml.Unmarshal(rawResp, &env); err != nil {
		return nil, fmt.Errorf("error al deserializar respuesta: %w", err)
	}

	if env.Body.RResEnviConsLoteDe == nil {
		return nil, fmt.Errorf("respuesta vacía o inválida")
	}

	return env.Body.RResEnviConsLoteDe, nil
}
