package events

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/rodascaar/sifen-go-py/sifen/types"
)

// ============================================================================
// Base Event Structure
// ============================================================================

// REvento: Envelope para eventos SIFEN
type REvento struct {
	XMLName  xml.Name `xml:"rEnviEventoDe"`
	XmlnsXsi string   `xml:"xmlns:xsi,attr"`
	DId      int64    `xml:"dId"`
	DDTE     string   `xml:"dFecFirma"`
	GEvento  GEvento  `xml:"gEvento"`
}

// GEvento: Grupo principal del evento
type GEvento struct {
	Id           string        `xml:"Id,attr"`
	DFecFirma    string        `xml:"dFecFirma"`
	DVerFor      int16         `xml:"dVerFor"`
	GGroupGesEvc *GGroupGesEve `xml:"gGroupGesEve,omitempty"` // Eventos del emisor
	GGroupTiEvt  *GGroupTiEvt  `xml:"gGroupTiEvt,omitempty"`  // Eventos del receptor
}

// GGroupGesEve: Grupo para eventos del emisor (Cancelación, Inutilización, Nominación, etc.)
type GGroupGesEve struct {
	RGesEveDE  *RGesEveDE           `xml:"rGeVeCanworking,omitempty"` // Cancelación
	RGesEveInu *EventoInutilizacion `xml:"rGeVeInu,omitempty"`        // Inutilización
	RGesEveNom *EventoNominacion    `xml:"rGeVeNom,omitempty"`        // Nominación
	RGesEveTra *EventoActTransporte `xml:"rGeVeTra,omitempty"`        // Actualización transporte
}

// GGroupTiEvt: Grupo para eventos del receptor
type GGroupTiEvt struct {
	RGeTrReceptorConf *EventoConformidad     `xml:"rGeTrReConf,omitempty"` // Conformidad
	RGeTrReceptorDisc *EventoDisconformidad  `xml:"rGeTrReDisc,omitempty"` // Disconformidad
	RGeTrReceptorDesc *EventoDesconocimiento `xml:"rGeTrReDesc,omitempty"` // Desconocimiento
	RGeTrReceptorNot  *EventoNotificacion    `xml:"rGeTrReNot,omitempty"`  // Notificación
}

// ============================================================================
// EventBuilder: Constructor de eventos base
// ============================================================================

type EventBuilder struct {
	dId       int64
	version   int16
	rucEmisor string
	dvEmisor  string
}

func NewEventBuilder(requestId int64, ruc, dv string) *EventBuilder {
	return &EventBuilder{
		dId:       requestId,
		version:   150,
		rucEmisor: ruc,
		dvEmisor:  dv,
	}
}

func (b *EventBuilder) generateEventId() string {
	return fmt.Sprintf("%s%s%d", b.rucEmisor, b.dvEmisor, time.Now().UnixNano())
}

// ============================================================================
// Evento Cancelación
// ============================================================================

// EventoCancelacion: Datos de evento de cancelación
type EventoCancelacion struct {
	CDC    string
	Motivo string
}

// RGesEveDE: Estructura de cancelación de DE
type RGesEveDE struct {
	DId      int64    `xml:"dId"`
	MEvento  string   `xml:"mEvento"`
	GEvEmiDE GEvEmiDE `xml:"gEvEmiDE"`
}

// GEvEmiDE: Grupo evento emisor - Cancelación
type GEvEmiDE struct {
	DRucEmi string     `xml:"dRucEmi"`
	DDVEmi  string     `xml:"dDVEmi"`
	ITipEvt int16      `xml:"iTipEvt"` // 1 = Cancelación
	MEvento string     `xml:"mMotEve"` // Motivo del evento
	GCamEve GCamEveCan `xml:"gGroupGesEve"`
}

type GCamEveCan struct {
	DCDC string `xml:"dCDC"` // CDC a cancelar
}

