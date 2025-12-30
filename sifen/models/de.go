package models

import (
	"encoding/xml"
	"time"

	"github.com/rodascaar/sifen-go-py/sifen/types"
)

// DocumentoElectronico represents the rDE XML structure
type DocumentoElectronico struct {
	XMLName      xml.Name `xml:"rDE"`
	XmlnsXsi     string   `xml:"xmlns:xsi,attr"`
	XsiSchemaLoc string   `xml:"xsi:schemaLocation,attr"`

	DVerFor  int       `xml:"dVerFor"` // Version del Formato (150)
	DE       DE        `xml:"DE"`
	GCamFuFD *GCamFuFD `xml:"gCamFuFD,omitempty"`
}

type DE struct {
	Id        string `xml:"Id,attr"`
	DDVId     string `xml:"dDVId"`
	DFecFirma string `xml:"dFecFirma"` // Format: yyyy-mm-ddThh:mi:ss
	DSisFact  int16  `xml:"dSisFact"`  // 1=Sistema Cliente, 2=Facturacion Gratuita SET in Java (short)

	GOpeDE      TgOpeDE       `xml:"gOpeDE"`
	GTimb       TgTimb        `xml:"gTimb"`
	GDatGralOpe TdDatGralOpe  `xml:"gDatGralOpe"`
	GDtipDE     TgDtipDE      `xml:"gDtipDE"`
	GTotSub     *TgTotSub     `xml:"gTotSub,omitempty"`
	GCamGen     *TgCamGen     `xml:"gCamGen,omitempty"`
	GCamDEAsoc  []TgCamDEAsoc `xml:"gCamDEAsoc,omitempty"`
}

type TgOpeDE struct {
	ITipEmi    types.TTipEmi `xml:"iTipEmi"`
	DDesTipEmi string        `xml:"dDesTipEmi"`
	DCodSeg    string        `xml:"dCodSeg"`
	DInfoEmi   string        `xml:"dInfoEmi,omitempty"`
	DInfoFisc  string        `xml:"dInfoFisc,omitempty"`
}

type TgTimb struct {
	ITiDE     types.TTiDE `xml:"iTiDE"`
	DDesTiDE  string      `xml:"dDesTiDE"`
	DNumTim   int32       `xml:"dNumTim"`
	DEst      string      `xml:"dEst"`    // 3 chars
	DPunExp   string      `xml:"dPunExp"` // 3 chars
	DNumDoc   string      `xml:"dNumDoc"` // 7 chars
	DSerieNum string      `xml:"dSerieNum,omitempty"`
	DFeIniT   string      `xml:"dFeIniT"` // yyyy-MM-dd
}

type TdDatGralOpe struct {
	DFeEmiDE string    `xml:"dFeEmiDE"` // yyyy-MM-ddTHH:mm:ss
	GOpeCom  *TgOpeCom `xml:"gOpeCom,omitempty"`
	GEmis    TgEmis    `xml:"gEmis"`
	GDatRec  TgDatRec  `xml:"gDatRec"`
}

type GCamFuFD struct {
	DCarQR string `xml:"dCarQR"`
}

// ... TgOpeDE, TgTimb, TdDatGralOpe ...

func NewDE(id string) *DocumentoElectronico {
	return &DocumentoElectronico{
		XmlnsXsi:     "http://www.w3.org/2001/XMLSchema-instance",
		XsiSchemaLoc: "http://ekuatia.set.gov.py/sifen/xsd siRecepDE_v150.xsd",
		DVerFor:      150,
		DE: DE{
			Id:        id,
			DSisFact:  1,
			DFecFirma: time.Now().Format("2006-01-02T15:04:05"),
		},
	}
}
