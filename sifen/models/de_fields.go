package models

import (
	"github.com/rodascaar/sifen-go-py/sifen/types"
)

// ============================================================================
// TgDtipDE: Tipo de Documento Electronico (Detalle)
// ============================================================================
type TgDtipDE struct {
	GCamFE       *TgCamFE    `xml:"gCamFE,omitempty"`   // Campos de Factura Electrónica
	GCamAE       *TgCamAE    `xml:"gCamAE,omitempty"`   // Campos de Autofactura
	GCamNCDE     *TgCamNCDE  `xml:"gCamNCDE,omitempty"` // Campos de Nota Crédito/Débito
	GCamNRE      *TgCamNRE   `xml:"gCamNRE,omitempty"`  // Campos de Nota de Remisión
	GCamCond     *TgCamCond  `xml:"gCamCond,omitempty"` // Condición de la operación
	GCamItemList []TgCamItem `xml:"gCamItem"`           // Items de la operación
	GCamEsp      *TgCamEsp   `xml:"gCamEsp,omitempty"`  // Campos por sector específico
	GTransp      *TgTransp   `xml:"gTransp,omitempty"`  // Campos de transporte
}

// ============================================================================
// TgCamFE: Campos de Factura Electrónica (E010-E099)
// ============================================================================
type TgCamFE struct {
	IIndPres    types.TiIndPres `xml:"iIndPres"`           // Indicador de presencia
	DDesIndPres string          `xml:"dDesIndPres"`        // Descripción del indicador de presencia
	DFecEmNR    string          `xml:"dFecEmNR,omitempty"` // Fecha de emisión de NR (yyyy-MM-dd)
	// Campos DNCP (Dirección Nacional de Contrataciones Públicas)
	GCompPub *TgCompPub `xml:"gCompPub,omitempty"` // Datos de compras públicas
}

// TgCompPub: Datos de Compras Públicas (E020-E029)
type TgCompPub struct {
	DModCont   string `xml:"dModCont"`             // Modalidad de contratación
	DEntCont   int32  `xml:"dEntCont"`             // Entidad contratante
	DAnoContP  int16  `xml:"dAnoContP"`            // Año del contrato
	DSecCont   int32  `xml:"dSecCont"`             // Secuencia del contrato
	DFeCodCont string `xml:"dFeCodCont,omitempty"` // Fecha del contrato (yyyy-MM-dd)
}

// ============================================================================
// TgCamAE: Campos de Autofactura Electrónica (E300-E399)
// ============================================================================
type TgCamAE struct {
	INatVen      types.TiNatVendedorAF `xml:"iNatVen"`              // Naturaleza del vendedor
	DDesNatVen   string                `xml:"dDesNatVen"`           // Descripción naturaleza vendedor
	ITipIDVen    types.TTipDocRec      `xml:"iTipIDVen"`            // Tipo de documento del vendedor
	DDesTipIDVen string                `xml:"dDesTipIDVen"`         // Descripción tipo documento
	DNumIDVen    string                `xml:"dNumIDVen"`            // Número de documento del vendedor
	DNomVen      string                `xml:"dNomVen"`              // Nombre del vendedor
	DDirVen      string                `xml:"dDirVen"`              // Dirección del vendedor
	DNumCasVen   int32                 `xml:"dNumCasVen"`           // Número de casa
	CDepVen      types.TDepartamento   `xml:"cDepVen"`              // Código departamento vendedor
	DDesDepVen   string                `xml:"dDesDepVen"`           // Descripción departamento
	CDisVen      int16                 `xml:"cDisVen,omitempty"`    // Código distrito vendedor
	DDesDisVen   string                `xml:"dDesDisVen,omitempty"` // Descripción distrito
	CCiuVen      int32                 `xml:"cCiuVen"`              // Código ciudad vendedor
	DDesCiuVen   string                `xml:"dDesCiuVen"`           // Descripción ciudad
	DLugarTrans  string                `xml:"dDirProv,omitempty"`   // Lugar de la transacción
	// Campos lugar de transacción
	GInfLugTran *TgInfLugTran `xml:"gLugRec,omitempty"` // Lugar de recepción
}

// TgInfLugTran: Lugar de la Transacción para Autofactura
type TgInfLugTran struct {
	DDirLug    string              `xml:"dDirLug"`              // Dirección del lugar
	CDepLug    types.TDepartamento `xml:"cDepLug"`              // Código departamento
	DDesDepLug string              `xml:"dDesDepLug"`           // Descripción departamento
	CDisLug    int16               `xml:"cDisLug,omitempty"`    // Código distrito
	DDesDisLug string              `xml:"dDesDisLug,omitempty"` // Descripción distrito
	CCiuLug    int32               `xml:"cCiuLug"`              // Código ciudad
	DDesCiuLug string              `xml:"dDesCiuLug"`           // Descripción ciudad
}