func (b *EventBuilder) BuildCancelacion(data EventoCancelacion) (*REvento, error) {
	if len(data.CDC) != 44 {
		return nil, fmt.Errorf("CDC debe tener 44 caracteres")
	}
	if data.Motivo == "" {
		return nil, fmt.Errorf("motivo es requerido")
	}

	evento := &REvento{
		XmlnsXsi: "http://www.w3.org/2001/XMLSchema-instance",
		DId:      b.dId,
		DDTE:     time.Now().Format("2006-01-02T15:04:05"),
		GEvento: GEvento{
			Id:        b.generateEventId(),
			DFecFirma: time.Now().Format("2006-01-02T15:04:05"),
			DVerFor:   b.version,
			GGroupGesEvc: &GGroupGesEve{
				RGesEveDE: &RGesEveDE{
					DId:     b.dId,
					MEvento: data.Motivo,
					GEvEmiDE: GEvEmiDE{
						DRucEmi: b.rucEmisor,
						DDVEmi:  b.dvEmisor,
						ITipEvt: 1,
						MEvento: data.Motivo,
						GCamEve: GCamEveCan{
							DCDC: data.CDC,
						},
					},
				},
			},
		},
	}

	return evento, nil
}

// ============================================================================
// Evento Inutilización
// ============================================================================

// EventoInutilizacionData: Datos para evento de inutilización
type EventoInutilizacionData struct {
	TipoDocumento   types.TTiDE
	Establecimiento string
	Punto           string
	Desde           int32
	Hasta           int32
	Motivo          string
}

// EventoInutilizacion: Estructura de inutilización
type EventoInutilizacion struct {
	DRucEmi    string `xml:"dRucEmi"`
	DDVEmi     string `xml:"dDVEmi"`
	ITipEvt    int16  `xml:"iTipEvt"`   // 2 = Inutilización
	MEventoInu string `xml:"mMotEvInu"` // Motivo
	DCoEvInu   GEvInu `xml:"gGroupGesEve"`
}

type GEvInu struct {
	ITiDE    types.TTiDE `xml:"iTiDE"`        // Tipo de documento
	DDesTiDE string      `xml:"dDesTiDE"`     // Descripción tipo doc
	DEst     string      `xml:"dEst"`         // Establecimiento
	DPunExp  string      `xml:"dPunExp"`      // Punto de expedición
	DNumDocI int32       `xml:"dNumDocDesde"` // Número desde
	DNumDocF int32       `xml:"dNumDocHasta"` // Número hasta
}

func (b *EventBuilder) BuildInutilizacion(data EventoInutilizacionData) (*REvento, error) {
	if data.Establecimiento == "" || len(data.Establecimiento) != 3 {
		return nil, fmt.Errorf("establecimiento debe tener 3 caracteres")
	}
	if data.Punto == "" || len(data.Punto) != 3 {
		return nil, fmt.Errorf("punto debe tener 3 caracteres")
	}
	if data.Desde > data.Hasta {
		return nil, fmt.Errorf("rango inválido: desde debe ser menor o igual a hasta")
	}
	if data.Motivo == "" {
		return nil, fmt.Errorf("motivo es requerido")
	}

	evento := &REvento{
		XmlnsXsi: "http://www.w3.org/2001/XMLSchema-instance",
		DId:      b.dId,
		DDTE:     time.Now().Format("2006-01-02T15:04:05"),
		GEvento: GEvento{
			Id:        b.generateEventId(),
			DFecFirma: time.Now().Format("2006-01-02T15:04:05"),
			DVerFor:   b.version,
			GGroupGesEvc: &GGroupGesEve{
				RGesEveInu: &EventoInutilizacion{
					DRucEmi:    b.rucEmisor,
					DDVEmi:     b.dvEmisor,
					ITipEvt:    2,
					MEventoInu: data.Motivo,
					DCoEvInu: GEvInu{
						ITiDE:    data.TipoDocumento,
						DDesTiDE: data.TipoDocumento.String(),
						DEst:     data.Establecimiento,
						DPunExp:  data.Punto,
						DNumDocI: data.Desde,
						DNumDocF: data.Hasta,
					},
				},
			},
		},
	}

	return evento, nil
}

// ============================================================================
// Evento Conformidad
// ============================================================================

type EventoConformidadData struct {
	CDC             string
	TipoConformidad types.TiTipoConformidad
	FechaRecepcion  time.Time
}

// EventoConformidad: Estructura de conformidad del receptor
type EventoConformidad struct {
	DRucRec     string                  `xml:"dRucRec"`
	DDVRec      string                  `xml:"dDVRec"`
	ITipEvt     int16                   `xml:"iTipEvt"` // 11 = Conformidad
	DCDC        string                  `xml:"dCDC"`
	ITiConf     types.TiTipoConformidad `xml:"iTiConf"`
	DDesTiConf  string                  `xml:"dDesTiConf"`
	DFecRecep   string                  `xml:"dFecRecep"` // Fecha recepción
	DFecEmiConf string                  `xml:"dFecEmiConf"`
}

