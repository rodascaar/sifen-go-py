package request

import (
	"encoding/xml"
)

// ============================================================================
// Consulta RUC Request
// ============================================================================
type REnviConsRUC struct {
	XMLName  xml.Name `xml:"http://ekuatia.set.gov.py/sifen/xsd rEnviConsRUC"`
	DId      int64    `xml:"dId"`
	DRUCCons string   `xml:"dRUCCons"`
}

// ============================================================================
// Recepción DE Request (Single Document)
// ============================================================================
type REnviDe struct {
	XMLName xml.Name `xml:"http://ekuatia.set.gov.py/sifen/xsd rEnviDe"`
	DId     int64    `xml:"dId"`
	XDE     XDE      `xml:"xDE"`
}

type XDE struct {
	RawRDE []byte `xml:",innerxml"`
}

// ============================================================================
// Recepción Lote DE Request (Batch Documents)
// ============================================================================
type REnviLoteDe struct {
	XMLName xml.Name `xml:"http://ekuatia.set.gov.py/sifen/xsd rEnviLoteDe"`
	DId     int64    `xml:"dId"`
	XDEList []XDE    `xml:"xDE"`
}

// ============================================================================
// Consulta DE Request (Single Document Query)
// ============================================================================
type REnviConsDE struct {
	XMLName xml.Name `xml:"http://ekuatia.set.gov.py/sifen/xsd rEnviConsDe"`
	DId     int64    `xml:"dId"`
	DCdCDE  string   `xml:"dCdCDE"` // CDC del documento a consultar
}

// ============================================================================
// Consulta Lote DE Request (Batch Query)
// ============================================================================
type REnviConsLoteDe struct {
	XMLName       xml.Name `xml:"http://ekuatia.set.gov.py/sifen/xsd rEnviConsLoteDe"`
	DId           int64    `xml:"dId"`
	DProtConsLote string   `xml:"dProtConsLote"` // Número de lote a consultar
}

// ============================================================================
// Evento Request
// ============================================================================
type REnviEventoDe struct {
	XMLName xml.Name `xml:"http://ekuatia.set.gov.py/sifen/xsd rEnviEventoDe"`
	DId     int64    `xml:"dId"`
	DEvReg  string   `xml:"dEvReg"` // Evento registrado (XML del evento firmado)
}