// ============================================================================
// TgCamNCDE: Campos de Nota de Crédito/Débito Electrónica (E400-E499)
// ============================================================================
type TgCamNCDE struct {
	IMotEmi    types.TiMotEmiNC `xml:"iMotEmi"`    // Motivo de emisión
	DDesMotEmi string           `xml:"dDesMotEmi"` // Descripción del motivo
}

// ============================================================================
// TgCamNRE: Campos de Nota de Remisión Electrónica (E500-E599)
// ============================================================================
type TgCamNRE struct {
	IMotEmiNR     types.TiMotEmiNR  `xml:"iMotEmiNR"`        // Motivo de la emisión
	DDesMotEmiNR  string            `xml:"dDesMotEmiNR"`     // Descripción del motivo
	IRespEmiNR    types.TiRespFlete `xml:"iRespEmiNR"`       // Responsable de la emisión
	DDesRespEmiNR string            `xml:"dDesRespEmiNR"`    // Descripción del responsable
	DKmR          float64           `xml:"dKmR,omitempty"`   // Kilómetros estimados de recorrido
	DFecEm        string            `xml:"dFecEm,omitempty"` // Fecha estimada de inicio de traslado
}

// ============================================================================
// TgCamCond: Condición de la Operación (E600-E699)
// ============================================================================
type TgCamCond struct {
	ICondOpe    types.TiCondOpe `xml:"iCondOpe"`             // Condición de la operación
	DDesCondOpe string          `xml:"dDCondOpe"`            // Descripción de la condición
	GPaConEIni  []TgPaConEIni   `xml:"gPaConEIni,omitempty"` // Entregas contado/efectivo
	GCredCond   *TgCredCond     `xml:"gPagCred,omitempty"`   // Condiciones de crédito
}

// TgPaConEIni: Pago Contado - Entregas (E610-E619)
type TgPaConEIni struct {
	ITiPago     types.TiTipPago `xml:"iTiPago"`               // Tipo de pago
	DDesTiPago  string          `xml:"dDesTiPag"`             // Descripción del tipo de pago
	DMonTiPag   float64         `xml:"dMonTiPag"`             // Monto del pago
	CMoneOpe    types.CMondT    `xml:"cMoneTiPag"`            // Moneda del pago
	DDesMoneOpe string          `xml:"dDMoneTiPag,omitempty"` // Descripción moneda
	DTiCamTiPag *float64        `xml:"dTiCamTiPag,omitempty"` // Tipo de cambio por pago
	// Campos para tarjeta
	GTarjeta *TgTarjeta `xml:"gPagTarCD,omitempty"` // Datos de tarjeta
	// Campos para cheque
	GCheque *TgCheque `xml:"gPagCheworking,omitempty"` // Datos de cheque
}

// TgTarjeta: Datos de Pago con Tarjeta (E620-E629)
type TgTarjeta struct {
	IDenTarj    int16  `xml:"iDenTarj"`              // Denominación de la tarjeta
	DDesDenTarj string `xml:"dDesDenTarj,omitempty"` // Descripción denominación
	DRSProTar   string `xml:"dRSProTar,omitempty"`   // Razón social de procesadora
	DRUCProTar  string `xml:"dRUCProTar,omitempty"`  // RUC de procesadora
	DDVProTar   int16  `xml:"dDVProTar,omitempty"`   // Dígito verificador RUC procesadora
	IForProPa   int16  `xml:"iForProPa,omitempty"`   // Forma de procesamiento del pago
	DCoAu662    string `xml:"dCodAu662,omitempty"`   // Código de autorización
}

// TgCheque: Datos de Pago con Cheque (E630-E639)
type TgCheque struct {
	DNumCheq    string `xml:"dNumCheq"` // Número de cheque
	DBanEmiCheq string `xml:"dBcoEmi"`  // Banco emisor del cheque
}

// TgCredCond: Condiciones de Crédito (E640-E669)
type TgCredCond struct {
	ICondCred    types.TiCondCredito `xml:"iCondCred"`           // Condición del crédito
	DDesCondCred string              `xml:"dDCondCred"`          // Descripción condición
	DPlazoCre    string              `xml:"dPlazoCre,omitempty"` // Plazo del crédito
	DCuotas      int16               `xml:"dCuotas,omitempty"`   // Cantidad de cuotas
	DMonEnt      float64             `xml:"dMonEnt,omitempty"`   // Monto de la entrega inicial
	GCuotas      []TgCuotas          `xml:"gCuotas,omitempty"`   // Detalle de cuotas
}

// TgCuotas: Detalle de Cuotas (E650-E659)
type TgCuotas struct {
	CMoneOpe    types.CMondT `xml:"cMoneCuo"`            // Moneda de la cuota
	DDesMoneCuo string       `xml:"dDMoneCuo,omitempty"` // Descripción moneda
	DMonCuota   float64      `xml:"dMonCuota"`           // Monto de la cuota
	DVencCuo    string       `xml:"dVencCuo,omitempty"`  // Fecha de vencimiento (yyyy-MM-dd)
}

