package sifen

import (
	"encoding/xml"
	"fmt"

	"github.com/rodascaar/sifen-go-py/internal/signature"
	"github.com/rodascaar/sifen-go-py/internal/soap"
	"github.com/rodascaar/sifen-go-py/sifen/cache"
	"github.com/rodascaar/sifen-go-py/sifen/errors"
	"github.com/rodascaar/sifen-go-py/sifen/events"
	"github.com/rodascaar/sifen-go-py/sifen/models"
	"github.com/rodascaar/sifen-go-py/sifen/request"
	"github.com/rodascaar/sifen-go-py/sifen/response"
)

// SifenClient provides methods to interact with the SIFEN API
type SifenClient struct {
	config     *SifenConfig
	soapClient *soap.Client
	requestID  int64
	cache      *cache.SifenCache
}

// NewSifenClient creates a new SIFEN client with the given configuration
func NewSifenClient(config *SifenConfig) (*SifenClient, error) {
	soapConfig := &soap.ClientConfig{
		TimeoutMs:          config.HttpConnectTimeout,
		UseClientCert:      config.UsarCertificadoCliente,
		ClientCertPath:     config.CertificadoCliente,
		ClientCertPassword: config.ContrasenaCertificadoCliente,
		UserAgent:          config.UserAgent,
	}

	sc, err := soap.NewClient(soapConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize SOAP client")
	}
	return &SifenClient{
		config:     config,
		soapClient: sc,
		requestID:  1,
		cache:      cache.NewSifenCacheWithConfig(config.CacheConfig),
	}, nil
}

// Close libera recursos del cliente
func (c *SifenClient) Close() {
	if c.cache != nil {
		c.cache.Close()
	}
}

// GetConfig returns the client configuration
func (c *SifenClient) GetConfig() *SifenConfig {
	return c.config
}

// ============================================================================
// RUC Consultation
// ============================================================================

// ConsultaRUC queries information about a RUC (tax ID)
func (c *SifenClient) ConsultaRUC(ruc string) (*response.RespuestaConsultaRUC, error) {
	// 1. Check Cache
	if resp, found := c.cache.RUC.GetRUC(ruc); found {
		return resp, nil
	}

	req := request.REnviConsRUC{
		DId:      c.nextID(),
		DRUCCons: ruc,
	}

	url := c.getURL(c.config.PathConsultaRUC)

	rawResp, err := c.soapClient.Send(url, req)
	if err != nil {
		return nil, errors.Wrap(err, "ConsultaRUC failed")
	}

	var env response.EnvelopeRefResponse
	if err := xml.Unmarshal(rawResp, &env); err != nil {
		return nil, errors.NewSifenResponseError("XML_ERROR", fmt.Sprintf("failed to unmarshal response: %v", err))
	}

	if env.Body.RResEnviConsRuc == nil {
		return nil, errors.NewSifenResponseError("EMPTY_RESPONSE", "response body is empty or invalid type")
	}

	resp := env.Body.RResEnviConsRuc

	// 2. Update Cache if successful
	if resp.DCodRes == "0500" || resp.DCodRes == "0501" {
		// Don't cache service errors
	} else {
		c.cache.RUC.SetRUC(ruc, resp)
	}

	return resp, nil
}

// ============================================================================
// Document Reception (Single)
// ============================================================================

// RecepcionDE sends a single electronic document for processing
func (c *SifenClient) RecepcionDE(de *models.DocumentoElectronico) (*response.RespuestaRecepcionDE, error) {
	// 1. Marshal DE to XML
	deBytes, err := xml.Marshal(de)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal DE")
	}

	// 2. Sign DE if configured
	signedBytes := deBytes
	if c.config.UsarCertificadoCliente {
		cert := c.soapClient.GetCertificate()
		if cert.PrivateKey != nil {
			signer := signature.NewSigner(cert)
			signedBytes, err = signer.Sign(deBytes, de.DE.Id)
			if err != nil {
				return nil, errors.NewCryptoError("failed to sign DE", err)
			}
		}
	}

	// 3. Wrap in Request
	req := request.REnviDe{
		DId: c.nextID(),
		XDE: request.XDE{
			RawRDE: signedBytes,
		},
	}

	url := c.getURL(c.config.PathRecibe)

	rawResp, err := c.soapClient.Send(url, req)
	if err != nil {
		return nil, errors.Wrap(err, "RecepcionDE failed")
	}

	var env response.EnvelopeRefResponse
	if err := xml.Unmarshal(rawResp, &env); err != nil {
		return nil, errors.NewSifenResponseError("XML_ERROR", fmt.Sprintf("failed to unmarshal response: %v", err))
	}

	if env.Body.RRetEnviDe == nil {
		return nil, errors.NewSifenResponseError("EMPTY_RESPONSE", "response body is empty or invalid type")
	}

	return env.Body.RRetEnviDe, nil
}

