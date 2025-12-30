package response

import (
	"encoding/xml"
	"time"
)

// ============================================================================
// Base Response
// ============================================================================
type BaseResponse struct {
	DCodRes string `xml:"dCodRes"` // Código de respuesta
	DMsgRes string `xml:"dMsgRes"` // Mensaje de respuesta
}

// IsSuccess returns true if the response code indicates success
func (r BaseResponse) IsSuccess() bool {
	return r.DCodRes == "0260" || r.DCodRes == "0300" || r.DCodRes == "0261"
}

// ============================================================================
// Consulta RUC Response
// ============================================================================
type RespuestaConsultaRUC struct {
	BaseResponse
	XContRUC TxContRuc `xml:"xContRUC"`
}

type TxContRuc struct {
	DRUCCons     string `xml:"dRUCCons"`     // RUC consultado
	DRazCons     string `xml:"dRazCons"`     // Razón social
	DCodEstCons  string `xml:"dCodEstCons"`  // Código de estado
	DDesEstCons  string `xml:"dDesEstCons"`  // Descripción de estado
	DRUCFactElec string `xml:"dRUCFactElec"` // RUC de facturación electrónica
}

// ============================================================================
// Recepción DE Response (Single Document)
// ============================================================================
type RespuestaRecepcionDE struct {
	BaseResponse
	RProtDe TxProtDe `xml:"rProtDe"`
}

type TxProtDe struct {
	Id       string      `xml:"Id"`       // CDC del documento
	DFecProc time.Time   `xml:"dFecProc"` // Fecha de procesamiento
	DDigVal  string      `xml:"dDigVal"`  // Digest del documento
	DEstRes  string      `xml:"dEstRes"`  // Estado de resultado
	DProtAut string      `xml:"dProtAut"` // Número de protocolo de autorización
	GResProc []TgResProc `xml:"gResProc"` // Resultados del procesamiento
}

type TgResProc struct {
	DCodRes string `xml:"dCodRes"` // Código de resultado
	DMsgRes string `xml:"dMsgRes"` // Mensaje de resultado
}

// IsApproved returns true if the document was approved
func (r RespuestaRecepcionDE) IsApproved() bool {
	return r.RProtDe.DEstRes == "Aprobado"
}

// ============================================================================
// Recepción Lote DE Response (Batch)
// ============================================================================
type RespuestaRecepcionLoteDE struct {
	BaseResponse
	DProtConsLot string `xml:"dProtConsLot"` // Número de lote para consulta posterior
	DTmpLot      int32  `xml:"dTmpLot"`      // Tiempo estimado de procesamiento (minutos)
}

// ============================================================================
// Consulta DE Response
// ============================================================================
type RespuestaConsultaDE struct {
	BaseResponse
	RProtDe *TxProtDe `xml:"rProtDe,omitempty"`   // Protocolo del DE
	RDE     []byte    `xml:"xContenDE,omitempty"` // Contenido del DE (XML)
}

// ============================================================================
// Consulta Lote DE Response
// ============================================================================
type RespuestaConsultaLoteDE struct {
	BaseResponse
	DEstLote     string          `xml:"dEstLote"`     // Estado del lote
	DProtConsLot string          `xml:"dProtConsLot"` // Número de lote
	GResProcLote []TgResProcLote `xml:"gResProcLot"`  // Resultados por documento
}

type TgResProcLote struct {
	Id       string      `xml:"Id"`       // CDC del documento
	DEstRes  string      `xml:"dEstRes"`  // Estado del resultado
	DProtAut string      `xml:"dProtAut"` // Protocolo de autorización
	GResProc []TgResProc `xml:"gResProc"` // Detalles del procesamiento
}

// ============================================================================
// Evento Response
// ============================================================================
type RespuestaEvento struct {
	BaseResponse
	RProtEve TxProtEve `xml:"rProtEve"`
}

type TxProtEve struct {
	Id       string    `xml:"Id"`       // ID del evento
	DFecProc time.Time `xml:"dFecProc"` // Fecha de procesamiento
	DCodRes  string    `xml:"dCodRes"`  // Código de resultado
	DMsgRes  string    `xml:"dMsgRes"`  // Mensaje de resultado
}

// IsApproved returns true if the event was approved
func (r RespuestaEvento) IsApproved() bool {
	return r.RProtEve.DCodRes == "0510" || r.RProtEve.DCodRes == "0520"
}

// ============================================================================
// SOAP Envelope Wrappers
// ============================================================================
type EnvelopeRefResponse struct {
	XMLName xml.Name        `xml:"Envelope"`
	Body    BodyRefResponse `xml:"Body"`
}

type BodyRefResponse struct {
	RResEnviConsRuc    *RespuestaConsultaRUC     `xml:"rResEnviConsRuc,omitempty"`
	RRetEnviDe         *RespuestaRecepcionDE     `xml:"rRetEnviDe,omitempty"`
	RRetEnviLoteDe     *RespuestaRecepcionLoteDE `xml:"rRetEnviLoteDe,omitempty"`
	RResEnviConsDe     *RespuestaConsultaDE      `xml:"rResEnviConsDe,omitempty"`
	RResEnviConsLoteDe *RespuestaConsultaLoteDE  `xml:"rResEnviConsLoteDe,omitempty"`
	RRetEnviEventoDe   *RespuestaEvento          `xml:"rRetEnviEventoDe,omitempty"`
}

// ============================================================================
// Error Codes
// ============================================================================
const (
	// Success codes
	CodeSuccess                = "0260" // DE recibido correctamente
	CodeSuccessProcessed       = "0300" // DE procesado
	CodeSuccessWithObservation = "0261" // DE aprobado con observaciones

	// Error codes - General
	CodeInternalError    = "0100" // Error interno del servidor
	CodeInvalidXML       = "0101" // XML inválido
	CodeInvalidSignature = "0102" // Firma inválida
	CodeInvalidCDC       = "0103" // CDC inválido

	// Error codes - Document
	CodeDuplicateCDC    = "0160" // CDC duplicado
	CodeInvalidIssuer   = "0161" // Emisor no autorizado
	CodeInvalidReceiver = "0162" // Receptor inválido

	// Event codes
	CodeEventSuccess  = "0510" // Evento procesado correctamente
	CodeEventAccepted = "0520" // Evento aceptado
)