// ============================================================================
// TgCamEsp: Campos por Sector Específico (E800-E899)
// ============================================================================
type TgCamEsp struct {
	GGrupEner *TgGrupEner `xml:"gGrupEner,omitempty"` // Sector energía eléctrica
	GGrupSeg  *TgGrupSeg  `xml:"gGrupSeg,omitempty"`  // Sector seguros
	GGrupSup  *TgGrupSup  `xml:"gGrupSup,omitempty"`  // Sector supermercados
	GGrupAdi  *TgGrupAdi  `xml:"gGrupAdi,omitempty"`  // Grupo de datos adicionales
}

// TgGrupEner: Sector Energía Eléctrica (E810-E819)
type TgGrupEner struct {
	DNroMed  string  `xml:"dNroMed"`            // Número de medidor
	DActEner int32   `xml:"dActEner,omitempty"` // Código de actividad
	DCatEner string  `xml:"dCatEner,omitempty"` // Categoría del servicio
	DLecAnt  float64 `xml:"dLecAnt,omitempty"`  // Lectura anterior
	DLecAct  float64 `xml:"dLecAct,omitempty"`  // Lectura actual
	DConKwh  float64 `xml:"dConKwh,omitempty"`  // Consumo en kWh
}

// TgGrupSeg: Sector Seguros (E820-E829)
type TgGrupSeg struct {
	DCodEmpSeg string    `xml:"dCodEmpSeg,omitempty"`  // Código de la aseguradora
	GPoliza    *TgPoliza `xml:"gGrupPolSeg,omitempty"` // Datos de la póliza
}

// TgPoliza: Datos Póliza de Seguro
type TgPoliza struct {
	DCodPolSeg string `xml:"dPoliza"`              // Código interno de la póliza
	DNumPolSeg string `xml:"dNumPoliza"`           // Número de póliza
	DVigencia  int16  `xml:"dVigencia,omitempty"`  // Vigencia de la póliza
	DUnidVig   string `xml:"dUnidVig,omitempty"`   // Unidad de medida de vigencia
	DFecIniVig string `xml:"dFecIniVig,omitempty"` // Fecha inicio vigencia (yyyy-MM-dd)
	DFecFinVig string `xml:"dFecFinVig,omitempty"` // Fecha fin vigencia (yyyy-MM-dd)
	DCodIntIt  string `xml:"dCodInt,omitempty"`    // Código interno del item asegurado
}

// TgGrupSup: Sector Supermercados (E830-E839)
type TgGrupSup struct {
	DNomCaj   string  `xml:"dNomCaj,omitempty"`   // Nombre del cajero
	DEfecivo  float64 `xml:"dEfectivo,omitempty"` // Monto efectivo recibido
	DVuelto   float64 `xml:"dVuelto,omitempty"`   // Monto del vuelto
	DDonac    float64 `xml:"dDonac,omitempty"`    // Monto de donación
	DDesDonac string  `xml:"dDesDonac,omitempty"` // Descripción de la donación
}

// TgGrupAdi: Grupo de Datos Adicionales (E840-E899)
type TgGrupAdi struct {
	DCiclo    string  `xml:"dCiclo,omitempty"`    // Ciclo facturado
	DFecIniC  string  `xml:"dFecIniC,omitempty"`  // Fecha inicio del ciclo
	DFecFinC  string  `xml:"dFecFinC,omitempty"`  // Fecha fin del ciclo
	DVencPag  string  `xml:"dVencPag,omitempty"`  // Fecha de vencimiento para pago
	DContrato string  `xml:"dContrato,omitempty"` // Número de contrato
	DSalAnt   float64 `xml:"dSalAnt,omitempty"`   // Saldo anterior
}

// ============================================================================
// TgTransp: Campos de Transporte (E900-E999)
// ============================================================================
type TgTransp struct {
	ITipTrans     types.TiTipoTransporte      `xml:"iTipTrans"`               // Tipo de transporte
	DDesTipTrans  string                      `xml:"dDesTipTrans"`            // Descripción tipo transporte
	IModTrans     types.TiModalidadTransporte `xml:"iModTrans"`               // Modalidad de transporte
	DDesModTrans  string                      `xml:"dDesModTrans"`            // Descripción modalidad
	IRepFlete     types.TiRespFlete           `xml:"iTipRep"`                 // Responsable del flete
	DDesRepFlete  string                      `xml:"dDesTipRep"`              // Descripción responsable
	DCodNegoci    string                      `xml:"dCondNeg,omitempty"`      // Condición de negociación
	DNuManworking string                      `xml:"dNuManworking,omitempty"` // Número de manifiesto
	DNumDesDI     string                      `xml:"dNuDespImp,omitempty"`    // Número de despacho importación
	DInIniTras    string                      `xml:"dIniTras,omitempty"`      // Fecha inicio traslado (yyyy-MM-ddT00:00:00)
	DFesFiworking string                      `xml:"dFinTras,omitempty"`      // Fecha fin estimada traslado
	CPaisDes      types.PaisType              `xml:"cPaisDest,omitempty"`     // País de destino
	DDesPaisDes   string                      `xml:"dDesPaisDest,omitempty"`  // Descripción país destino

	// Lugares de salida y entrega
	GSalida  *TgDirSaliEnt `xml:"gCamSal,omitempty"` // Datos de salida
	GEntrega *TgDirSaliEnt `xml:"gCamEnt,omitempty"` // Datos de entrega

	// Vehículo y transportista
	GVehiculo      *TgVehiculo      `xml:"gVehTras,omitempty"`  // Datos del vehículo
	GTransportista *TgTransportista `xml:"gCamTrans,omitempty"` // Datos del transportista
}