// ============================================================================
// Document Reception (Batch)
// ============================================================================

// RecepcionLoteDE sends multiple electronic documents for batch processing
func (c *SifenClient) RecepcionLoteDE(docs []*models.DocumentoElectronico) (*response.RespuestaRecepcionLoteDE, error) {
	if len(docs) == 0 {
		return nil, errors.ErrLoteVacio
	}
	if len(docs) > 50 {
		return nil, errors.ErrLoteExcedeMaximo
	}

	var xdeList []request.XDE
	for _, de := range docs {
		deBytes, err := xml.Marshal(de)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to marshal DE %s", de.DE.Id))
		}

		// Sign if configured
		signedBytes := deBytes
		if c.config.UsarCertificadoCliente {
			cert := c.soapClient.GetCertificate()
			if cert.PrivateKey != nil {
				signer := signature.NewSigner(cert)
				signedBytes, err = signer.Sign(deBytes, de.DE.Id)
				if err != nil {
					return nil, errors.NewCryptoError(fmt.Sprintf("failed to sign DE %s", de.DE.Id), err)
				}
			}
		}

		xdeList = append(xdeList, request.XDE{RawRDE: signedBytes})
	}

	req := request.REnviLoteDe{
		DId:     c.nextID(),
		XDEList: xdeList,
	}

	url := c.getURL(c.config.PathRecibeLote)

	rawResp, err := c.soapClient.Send(url, req)
	if err != nil {
		return nil, errors.Wrap(err, "RecepcionLoteDE failed")
	}

	var env response.EnvelopeRefResponse
	if err := xml.Unmarshal(rawResp, &env); err != nil {
		return nil, errors.NewSifenResponseError("XML_ERROR", fmt.Sprintf("failed to unmarshal response: %v", err))
	}

	if env.Body.RRetEnviLoteDe == nil {
		return nil, errors.NewSifenResponseError("EMPTY_RESPONSE", "response body is empty or invalid type")
	}

	return env.Body.RRetEnviLoteDe, nil
}

// ============================================================================
// Document Query
// ============================================================================

// ConsultaDE queries the status of a single electronic document by CDC
func (c *SifenClient) ConsultaDE(cdc string) (*response.RespuestaConsultaDE, error) {
	if len(cdc) != 44 {
		return nil, errors.ErrCDCInvalido
	}

	// 1. Check Cache
	if resp, found := c.cache.DE.GetDE(cdc); found {
		return resp, nil
	}

	req := request.REnviConsDE{
		DId:    c.nextID(),
		DCdCDE: cdc,
	}

	url := c.getURL(c.config.PathConsulta)

	rawResp, err := c.soapClient.Send(url, req)
	if err != nil {
		return nil, errors.Wrap(err, "ConsultaDE failed")
	}

	var env response.EnvelopeRefResponse
	if err := xml.Unmarshal(rawResp, &env); err != nil {
		return nil, errors.NewSifenResponseError("XML_ERROR", fmt.Sprintf("failed to unmarshal response: %v", err))
	}

	if env.Body.RResEnviConsDe == nil {
		return nil, errors.NewSifenResponseError("EMPTY_RESPONSE", "response body is empty or invalid type")
	}

	resp := env.Body.RResEnviConsDe

	// 2. Cache if final state
	if resp.DCodRes == "0260" || resp.DCodRes == "0530" { // Aprobado o Rechazado (estados finales)
		c.cache.DE.SetDE(cdc, resp)
	}

	return resp, nil
}