func (b *EventBuilder) BuildConformidad(data EventoConformidadData) (*REvento, error) {
	if len(data.CDC) != 44 {
		return nil, fmt.Errorf("CDC debe tener 44 caracteres")
	}

	fechaRecep := data.FechaRecepcion
	if fechaRecep.IsZero() {
		fechaRecep = time.Now()
	}

	evento := &REvento{
		XmlnsXsi: "http://www.w3.org/2001/XMLSchema-instance",
		DId:      b.dId,
		DDTE:     time.Now().Format("2006-01-02T15:04:05"),
		GEvento: GEvento{
			Id:        b.generateEventId(),
			DFecFirma: time.Now().Format("2006-01-02T15:04:05"),
			DVerFor:   b.version,
			GGroupTiEvt: &GGroupTiEvt{
				RGeTrReceptorConf: &EventoConformidad{
					DRucRec:     b.rucEmisor,
					DDVRec:      b.dvEmisor,
					ITipEvt:     11,
					DCDC:        data.CDC,
					ITiConf:     data.TipoConformidad,
					DDesTiConf:  data.TipoConformidad.String(),
					DFecRecep:   fechaRecep.Format("2006-01-02T15:04:05"),
					DFecEmiConf: time.Now().Format("2006-01-02T15:04:05"),
				},
			},
		},
	}

	return evento, nil
}

// ============================================================================
// Evento Disconformidad
// ============================================================================

type EventoDisconformidadData struct {
	CDC    string
	Motivo string
}

// EventoDisconformidad: Estructura de disconformidad del receptor
type EventoDisconformidad struct {
	DRucRec     string `xml:"dRucRec"`
	DDVRec      string `xml:"dDVRec"`
	ITipEvt     int16  `xml:"iTipEvt"` // 12 = Disconformidad
	DCDC        string `xml:"dCDC"`
	MMotDisc    string `xml:"mMotDisc"`
	DFecEmiDisc string `xml:"dFecEmiDisc"`
}

func (b *EventBuilder) BuildDisconformidad(data EventoDisconformidadData) (*REvento, error) {
	if len(data.CDC) != 44 {
		return nil, fmt.Errorf("CDC debe tener 44 caracteres")
	}
	if data.Motivo == "" {
		return nil, fmt.Errorf("motivo es requerido")
	}

	evento := &REvento{
		XmlnsXsi: "http://www.w3.org/2001/XMLSchema-instance",
		DId:      b.dId,
		DDTE:     time.Now().Format("2006-01-02T15:04:05"),
		GEvento: GEvento{
			Id:        b.generateEventId(),
			DFecFirma: time.Now().Format("2006-01-02T15:04:05"),
			DVerFor:   b.version,
			GGroupTiEvt: &GGroupTiEvt{
				RGeTrReceptorDisc: &EventoDisconformidad{
					DRucRec:     b.rucEmisor,
					DDVRec:      b.dvEmisor,
					ITipEvt:     12,
					DCDC:        data.CDC,
					MMotDisc:    data.Motivo,
					DFecEmiDisc: time.Now().Format("2006-01-02T15:04:05"),
				},
			},
		},
	}

	return evento, nil
}

// ============================================================================
// Evento Desconocimiento
// ============================================================================

type EventoDesconocimientoData struct {
	CDC            string
	FechaEmision   time.Time
	FechaRecepcion time.Time
	TipoReceptor   types.TiNatRec
	Nombre         string
	RUC            string
	TipoDocumento  types.TTipDocRec
	NumeroDoc      string
	Motivo         string
}

// EventoDesconocimiento: Estructura de desconocimiento del DE
type EventoDesconocimiento struct {
	DRucRec     string           `xml:"dRucRec,omitempty"`
	DDVRec      string           `xml:"dDVRec,omitempty"`
	ITipEvt     int16            `xml:"iTipEvt"` // 13 = Desconocimiento
	DCDC        string           `xml:"dCDC"`
	DFecEmi     string           `xml:"dFecEmi"`
	DFecRecep   string           `xml:"dFecRecep"`
	ITiRec      types.TiNatRec   `xml:"iTiRec"`
	DNomRec     string           `xml:"dNomRec"`
	DRucRecDes  string           `xml:"dRucRecDes,omitempty"`
	ITipDocRec  types.TTipDocRec `xml:"iTipIDRec,omitempty"`
	DNumDocRec  string           `xml:"dNumIDRec,omitempty"`
	MMotDesc    string           `xml:"mMotDesc"`
	DFecEmiDesc string           `xml:"dFecEmiDesc"`
}