// TgDirSaliEnt: Dirección de Salida/Entrega (E920-E939)
type TgDirSaliEnt struct {
	DDirLoc   string              `xml:"dDirLocSal"`            // Dirección de salida/entrega
	DNumCas   string              `xml:"dNumCasSal,omitempty"`  // Número de casa
	DCompDir1 string              `xml:"dComp1Sal,omitempty"`   // Complemento dirección 1
	DCompDir2 string              `xml:"dComp2Sal,omitempty"`   // Complemento dirección 2
	CDep      types.TDepartamento `xml:"cDepSal,omitempty"`     // Código departamento
	DDesDep   string              `xml:"dDesDepSal,omitempty"`  // Descripción departamento
	CDis      int16               `xml:"cDisSal,omitempty"`     // Código distrito
	DDesDis   string              `xml:"dDesDisSal,omitempty"`  // Descripción distrito
	CCiu      int32               `xml:"cCiuSal,omitempty"`     // Código ciudad
	DDesCiu   string              `xml:"dDesCiuSal,omitempty"`  // Descripción ciudad
	CPais     types.PaisType      `xml:"cPaisSal,omitempty"`    // País
	DDesPais  string              `xml:"dDesPaisSal,omitempty"` // Descripción país
	DTelCont  string              `xml:"dTelSal,omitempty"`     // Teléfono de contacto
}

// TgVehiculo: Datos del Vehículo (E940-E959)
type TgVehiculo struct {
	DTipVeh   string `xml:"dTipVeh"`              // Tipo de vehículo
	DMarca    string `xml:"dMarVeh,omitempty"`    // Marca del vehículo
	DTipIdent int16  `xml:"dTipIdeVeh,omitempty"` // Tipo de identificación vehículo
	DNumIdent string `xml:"dNroIDVeh,omitempty"`  // Número de identificación
	DAdicVeh  string `xml:"dAdicVeh,omitempty"`   // Información adicional del vehículo
	DNumMat   string `xml:"dNroMatVeh,omitempty"` // Número de matrícula
	DNumVuelo string `xml:"dNroVuelo,omitempty"`  // Número de vuelo (aéreo)
}

// TgTransportista: Datos del Transportista (E960-E999)
type TgTransportista struct {
	IContTrans     types.TiNatRec   `xml:"iNatTrans"`              // Contribuyente o no
	DNomTrans      string           `xml:"dNomTrans"`              // Nombre del transportista
	DRucTrans      string           `xml:"dRucTrans,omitempty"`    // RUC del transportista
	DDVTrans       int16            `xml:"dDVTrans,omitempty"`     // Dígito verificador RUC
	ITipIdTrans    types.TTipDocRec `xml:"iTipIDTrans,omitempty"`  // Tipo de documento
	DDesTipIdTrans string           `xml:"dDTipIDTrans,omitempty"` // Descripción tipo documento
	DNumIdTrans    string           `xml:"dNumIDTrans,omitempty"`  // Número de documento
	CPaisTrans     types.PaisType   `xml:"cNacTrans,omitempty"`    // País del transportista
	DDesPaisTrans  string           `xml:"dDesNacTrans,omitempty"` // Descripción país
	DDirTrans      string           `xml:"dDirTrans,omitempty"`    // Dirección del transportista

	// Datos del chofer
	DChofer *TgChofer `xml:"gCamChof,omitempty"` // Datos del chofer
	// Datos del agente
	DAgente *TgAgente `xml:"gCamAgente,omitempty"` // Datos del agente
}

// TgChofer: Datos del Chofer (E980-E989)
type TgChofer struct {
	DNomChofer   string `xml:"dNomChof"`             // Nombre del chofer
	DNumIdChofer string `xml:"dNumIDChof,omitempty"` // Número de documento del chofer
	DDirChofer   string `xml:"dDirChof,omitempty"`   // Dirección del chofer
}

// TgAgente: Datos del Agente (E990-E999)
type TgAgente struct {
	DNomAgente string `xml:"dNomAg,omitempty"` // Nombre del agente
	DRucAgente string `xml:"dRucAg,omitempty"` // RUC del agente
	DDVAgente  int16  `xml:"dDVAg,omitempty"`  // Dígito verificador
	DDirAgente string `xml:"dDirAg,omitempty"` // Dirección del agente
}