// ConsultaLoteDE queries the status of a batch of documents
func (c *SifenClient) ConsultaLoteDE(protocoloLote string) (*response.RespuestaConsultaLoteDE, error) {
	req := request.REnviConsLoteDe{
		DId:           c.nextID(),
		DProtConsLote: protocoloLote,
	}

	url := c.getURL(c.config.PathConsultaLote)

	rawResp, err := c.soapClient.Send(url, req)
	if err != nil {
		return nil, err
	}

	var env response.EnvelopeRefResponse
	if err := xml.Unmarshal(rawResp, &env); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if env.Body.RResEnviConsLoteDe == nil {
		return nil, fmt.Errorf("response body is empty or invalid type")
	}

	return env.Body.RResEnviConsLoteDe, nil
}

// ============================================================================
// Event Submission
// ============================================================================

// EnviarEvento sends an event to SIFEN
func (c *SifenClient) EnviarEvento(evento *events.REvento) (*response.RespuestaEvento, error) {
	// Marshal event to XML
	eventoBytes, err := xml.Marshal(evento)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event: %w", err)
	}

	// Sign event if configured
	signedBytes := eventoBytes
	if c.config.UsarCertificadoCliente {
		cert := c.soapClient.GetCertificate()
		if cert.PrivateKey != nil {
			signer := signature.NewSigner(cert)
			signedBytes, err = signer.Sign(eventoBytes, evento.GEvento.Id)
			if err != nil {
				return nil, fmt.Errorf("failed to sign event: %w", err)
			}
		}
	}

	req := request.REnviEventoDe{
		DId:    c.nextID(),
		DEvReg: string(signedBytes),
	}

	url := c.getURL(c.config.PathEvento)

	rawResp, err := c.soapClient.Send(url, req)
	if err != nil {
		return nil, err
	}

	var env response.EnvelopeRefResponse
	if err := xml.Unmarshal(rawResp, &env); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if env.Body.RRetEnviEventoDe == nil {
		return nil, fmt.Errorf("response body is empty or invalid type")
	}

	return env.Body.RRetEnviEventoDe, nil
}

// ============================================================================
// Event Convenience Methods
// ============================================================================

// CancelarDE cancels an electronic document
func (c *SifenClient) CancelarDE(cdc, motivo string) (*response.RespuestaEvento, error) {
	ruc, dv := c.splitRUC()
	builder := events.NewEventBuilder(c.nextID(), ruc, dv)

	evento, err := builder.BuildCancelacion(events.EventoCancelacion{
		CDC:    cdc,
		Motivo: motivo,
	})
	if err != nil {
		return nil, err
	}

	return c.EnviarEvento(evento)
}

// InutilizarNumeracion invalidates a range of document numbers
func (c *SifenClient) InutilizarNumeracion(data events.EventoInutilizacionData) (*response.RespuestaEvento, error) {
	ruc, dv := c.splitRUC()
	builder := events.NewEventBuilder(c.nextID(), ruc, dv)

	evento, err := builder.BuildInutilizacion(data)
	if err != nil {
		return nil, err
	}

	return c.EnviarEvento(evento)
}

// ConfirmarRecepcion confirms receipt of a document
func (c *SifenClient) ConfirmarRecepcion(data events.EventoConformidadData) (*response.RespuestaEvento, error) {
	ruc, dv := c.splitRUC()
	builder := events.NewEventBuilder(c.nextID(), ruc, dv)

	evento, err := builder.BuildConformidad(data)
	if err != nil {
		return nil, err
	}

	return c.EnviarEvento(evento)
}

// ReportarDisconformidad reports non-conformity of a document
func (c *SifenClient) ReportarDisconformidad(cdc, motivo string) (*response.RespuestaEvento, error) {
	ruc, dv := c.splitRUC()
	builder := events.NewEventBuilder(c.nextID(), ruc, dv)

	evento, err := builder.BuildDisconformidad(events.EventoDisconformidadData{
		CDC:    cdc,
		Motivo: motivo,
	})
	if err != nil {
		return nil, err
	}

	return c.EnviarEvento(evento)
}

// ============================================================================
// Helper Methods
// ============================================================================

func (c *SifenClient) nextID() int64 {
	c.requestID++
	return c.requestID
}

func (c *SifenClient) getURL(path string) string {
	if c.config.UrlBase != "" {
		return c.config.UrlBase + path
	}
	return c.config.UrlBaseLocal + path
}

func (c *SifenClient) splitRUC() (ruc, dv string) {
	// Extract RUC and DV from config if available
	// This is a simplified version - in production would need proper extraction
	// Assuming RUC is stored in CSC or needs to be passed via config
	return c.config.IdCSC, "0" // Placeholder - should be properly configured
}