func (b *EventBuilder) BuildDesconocimiento(data EventoDesconocimientoData) (*REvento, error) {
	if len(data.CDC) != 44 {
		return nil, fmt.Errorf("CDC debe tener 44 caracteres")
	}
	if data.Motivo == "" {
		return nil, fmt.Errorf("motivo es requerido")
	}
	if data.Nombre == "" {
		return nil, fmt.Errorf("nombre es requerido")
	}

	fechaEmi := data.FechaEmision
	if fechaEmi.IsZero() {
		fechaEmi = time.Now()
	}
	fechaRecep := data.FechaRecepcion
	if fechaRecep.IsZero() {
		fechaRecep = time.Now()
	}

	evento := &REvento{
		XmlnsXsi: "http://www.w3.org/2001/XMLSchema-instance",
		DId:      b.dId,
		DDTE:     time.Now().Format("2006-01-02T15:04:05"),
		GEvento: GEvento{
			Id:        b.generateEventId(),
			DFecFirma: time.Now().Format("2006-01-02T15:04:05"),
			DVerFor:   b.version,
			GGroupTiEvt: &GGroupTiEvt{
				RGeTrReceptorDesc: &EventoDesconocimiento{
					DRucRec:     b.rucEmisor,
					DDVRec:      b.dvEmisor,
					ITipEvt:     13,
					DCDC:        data.CDC,
					DFecEmi:     fechaEmi.Format("2006-01-02T15:04:05"),
					DFecRecep:   fechaRecep.Format("2006-01-02T15:04:05"),
					ITiRec:      data.TipoReceptor,
					DNomRec:     data.Nombre,
					DRucRecDes:  data.RUC,
					ITipDocRec:  data.TipoDocumento,
					DNumDocRec:  data.NumeroDoc,
					MMotDesc:    data.Motivo,
					DFecEmiDesc: time.Now().Format("2006-01-02T15:04:05"),
				},
			},
		},
	}

	return evento, nil
}

// ============================================================================
// Evento Notificación
// ============================================================================

type EventoNotificacionData struct {
	CDC            string
	FechaEmision   time.Time
	FechaRecepcion time.Time
	TipoReceptor   types.TiNatRec
	Nombre         string
	RUC            string
	TipoDocumento  types.TTipDocRec
	NumeroDoc      string
	TotalPYG       float64
}

// EventoNotificacion: Estructura de notificación de recepción
type EventoNotificacion struct {
	DRucRec    string           `xml:"dRucRec,omitempty"`
	DDVRec     string           `xml:"dDVRec,omitempty"`
	ITipEvt    int16            `xml:"iTipEvt"` // 14 = Notificación
	DCDC       string           `xml:"dCDC"`
	DFecEmi    string           `xml:"dFecEmi"`
	DFecRecep  string           `xml:"dFecRecep"`
	ITiRec     types.TiNatRec   `xml:"iTiRec"`
	DNomRec    string           `xml:"dNomRec"`
	DRucRecNot string           `xml:"dRucRecNot,omitempty"`
	ITipDocRec types.TTipDocRec `xml:"iTipIDRec,omitempty"`
	DNumDocRec string           `xml:"dNumIDRec,omitempty"`
	DTotPyg    float64          `xml:"dTotPyg"`
	DFecEmiNot string           `xml:"dFecEmiNot"`
}