// ============================================================================
// TgOpeCom: Campos que describen la operación comercial (D010-D099)
// ============================================================================
type TgOpeCom struct {
	ITipTra     *types.TTipTra `xml:"iTipTra,omitempty"`     // Tipo de transacción
	DDesTipTra  string         `xml:"dDesTipTra,omitempty"`  // Descripción tipo transacción
	ITImp       types.TTImp    `xml:"iTImp"`                 // Tipo de impuesto
	DDesTImp    string         `xml:"dDesTImp"`              // Descripción tipo impuesto
	CMoneOpe    types.CMondT   `xml:"cMoneOpe"`              // Moneda de la operación
	DDesMoneOpe string         `xml:"dDesMoneOpe"`           // Descripción moneda
	DCondTiCam  *int16         `xml:"dCondTiCam,omitempty"`  // Condición tipo de cambio
	DTiCam      *float64       `xml:"dTiCam,omitempty"`      // Tipo de cambio
	ICondAnt    *int16         `xml:"iCondAnt,omitempty"`    // Condición del anticipo
	DDesCondAnt string         `xml:"dDesCondAnt,omitempty"` // Descripción condición anticipo
}

// ============================================================================
// TgEmis: Datos del Emisor (D100-D199)
// ============================================================================
type TgEmis struct {
	DRucEm     string              `xml:"dRucEm"`               // RUC del emisor
	DDVEmi     string              `xml:"dDVEmi"`               // Dígito verificador
	ITipCont   types.TiTipCont     `xml:"iTipCont"`             // Tipo de contribuyente
	CTipReg    *int16              `xml:"cTipReg,omitempty"`    // Código tipo régimen
	DNomEmi    string              `xml:"dNomEmi"`              // Nombre o razón social
	DNomFanEmi string              `xml:"dNomFanEmi,omitempty"` // Nombre de fantasía
	DDirEmi    string              `xml:"dDirEmi"`              // Dirección del emisor
	DNumCas    string              `xml:"dNumCas"`              // Número de casa
	DCompDir1  string              `xml:"dCompDir1,omitempty"`  // Complemento dirección 1
	DCompDir2  string              `xml:"dCompDir2,omitempty"`  // Complemento dirección 2
	CDepEmi    types.TDepartamento `xml:"cDepEmi"`              // Código departamento
	DDesDepEmi string              `xml:"dDesDepEmi"`           // Descripción departamento
	CDisEmi    int16               `xml:"cDisEmi,omitempty"`    // Código distrito
	DDesDisEmi string              `xml:"dDesDisEmi,omitempty"` // Descripción distrito
	CCiuEmi    int32               `xml:"cCiuEmi"`              // Código ciudad
	DDesCiuEmi string              `xml:"dDesCiuEmi"`           // Descripción ciudad
	DTelEmi    string              `xml:"dTelEmi"`              // Teléfono del emisor
	DEmailE    string              `xml:"dEmailE"`              // Email del emisor
	DDenSuc    string              `xml:"dDenSuc,omitempty"`    // Denominación del establecimiento

	GActEcoList []TgActEco `xml:"gActEco"`           // Actividades económicas
	GRespDE     *TgRespDE  `xml:"gRespDE,omitempty"` // Responsable del DE
}

// TgActEco: Actividad Económica del Emisor
type TgActEco struct {
	CActEco    string `xml:"cActEco"`    // Código de actividad económica
	DDesActEco string `xml:"dDesActEco"` // Descripción de la actividad
}

// TgRespDE: Responsable del Documento Electrónico
type TgRespDE struct {
	ITipIDRespDE  int16  `xml:"iTipIDRespDE"`  // Tipo de documento del responsable
	DDTipIDRespDE string `xml:"dDTipIDRespDE"` // Descripción tipo documento
	DNumIDRespDE  string `xml:"dNumIDRespDE"`  // Número de documento
	DNomRespDE    string `xml:"dNomRespDE"`    // Nombre del responsable
	DCarRespDE    string `xml:"dCarRespDE"`    // Cargo del responsable
}

// ============================================================================
// TgDatRec: Datos del Receptor (D200-D299)
// ============================================================================
type TgDatRec struct {
	INatRec     types.TiNatRec       `xml:"iNatRec"`               // Naturaleza del receptor
	ITiOpe      types.TiTiOpe        `xml:"iTiOpe"`                // Tipo de operación
	CPaisRec    types.PaisType       `xml:"cPaisRec"`              // País del receptor
	DDesPaisRe  string               `xml:"dDesPaisRe"`            // Descripción país
	ITiContRec  *types.TiTipCont     `xml:"iTiContRec,omitempty"`  // Tipo de contribuyente receptor
	DRucRec     string               `xml:"dRucRec,omitempty"`     // RUC del receptor
	DDVRec      *int16               `xml:"dDVRec,omitempty"`      // Dígito verificador
	ITipIDRec   *int16               `xml:"iTipIDRec,omitempty"`   // Tipo de documento identidad
	DDTipIDRec  string               `xml:"dDTipIDRec,omitempty"`  // Descripción tipo documento
	DNumIDRec   string               `xml:"dNumIDRec,omitempty"`   // Número de documento
	DNomRec     string               `xml:"dNomRec"`               // Nombre del receptor
	DNomFanRec  string               `xml:"dNomFanRec,omitempty"`  // Nombre de fantasía
	DDirRec     string               `xml:"dDirRec,omitempty"`     // Dirección del receptor
	DNumCasRec  int32                `xml:"dNumCasRec,omitempty"`  // Número de casa
	CDepRec     *types.TDepartamento `xml:"cDepRec,omitempty"`     // Código departamento
	DDesDepRec  string               `xml:"dDesDepRec,omitempty"`  // Descripción departamento
	CDisRec     *int16               `xml:"cDisRec,omitempty"`     // Código distrito
	DDesDisRec  string               `xml:"dDesDisRec,omitempty"`  // Descripción distrito
	CCiuRec     *int32               `xml:"cCiuRec,omitempty"`     // Código ciudad
	DDesCiuRec  string               `xml:"dDesCiuRec,omitempty"`  // Descripción ciudad
	DTelRec     string               `xml:"dTelRec,omitempty"`     // Teléfono del receptor
	DCelRec     string               `xml:"dCelRec,omitempty"`     // Celular del receptor
	DEmailRec   string               `xml:"dEmailRec,omitempty"`   // Email del receptor
	DCodCliente string               `xml:"dCodCliente,omitempty"` // Código interno del cliente
}

// ============================================================================
// TgCamItem: Descripción del Item (E700-E899)
// ============================================================================
type TgCamItem struct {
	DCodInt      string          `xml:"dCodInt"`                // Código interno del item
	DParAranc    int16           `xml:"dParAranc,omitempty"`    // Partida arancelaria
	DNCM         int32           `xml:"dNCM,omitempty"`         // Código NCM
	DDncpG       string          `xml:"dDncpG,omitempty"`       // Código DNCP nivel general
	DDncpE       string          `xml:"dDncpE,omitempty"`       // Código DNCP nivel específico
	DGtin        int64           `xml:"dGtin,omitempty"`        // GTIN del producto
	DGtinPq      int64           `xml:"dGtinPq,omitempty"`      // GTIN del paquete
	DDesProSer   string          `xml:"dDesProSer"`             // Descripción del producto/servicio
	CUniMed      types.TcUniMed  `xml:"cUniMed"`                // Código unidad de medida
	DDesUniMed   string          `xml:"dDesUniMed"`             // Descripción unidad de medida
	DCantProSer  float64         `xml:"dCantProSer"`            // Cantidad del producto/servicio
	CPaisOrig    *types.PaisType `xml:"cPaisOrig,omitempty"`    // País de origen
	DDesPaisOrig string          `xml:"dDesPaisOrig,omitempty"` // Descripción país de origen
	DInfItem     string          `xml:"dInfItem,omitempty"`     // Info adicional del item
	CRelMerc     *int16          `xml:"cRelMerc,omitempty"`     // Relevancia de la mercadería
	DDesRelMerc  string          `xml:"dDesRelMerc,omitempty"`  // Descripción relevancia
	DCanQuiMer   *float64        `xml:"dCanQuiMer,omitempty"`   // Cantidad que acepta
	DPorQuiMer   *float64        `xml:"dPorQuiMer,omitempty"`   // Porcentaje tolerancia merma
	DCDCAnticipo string          `xml:"dCDCAnticipo,omitempty"` // CDC del anticipo

	GValorItem TgValorItem `xml:"gValorItem"`          // Valores del item
	GCamIVA    *TgCamIVA   `xml:"gCamIVA,omitempty"`   // Campos del IVA
	GRasMerc   *TgRasMerc  `xml:"gRasMerc,omitempty"`  // Rastreo de mercadería
	GVehNuevo  *TgVehNuevo `xml:"gVehNuevo,omitempty"` // Sector vehículos
}

// TgValorItem: Valores del Item (E720-E729)
type TgValorItem struct {
	DPUniProSer     float64          `xml:"dPUniProSer"`        // Precio unitario
	DTiCamIt        *float64         `xml:"dTiCamIt,omitempty"` // Tipo de cambio por item
	DTotBruOpeItem  float64          `xml:"dTotBruOpeItem"`     // Total bruto de la operación
	GValorRestaItem TgValorRestaItem `xml:"gValorRestaItem"`    // Valores restantes
}

// TgValorRestaItem: Descuentos y Anticipos del Item (EA001-EA050)
type TgValorRestaItem struct {
	DDescItem       *float64 `xml:"dDescItem,omitempty"`       // Descuento particular del item
	DPorcDesIt      *float64 `xml:"dPorcDesIt,omitempty"`      // Porcentaje de descuento
	DDescGloItem    *float64 `xml:"dDescGloItem,omitempty"`    // Descuento global del item
	DAntPreUniIt    *float64 `xml:"dAntPreUniIt,omitempty"`    // Anticipo particular del item
	DAntGloPreUniIt *float64 `xml:"dAntGloPreUniIt,omitempty"` // Anticipo global del item
	DTotOpeItem     float64  `xml:"dTotOpeItem"`               // Total de la operación
	DTotOpeGs       *float64 `xml:"dTotOpeGs,omitempty"`       // Total operación en Guaraníes
}