func (b *EventBuilder) BuildNotificacion(data EventoNotificacionData) (*REvento, error) {
	if len(data.CDC) != 44 {
		return nil, fmt.Errorf("CDC debe tener 44 caracteres")
	}
	if data.Nombre == "" {
		return nil, fmt.Errorf("nombre es requerido")
	}

	fechaEmi := data.FechaEmision
	if fechaEmi.IsZero() {
		fechaEmi = time.Now()
	}
	fechaRecep := data.FechaRecepcion
	if fechaRecep.IsZero() {
		fechaRecep = time.Now()
	}

	evento := &REvento{
		XmlnsXsi: "http://www.w3.org/2001/XMLSchema-instance",
		DId:      b.dId,
		DDTE:     time.Now().Format("2006-01-02T15:04:05"),
		GEvento: GEvento{
			Id:        b.generateEventId(),
			DFecFirma: time.Now().Format("2006-01-02T15:04:05"),
			DVerFor:   b.version,
			GGroupTiEvt: &GGroupTiEvt{
				RGeTrReceptorNot: &EventoNotificacion{
					DRucRec:    b.rucEmisor,
					DDVRec:     b.dvEmisor,
					ITipEvt:    14,
					DCDC:       data.CDC,
					DFecEmi:    fechaEmi.Format("2006-01-02T15:04:05"),
					DFecRecep:  fechaRecep.Format("2006-01-02T15:04:05"),
					ITiRec:     data.TipoReceptor,
					DNomRec:    data.Nombre,
					DRucRecNot: data.RUC,
					ITipDocRec: data.TipoDocumento,
					DNumDocRec: data.NumeroDoc,
					DTotPyg:    data.TotalPYG,
					DFecEmiNot: time.Now().Format("2006-01-02T15:04:05"),
				},
			},
		},
	}

	return evento, nil
}

// ============================================================================
// Evento Nominación
// ============================================================================

type EventoNominacionData struct {
	CDC            string
	RUCNominado    string
	DVNominado     string
	NombreNominado string
}

type EventoNominacion struct {
	DRucEmi    string `xml:"dRucEmi"`
	DDVEmi     string `xml:"dDVEmi"`
	ITipEvt    int16  `xml:"iTipEvt"` // 20 = Nominación
	DCDC       string `xml:"dCDC"`
	DRucNom    string `xml:"dRucNom"`    // RUC del nominado
	DDVNom     string `xml:"dDVNom"`     // DV del nominado
	DNomNom    string `xml:"dNomNom"`    // Nombre del nominado
	DFecEmiNom string `xml:"dFecEmiNom"` // Fecha emisión nominación
}

func (b *EventBuilder) BuildNominacion(data EventoNominacionData) (*REvento, error) {
	if len(data.CDC) != 44 {
		return nil, fmt.Errorf("CDC debe tener 44 caracteres")
	}
	if data.RUCNominado == "" {
		return nil, fmt.Errorf("RUC del nominado es requerido")
	}
	if data.NombreNominado == "" {
		return nil, fmt.Errorf("nombre del nominado es requerido")
	}

	evento := &REvento{
		XmlnsXsi: "http://www.w3.org/2001/XMLSchema-instance",
		DId:      b.dId,
		DDTE:     time.Now().Format("2006-01-02T15:04:05"),
		GEvento: GEvento{
			Id:        b.generateEventId(),
			DFecFirma: time.Now().Format("2006-01-02T15:04:05"),
			DVerFor:   b.version,
			GGroupGesEvc: &GGroupGesEve{
				RGesEveNom: &EventoNominacion{
					DRucEmi:    b.rucEmisor,
					DDVEmi:     b.dvEmisor,
					ITipEvt:    20,
					DCDC:       data.CDC,
					DRucNom:    data.RUCNominado,
					DDVNom:     data.DVNominado,
					DNomNom:    data.NombreNominado,
					DFecEmiNom: time.Now().Format("2006-01-02T15:04:05"),
				},
			},
		},
	}

	return evento, nil
}

// ============================================================================
// Evento Actualización de Datos del Transporte (rGeVeTr)
// ============================================================================

// MotivoActualizacionTransporte: Códigos de motivo para actualización de transporte
type MotivoActualizacionTransporte int16

const (
	// MotivoActTransp_CambioLocalEntrega: Cambio del local de entrega
	// Requiere: nueva dirección, departamento, distrito y ciudad
	MotivoActTransp_CambioLocalEntrega MotivoActualizacionTransporte = 1

	// MotivoActTransp_CambioChofer: Cambio del chofer
	// Requiere: nombre y número de documento del nuevo conductor
	MotivoActTransp_CambioChofer MotivoActualizacionTransporte = 2

	// MotivoActTransp_CambioTransportista: Cambio del transportista
	// Requiere: RUC (si contribuyente) o documento de identidad, y nombre
	MotivoActTransp_CambioTransportista MotivoActualizacionTransporte = 3

	// MotivoActTransp_CambioVehiculo: Cambio de vehículo
	// Requiere: tipo de vehículo, marca y número de placa o chasis
	MotivoActTransp_CambioVehiculo MotivoActualizacionTransporte = 4
)