// TgCamIVA: Campos del IVA por Item (E730-E739)
type TgCamIVA struct {
	IAfecIVA    types.TiAfecIVA `xml:"iAfecIVA"`          // Afectación tributaria IVA
	DDesAfecIVA string          `xml:"dDesAfecIVA"`       // Descripción afectación
	DPropIVA    float64         `xml:"dPropIVA"`          // Proporción gravada de IVA
	DTasaIVA    float64         `xml:"dTasaIVA"`          // Tasa del IVA
	DBasGravIVA float64         `xml:"dBasGravIVA"`       // Base gravada del IVA
	DLiqIVAItem float64         `xml:"dLiqIVAItem"`       // Liquidación del IVA
	DBasExe     *float64        `xml:"dBasExe,omitempty"` // Base exenta
}

// ============================================================================
// TgRasMerc: Rastreo de Mercadería (E750-E760)
// ============================================================================
type TgRasMerc struct {
	DNLote    string `xml:"dNLote,omitempty"`    // Número de lote
	DVencMerc string `xml:"dVencMerc,omitempty"` // Fecha de vencimiento (yyyy-MM-dd)
	DNSerie   string `xml:"dNSerie,omitempty"`   // Número de serie
	DNPedido  string `xml:"dNPedido,omitempty"`  // Número de pedido
	DNSeguim  string `xml:"dNSeguim,omitempty"`  // Número de seguimiento
	// Campos del importador
	GCamImp     *TgCamImp `xml:"dInfoImport,omitempty"` // Datos del importador
	DRegistroS  string    `xml:"dNRegSenave,omitempty"` // Registro SENAVE
	DRegistroEn string    `xml:"dNRegEntCom,omitempty"` // Registro entidad comercial
}

// TgCamImp: Datos del Importador
type TgCamImp struct {
	DNomImp  string `xml:"dNomImp,omitempty"` // Nombre del importador
	DDirImp  string `xml:"dDirImp,omitempty"` // Dirección del importador
	DNRegImp string `xml:"dNumReg,omitempty"` // Número de registro
}

// ============================================================================
// TgVehNuevo: Sector de Vehículos Nuevos/Usados (E770-E789)
// ============================================================================
type TgVehNuevo struct {
	ITipOpVN    int16                   `xml:"iTipOpVN"`                 // Tipo de operación de vehículo
	DDesTipOpVN string                  `xml:"dDesTipOpVN"`              // Descripción tipo operación
	DChasis     string                  `xml:"dChworking,omitempty"`     // Número de chasis
	DColor      string                  `xml:"dColor,omitempty"`         // Color del vehículo
	DPotWorking int32                   `xml:"dPotVeh,omitempty"`        // Potencia del motor (HP)
	DCapMot     int32                   `xml:"dCapMot,omitempty"`        // Capacidad del motor (CC)
	DPNet       float64                 `xml:"dPNet,omitempty"`          // Peso neto
	DPBrut      float64                 `xml:"dPBrut,omitempty"`         // Peso bruto
	ITipCom     types.TiTipoCombustible `xml:"iTipComb,omitempty"`       // Tipo de combustible
	DDesTipCom  string                  `xml:"dDesTipComb,omitempty"`    // Descripción tipo combustible
	DNMotor     string                  `xml:"dNroMotor,omitempty"`      // Número de motor
	DCapTracc   float64                 `xml:"dCapTraworking,omitempty"` // Capacidad de tracción
	DAnoFab     int16                   `xml:"dAnoFab,omitempty"`        // Año de fabricación
	DTipVeh     string                  `xml:"dTipVeh,omitempty"`        // Tipo de vehículo
	DCap        int16                   `xml:"dCap,omitempty"`           // Capacidad (pasajeros)
	DCil        string                  `xml:"dCil,omitempty"`           // Cilindrada
}