func (m MotivoActualizacionTransporte) String() string {
	switch m {
	case MotivoActTransp_CambioLocalEntrega:
		return "Cambio del local de entrega"
	case MotivoActTransp_CambioChofer:
		return "Cambio del chofer"
	case MotivoActTransp_CambioTransportista:
		return "Cambio del transportista"
	case MotivoActTransp_CambioVehiculo:
		return "Cambio de vehículo"
	default:
		return fmt.Sprintf("Motivo %d", m)
	}
}

// EventoActTransporteData: Datos para evento de actualización de transporte
type EventoActTransporteData struct {
	CDC    string                        // CDC del documento afectado (ej: Nota de Remisión)
	Motivo MotivoActualizacionTransporte // Motivo de la actualización (1-4)

	// Datos para Cambio de Local de Entrega (Motivo 1)
	NuevaDireccion    string
	NuevoDepartamento types.TDepartamento
	NuevoDistrito     int16
	NuevaCiudad       int32
	DescDepartamento  string
	DescDistrito      string
	DescCiudad        string

	// Datos para Cambio de Chofer (Motivo 2)
	NombreChofer    string
	DocumentoChofer string

	// Datos para Cambio de Transportista (Motivo 3)
	RUCTransportista     string // Si es contribuyente
	DVTransportista      string
	DocTransportista     string // Si no es contribuyente
	TipoDocTransportista types.TTipDocRec
	NombreTransportista  string
	EsContribuyente      bool

	// Datos para Cambio de Vehículo (Motivo 4)
	TipoVehiculo  string
	MarcaVehiculo string
	PlacaVehiculo string // Número de placa o chasis
}

// EventoActTransporte: Estructura de actualización de datos del transporte
type EventoActTransporte struct {
	DRucEmi    string                        `xml:"dRucEmi"`
	DDVEmi     string                        `xml:"dDVEmi"`
	ITipEvt    int16                         `xml:"iTipEvt"`    // 21 = Actualización transporte
	DCDC       string                        `xml:"dCDC"`       // CDC del documento afectado
	DMotEv     MotivoActualizacionTransporte `xml:"dMotEv"`     // Motivo de la actualización
	DDesMotEv  string                        `xml:"dDesMotEv"`  // Descripción del motivo
	DFecEmiEvt string                        `xml:"dFecEmiEvt"` // Fecha emisión del evento

	// Grupo para cambio de local de entrega
	GCamLocEnt *GCamLocEntrega `xml:"gCamLocEnt,omitempty"`

	// Grupo para cambio de chofer
	GCamChof *GCamChofer `xml:"gCamChof,omitempty"`

	// Grupo para cambio de transportista
	GCamTrans *GCamTransportista `xml:"gCamTrans,omitempty"`

	// Grupo para cambio de vehículo
	GCamVeh *GCamVehiculo `xml:"gCamVeh,omitempty"`
}

// GCamLocEntrega: Grupo para cambio de local de entrega
type GCamLocEntrega struct {
	DDirLocEnt string              `xml:"dDirLocEnt"`           // Nueva dirección
	CDepEnt    types.TDepartamento `xml:"cDepEnt"`              // Código departamento
	DDesDepEnt string              `xml:"dDesDepEnt"`           // Descripción departamento
	CDisEnt    int16               `xml:"cDisEnt,omitempty"`    // Código distrito
	DDesDisEnt string              `xml:"dDesDisEnt,omitempty"` // Descripción distrito
	CCiuEnt    int32               `xml:"cCiuEnt"`              // Código ciudad
	DDesCiuEnt string              `xml:"dDesCiuEnt"`           // Descripción ciudad
}

// GCamChofer: Grupo para cambio de chofer
type GCamChofer struct {
	DNomChof   string `xml:"dNomChof"`   // Nombre del nuevo chofer
	DNumIDChof string `xml:"dNumIDChof"` // Número de documento del chofer
}

// GCamTransportista: Grupo para cambio de transportista
type GCamTransportista struct {
	INatTrans    types.TiNatRec   `xml:"iNatTrans"`              // 1=Contribuyente, 2=No contribuyente
	DNomTrans    string           `xml:"dNomTrans"`              // Nombre del transportista
	DRucTrans    string           `xml:"dRucTrans,omitempty"`    // RUC (si contribuyente)
	DDVTrans     string           `xml:"dDVTrans,omitempty"`     // DV del RUC
	ITipIDTrans  types.TTipDocRec `xml:"iTipIDTrans,omitempty"`  // Tipo documento (si no contribuyente)
	DDTipIDTrans string           `xml:"dDTipIDTrans,omitempty"` // Descripción tipo documento
	DNumIDTrans  string           `xml:"dNumIDTrans,omitempty"`  // Número documento
}

// GCamVehiculo: Grupo para cambio de vehículo
type GCamVehiculo struct {
	DTipVeh   string `xml:"dTipVeh"`             // Tipo de vehículo
	DMarVeh   string `xml:"dMarVeh,omitempty"`   // Marca del vehículo
	DNumPlaca string `xml:"dNumPlaca,omitempty"` // Número de placa o chasis
}

func (b *EventBuilder) BuildActualizacionTransporte(data EventoActTransporteData) (*REvento, error) {
	if len(data.CDC) != 44 {
		return nil, fmt.Errorf("CDC debe tener 44 caracteres")
	}
	if data.Motivo < 1 || data.Motivo > 4 {
		return nil, fmt.Errorf("motivo debe ser 1, 2, 3 o 4")
	}

	// Construir estructura base del evento
	evtTransp := &EventoActTransporte{
		DRucEmi:    b.rucEmisor,
		DDVEmi:     b.dvEmisor,
		ITipEvt:    21,
		DCDC:       data.CDC,
		DMotEv:     data.Motivo,
		DDesMotEv:  data.Motivo.String(),
		DFecEmiEvt: time.Now().Format("2006-01-02T15:04:05"),
	}

	// Agregar datos según el motivo
	switch data.Motivo {
	case MotivoActTransp_CambioLocalEntrega:
		if data.NuevaDireccion == "" {
			return nil, fmt.Errorf("nueva dirección es requerida para cambio de local")
		}
		evtTransp.GCamLocEnt = &GCamLocEntrega{
			DDirLocEnt: data.NuevaDireccion,
			CDepEnt:    data.NuevoDepartamento,
			DDesDepEnt: data.DescDepartamento,
			CDisEnt:    data.NuevoDistrito,
			DDesDisEnt: data.DescDistrito,
			CCiuEnt:    data.NuevaCiudad,
			DDesCiuEnt: data.DescCiudad,
		}

	case MotivoActTransp_CambioChofer:
		if data.NombreChofer == "" || data.DocumentoChofer == "" {
			return nil, fmt.Errorf("nombre y documento del chofer son requeridos")
		}
		evtTransp.GCamChof = &GCamChofer{
			DNomChof:   data.NombreChofer,
			DNumIDChof: data.DocumentoChofer,
		}

	case MotivoActTransp_CambioTransportista:
		if data.NombreTransportista == "" {
			return nil, fmt.Errorf("nombre del transportista es requerido")
		}
		camTrans := &GCamTransportista{
			DNomTrans: data.NombreTransportista,
		}
		if data.EsContribuyente {
			camTrans.INatTrans = types.TiNatRec_Contribuyente
			camTrans.DRucTrans = data.RUCTransportista
			camTrans.DDVTrans = data.DVTransportista
		} else {
			camTrans.INatTrans = types.TiNatRec_NoContribuyente
			camTrans.ITipIDTrans = data.TipoDocTransportista
			camTrans.DDTipIDTrans = data.TipoDocTransportista.String()
			camTrans.DNumIDTrans = data.DocTransportista
		}
		evtTransp.GCamTrans = camTrans

	case MotivoActTransp_CambioVehiculo:
		if data.TipoVehiculo == "" || data.PlacaVehiculo == "" {
			return nil, fmt.Errorf("tipo y placa del vehículo son requeridos")
		}
		evtTransp.GCamVeh = &GCamVehiculo{
			DTipVeh:   data.TipoVehiculo,
			DMarVeh:   data.MarcaVehiculo,
			DNumPlaca: data.PlacaVehiculo,
		}
	}

	evento := &REvento{
		XmlnsXsi: "http://www.w3.org/2001/XMLSchema-instance",
		DId:      b.dId,
		DDTE:     time.Now().Format("2006-01-02T15:04:05"),
		GEvento: GEvento{
			Id:        b.generateEventId(),
			DFecFirma: time.Now().Format("2006-01-02T15:04:05"),
			DVerFor:   b.version,
			GGroupGesEvc: &GGroupGesEve{
				RGesEveTra: evtTransp,
			},
		},
	}

	return evento, nil
}