// ============================================================================
// TgTotSub: Totales del Documento (F001-F099)
// ============================================================================
type TgTotSub struct {
	DSubExe        float64  `xml:"dSubExe"`                // Subtotal exentas
	DSubExo        float64  `xml:"dSubExo"`                // Subtotal exoneradas
	DSub5          float64  `xml:"dSub5,omitempty"`        // Subtotal 5%
	DSub10         float64  `xml:"dSub10,omitempty"`       // Subtotal 10%
	DTotOpe        float64  `xml:"dTotOpe"`                // Total bruto de la operación
	DTotDesc       float64  `xml:"dTotDesc"`               // Total de descuentos
	DTotDescGlotem float64  `xml:"dTotDescGlotem"`         // Total descuento global
	DTotAntItem    float64  `xml:"dTotAntItem"`            // Total anticipo por item
	DTotAnt        float64  `xml:"dTotAnt"`                // Total de anticipos
	DPorcDescTotal float64  `xml:"dPorcDescTotal"`         // Porcentaje descuento total
	DDescTotal     float64  `xml:"dDescTotal"`             // Descuento total
	DAnticipo      float64  `xml:"dAnticipo"`              // Anticipo
	DRedon         float64  `xml:"dRedon"`                 // Redondeo
	DComi          *float64 `xml:"dComi,omitempty"`        // Comisión
	DTotGralOpe    float64  `xml:"dTotGralOpe"`            // Total general de la operación
	DIVA5          float64  `xml:"dIVA5,omitempty"`        // IVA 5%
	DIVA10         float64  `xml:"dIVA10,omitempty"`       // IVA 10%
	DLiqTotIVA5    float64  `xml:"dLiqTotIVA5,omitempty"`  // Liquidación total IVA 5%
	DLiqTotIVA10   float64  `xml:"dLiqTotIVA10,omitempty"` // Liquidación total IVA 10%
	DIVAComi       *float64 `xml:"dIVAComi,omitempty"`     // IVA de la comisión
	DTotIVA        float64  `xml:"dTotIVA,omitempty"`      // Total IVA
	DBaseGrav5     float64  `xml:"dBaseGrav5,omitempty"`   // Base gravada 5%
	DBaseGrav10    float64  `xml:"dBaseGrav10,omitempty"`  // Base gravada 10%
	DTBasGraIVA    float64  `xml:"dTBasGraIVA,omitempty"`  // Total base gravada IVA
	DTotalGs       *float64 `xml:"dTotalGs,omitempty"`     // Total en Guaraníes
}

// ============================================================================
// TgCamGen: Campos Generales Complementarios (G001-G099)
// ============================================================================
type TgCamGen struct {
	DOrdCompra string     `xml:"dOrdCompra,omitempty"` // Número de orden de compra
	DOrdVta    string     `xml:"dOrdVta,omitempty"`    // Número de orden de venta
	DAsiento   string     `xml:"dAsiento,omitempty"`   // Número de asiento contable
	GCamCarg   *TgCamCarg `xml:"gCamCarg,omitempty"`   // Campos de carga
}

// TgCamCarg: Datos de Carga
type TgCamCarg struct {
	CUniMedTotVol    types.TcUniMed   `xml:"cUniMedTotVol,omitempty"`    // Unidad de medida volumen
	DDesUniMedTotVol string           `xml:"dDesUniMedTotVol,omitempty"` // Descripción unidad
	DTotVolMerc      int64            `xml:"dTotVolMerc,omitempty"`      // Volumen total mercadería
	CUniMedTotPes    types.TcUniMed   `xml:"cUniMedTotPes,omitempty"`    // Unidad de medida peso
	DDesUniMedTotPes string           `xml:"dDesUniMedTotPes,omitempty"` // Descripción unidad
	DTotPesMerc      int64            `xml:"dTotPesMerc,omitempty"`      // Peso total mercadería
	ICarCarga        types.TiCarCarga `xml:"iCarCarga,omitempty"`        // Característica de la carga
	DDesCarCarga     string           `xml:"dDesCarCarga,omitempty"`     // Descripción característica
}

// ============================================================================
// TgCamDEAsoc: Documento Electrónico Asociado (H001-H049)
// ============================================================================
type TgCamDEAsoc struct {
	ITipDocAso    types.TiTipDocAso `xml:"iTipDocAso"`             // Tipo de documento asociado
	DDesTipDocAso string            `xml:"dDesTipDocAso"`          // Descripción tipo documento
	DCdCDERef     string            `xml:"dCdCDERef,omitempty"`    // CDC del DE referenciado
	DNTimDI       string            `xml:"dNTimDI,omitempty"`      // Timbrado del documento impreso
	DEstDocAso    string            `xml:"dEstDocAso,omitempty"`   // Establecimiento del doc asociado
	DPExpDocAso   string            `xml:"dPExpDocAso,omitempty"`  // Punto de expedición
	DNumDocAso    string            `xml:"dNumDocAso,omitempty"`   // Número del documento asociado
	ITipoDocAso   *types.TiTIpoDoc  `xml:"iTipoDocAso,omitempty"`  // Tipo de documento impreso
	DDTipoDocAso  string            `xml:"dDTipoDocAso,omitempty"` // Descripción tipo documento
	DFecEmiDI     string            `xml:"dFecEmiDI,omitempty"`    // Fecha emisión (yyyy-MM-dd)
	DNumComRet    string            `xml:"dNumComRet,omitempty"`   // Número de comprobante retención
	DNumResCF     string            `xml:"dNumResCF,omitempty"`    // Número resolución crédito fiscal
	ITipCons      *types.TdTipCons  `xml:"iTipCons,omitempty"`     // Tipo de constancia
	DDesTipCons   string            `xml:"dDesTipCons,omitempty"`  // Descripción tipo constancia
	DNumCons      *int64            `xml:"dNumCons,omitempty"`     // Número de constancia
	DNumControl   string            `xml:"dNumControl,omitempty"`  // Número de control de la constancia
}
