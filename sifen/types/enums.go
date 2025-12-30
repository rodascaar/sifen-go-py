package types

import "fmt"

// ============================================================================
// TTiDE: Tipo de Documento Electronico (Document Type)
// ============================================================================
type TTiDE int16

const (
	TTiDE_FacturaElectronica              TTiDE = 1
	TTiDE_FacturaElectronicaExportacion   TTiDE = 2
	TTiDE_FacturaElectronicaImportacion   TTiDE = 3
	TTiDE_AutofacturaElectronica          TTiDE = 4
	TTiDE_NotaCreditoElectronica          TTiDE = 5
	TTiDE_NotaDebitoElectronica           TTiDE = 6
	TTiDE_NotaRemisionElectronica         TTiDE = 7
	TTiDE_ComprobanteRetencionElectronico TTiDE = 8
)

func (t TTiDE) String() string {
	switch t {
	case TTiDE_FacturaElectronica:
		return "Factura electrónica"
	case TTiDE_FacturaElectronicaExportacion:
		return "Factura electrónica de exportación"
	case TTiDE_FacturaElectronicaImportacion:
		return "Factura electrónica de importación"
	case TTiDE_AutofacturaElectronica:
		return "Autofactura electrónica"
	case TTiDE_NotaCreditoElectronica:
		return "Nota de crédito electrónica"
	case TTiDE_NotaDebitoElectronica:
		return "Nota de débito electrónica"
	case TTiDE_NotaRemisionElectronica:
		return "Nota de remisión electrónica"
	case TTiDE_ComprobanteRetencionElectronico:
		return "Comprobante de retención electrónico"
	default:
		return fmt.Sprintf("Tipo desconocido: %d", t)
	}
}

// ============================================================================
// TTipEmi: Tipo de Emision (Emission Type)
// ============================================================================
type TTipEmi int16

const (
	TTipEmi_Normal       TTipEmi = 1
	TTipEmi_Contingencia TTipEmi = 2
)

func (t TTipEmi) String() string {
	switch t {
	case TTipEmi_Normal:
		return "Normal"
	case TTipEmi_Contingencia:
		return "Contingencia"
	default:
		return fmt.Sprintf("%d", t)
	}
}

// ============================================================================
// TTipTra: Tipo de Transaccion (Transaction Type)
// ============================================================================
type TTipTra int16

const (
	TTipTra_VentaMercaderia     TTipTra = 1
	TTipTra_PrestacionServicios TTipTra = 2
	TTipTra_Mixto               TTipTra = 3
	TTipTra_VentaActivoFijo     TTipTra = 4
	TTipTra_VentaDivisas        TTipTra = 5
	TTipTra_CompraDivisas       TTipTra = 6
	TTipTra_PromocionMuestras   TTipTra = 7
	TTipTra_Donacion            TTipTra = 8
	TTipTra_Anticipo            TTipTra = 9
	TTipTra_CompraProductos     TTipTra = 10
	TTipTra_CompraServicios     TTipTra = 11
	TTipTra_VentaCreditoFiscal  TTipTra = 12
	TTipTra_MuestrasMedicas     TTipTra = 13
)

func (t TTipTra) String() string {
	switch t {
	case TTipTra_VentaMercaderia:
		return "Venta de mercadería"
	case TTipTra_PrestacionServicios:
		return "Prestación de servicios"
	case TTipTra_Mixto:
		return "Mixto (Venta de mercadería y servicios)"
	case TTipTra_VentaActivoFijo:
		return "Venta de activo fijo"
	case TTipTra_VentaDivisas:
		return "Venta de divisas"
	case TTipTra_CompraDivisas:
		return "Compra de divisas"
	case TTipTra_PromocionMuestras:
		return "Promoción o entrega de muestras"
	case TTipTra_Donacion:
		return "Donación"
	case TTipTra_Anticipo:
		return "Anticipo"
	case TTipTra_CompraProductos:
		return "Compra de productos"
	case TTipTra_CompraServicios:
		return "Compra de servicios"
	case TTipTra_VentaCreditoFiscal:
		return "Venta de crédito fiscal"
	case TTipTra_MuestrasMedicas:
		return "Muestras médicas (Art. 3 RG 24/2014)"
	default:
		return fmt.Sprintf("%d", t)
	}
}

// ============================================================================
// TTImp: Tipo de Impuesto (Tax Type)
// ============================================================================
type TTImp int16

const (
	TTImp_IVA      TTImp = 1
	TTImp_ISC      TTImp = 2
	TTImp_Renta    TTImp = 3
	TTImp_Ninguno  TTImp = 4
	TTImp_IVARenta TTImp = 5
)

func (t TTImp) String() string {
	switch t {
	case TTImp_IVA:
		return "IVA"
	case TTImp_ISC:
		return "ISC"
	case TTImp_Renta:
		return "Renta"
	case TTImp_Ninguno:
		return "Ninguno"
	case TTImp_IVARenta:
		return "IVA - Renta"
	default:
		return fmt.Sprintf("%d", t)
	}
}

// ============================================================================
// TiTipCont: Tipo de Contribuyente (Taxpayer Type)
// ============================================================================
type TiTipCont int16

const (
	TiTipCont_PersonaFisica   TiTipCont = 1
	TiTipCont_PersonaJuridica TiTipCont = 2
)

func (t TiTipCont) String() string {
	switch t {
	case TiTipCont_PersonaFisica:
		return "Persona Física"
	case TiTipCont_PersonaJuridica:
		return "Persona Jurídica"
	default:
		return fmt.Sprintf("%d", t)
	}
}

// ============================================================================
// TiNatRec: Naturaleza del Receptor (Receiver Nature)
// ============================================================================
type TiNatRec int16

const (
	TiNatRec_Contribuyente   TiNatRec = 1
	TiNatRec_NoContribuyente TiNatRec = 2
)

func (t TiNatRec) String() string {
	switch t {
	case TiNatRec_Contribuyente:
		return "Contribuyente"
	case TiNatRec_NoContribuyente:
		return "No Contribuyente"
	default:
		return fmt.Sprintf("%d", t)
	}
}

// ============================================================================
// TiTiOpe: Tipo de Operacion (Operation Type - B2B/B2C/B2G/B2F)
// ============================================================================
type TiTiOpe int16

const (
	TiTiOpe_B2B TiTiOpe = 1
	TiTiOpe_B2C TiTiOpe = 2
	TiTiOpe_B2G TiTiOpe = 3
	TiTiOpe_B2F TiTiOpe = 4
)

func (t TiTiOpe) String() string {
	switch t {
	case TiTiOpe_B2B:
		return "B2B"
	case TiTiOpe_B2C:
		return "B2C"
	case TiTiOpe_B2G:
		return "B2G"
	case TiTiOpe_B2F:
		return "B2F"
	default:
		return fmt.Sprintf("%d", t)
	}
}

// ============================================================================
// TiAfecIVA: Afectacion Tributaria IVA (IVA Tax Affectation)
// ============================================================================
type TiAfecIVA int16

const (
	TiAfecIVA_GravadoIVA     TiAfecIVA = 1
	TiAfecIVA_Exonerado      TiAfecIVA = 2
	TiAfecIVA_Exento         TiAfecIVA = 3
	TiAfecIVA_GravadoParcial TiAfecIVA = 4
)

func (t TiAfecIVA) String() string {
	switch t {
	case TiAfecIVA_GravadoIVA:
		return "Gravado IVA"
	case TiAfecIVA_Exonerado:
		return "Exonerado (Art. 83-Ley 125/91)"
	case TiAfecIVA_Exento:
		return "Exento"
	case TiAfecIVA_GravadoParcial:
		return "Gravado parcial (Grav-Exento)"
	default:
		return fmt.Sprintf("%d", t)
	}
}

// ============================================================================
// TiIndPres: Indicador de Presencia (Presence Indicator)
// ============================================================================
type TiIndPres int16

const (
	TiIndPres_Presencial     TiIndPres = 1
	TiIndPres_Electronica    TiIndPres = 2
	TiIndPres_Telemarketing  TiIndPres = 3
	TiIndPres_VentaDomicilio TiIndPres = 4
	TiIndPres_Bancaria       TiIndPres = 5
	TiIndPres_Ciclica        TiIndPres = 6
	TiIndPres_Otro           TiIndPres = 9
)

func (t TiIndPres) String() string {
	switch t {
	case TiIndPres_Presencial:
		return "Operación presencial"
	case TiIndPres_Electronica:
		return "Operación electrónica"
	case TiIndPres_Telemarketing:
		return "Operación telemarketing"
	case TiIndPres_VentaDomicilio:
		return "Venta a domicilio"
	case TiIndPres_Bancaria:
		return "Operación bancaria"
	case TiIndPres_Ciclica:
		return "Operación cíclica"
	case TiIndPres_Otro:
		return "Otro"
	default:
		return fmt.Sprintf("%d", t)
	}
}

// ============================================================================
// TiCondOpe: Condicion de la Operacion (Operation Condition)
// ============================================================================
type TiCondOpe int16

const (
	TiCondOpe_Contado TiCondOpe = 1
	TiCondOpe_Credito TiCondOpe = 2
)

func (t TiCondOpe) String() string {
	switch t {
	case TiCondOpe_Contado:
		return "Contado"
	case TiCondOpe_Credito:
		return "Crédito"
	default:
		return fmt.Sprintf("%d", t)
	}
}

// ============================================================================
// TiCondCredito: Condicion de Credito (Credit Condition)
// ============================================================================
type TiCondCredito int16

const (
	TiCondCredito_Plazo  TiCondCredito = 1
	TiCondCredito_Cuotas TiCondCredito = 2
)

func (t TiCondCredito) String() string {
	switch t {
	case TiCondCredito_Plazo:
		return "Plazo"
	case TiCondCredito_Cuotas:
		return "Cuotas"
	default:
		return fmt.Sprintf("%d", t)
	}
}

// ============================================================================
// TiTipPago: Tipo de Pago (Payment Type)
// ============================================================================
type TiTipPago int16

const (
	TiTipPago_Efectivo            TiTipPago = 1
	TiTipPago_Cheque              TiTipPago = 2
	TiTipPago_TarjetaCredito      TiTipPago = 3
	TiTipPago_TarjetaDebito       TiTipPago = 4
	TiTipPago_TransferenciaBanco  TiTipPago = 5
	TiTipPago_GirosBancarios      TiTipPago = 6
	TiTipPago_BilleteraMobile     TiTipPago = 7
	TiTipPago_CreditosFiscales    TiTipPago = 8
	TiTipPago_Voucher             TiTipPago = 9
	TiTipPago_RetencionParcial    TiTipPago = 10
	TiTipPago_RetencionTotal      TiTipPago = 11
	TiTipPago_PagoPorAnticipo     TiTipPago = 12
	TiTipPago_ValorFiscal         TiTipPago = 13
	TiTipPago_ValorComercial      TiTipPago = 14
	TiTipPago_Compensacion        TiTipPago = 15
	TiTipPago_Permuta             TiTipPago = 16
	TiTipPago_PagoMovilBilletera  TiTipPago = 17
	TiTipPago_InterbancariaCuenta TiTipPago = 18
	TiTipPago_InterbancariaTarj   TiTipPago = 19
	TiTipPago_Otro                TiTipPago = 99
)

func (t TiTipPago) String() string {
	switch t {
	case TiTipPago_Efectivo:
		return "Efectivo"
	case TiTipPago_Cheque:
		return "Cheque"
	case TiTipPago_TarjetaCredito:
		return "Tarjeta de crédito"
	case TiTipPago_TarjetaDebito:
		return "Tarjeta de débito"
	case TiTipPago_TransferenciaBanco:
		return "Transferencia bancaria"
	case TiTipPago_GirosBancarios:
		return "Giros bancarios"
	case TiTipPago_BilleteraMobile:
		return "Billetera electrónica"
	default:
		return fmt.Sprintf("Tipo pago: %d", t)
	}
}

// ============================================================================
// TiMotEmiNC: Motivo de Emision de Nota de Credito
// ============================================================================
type TiMotEmiNC int16

const (
	TiMotEmiNC_DevolucionAjuste  TiMotEmiNC = 1
	TiMotEmiNC_Devolucion        TiMotEmiNC = 2
	TiMotEmiNC_Descuento         TiMotEmiNC = 3
	TiMotEmiNC_Bonificacion      TiMotEmiNC = 4
	TiMotEmiNC_CreditoIncobrable TiMotEmiNC = 5
	TiMotEmiNC_RecuperoCosto     TiMotEmiNC = 6
	TiMotEmiNC_RecuperoGasto     TiMotEmiNC = 7
	TiMotEmiNC_AjustePrecio      TiMotEmiNC = 8
)

func (t TiMotEmiNC) String() string {
	switch t {
	case TiMotEmiNC_DevolucionAjuste:
		return "Devolución y Ajuste de precios"
	case TiMotEmiNC_Devolucion:
		return "Devolución"
	case TiMotEmiNC_Descuento:
		return "Descuento"
	case TiMotEmiNC_Bonificacion:
		return "Bonificación"
	case TiMotEmiNC_CreditoIncobrable:
		return "Crédito incobrable"
	case TiMotEmiNC_RecuperoCosto:
		return "Recupero de costo"
	case TiMotEmiNC_RecuperoGasto:
		return "Recupero de gasto"
	case TiMotEmiNC_AjustePrecio:
		return "Ajuste de precio"
	default:
		return fmt.Sprintf("%d", t)
	}
}

// ============================================================================
// TiMotEmiNR: Motivo de Emision de Nota de Remision
// ============================================================================
type TiMotEmiNR int16

const (
	TiMotEmiNR_TrasladoVentas         TiMotEmiNR = 1
	TiMotEmiNR_TrasladoConsignacion   TiMotEmiNR = 2
	TiMotEmiNR_Exportacion            TiMotEmiNR = 3
	TiMotEmiNR_TrasladoCompra         TiMotEmiNR = 4
	TiMotEmiNR_Importacion            TiMotEmiNR = 5
	TiMotEmiNR_TrasladoDevolucion     TiMotEmiNR = 6
	TiMotEmiNR_TrasladoEntreLocales   TiMotEmiNR = 7
	TiMotEmiNR_TrasladoTransformacion TiMotEmiNR = 8
	TiMotEmiNR_TrasladoReparacion     TiMotEmiNR = 9
	TiMotEmiNR_TrasladoEmisorMovil    TiMotEmiNR = 10
	TiMotEmiNR_ExhibicionDemostracion TiMotEmiNR = 11
	TiMotEmiNR_ParticipacionFerias    TiMotEmiNR = 12
	TiMotEmiNR_TrasladoEncomienda     TiMotEmiNR = 13
	TiMotEmiNR_Decomiso               TiMotEmiNR = 14
	TiMotEmiNR_Otro                   TiMotEmiNR = 99
)

func (t TiMotEmiNR) String() string {
	switch t {
	case TiMotEmiNR_TrasladoVentas:
		return "Traslado por ventas"
	case TiMotEmiNR_TrasladoConsignacion:
		return "Traslado por consignación"
	case TiMotEmiNR_Exportacion:
		return "Exportación"
	case TiMotEmiNR_TrasladoCompra:
		return "Traslado por compra"
	case TiMotEmiNR_Importacion:
		return "Importación"
	case TiMotEmiNR_TrasladoDevolucion:
		return "Traslado por devolución"
	case TiMotEmiNR_TrasladoEntreLocales:
		return "Traslado entre locales de la empresa"
	case TiMotEmiNR_TrasladoTransformacion:
		return "Traslado de bienes por transformación"
	case TiMotEmiNR_TrasladoReparacion:
		return "Traslado de bienes por reparación"
	case TiMotEmiNR_TrasladoEmisorMovil:
		return "Traslado por emisor móvil"
	case TiMotEmiNR_ExhibicionDemostracion:
		return "Exhibición o demostración"
	case TiMotEmiNR_ParticipacionFerias:
		return "Participación en ferias"
	case TiMotEmiNR_TrasladoEncomienda:
		return "Traslado de encomienda"
	case TiMotEmiNR_Decomiso:
		return "Decomiso"
	case TiMotEmiNR_Otro:
		return "Otro"
	default:
		return fmt.Sprintf("%d", t)
	}
}

// ============================================================================
// TiRespFlete: Responsable del Flete
// ============================================================================
type TiRespFlete int16

const (
	TiRespFlete_EmisorFactura         TiRespFlete = 1
	TiRespFlete_PoseedorFacturaBienes TiRespFlete = 2
	TiRespFlete_EmpresaTransportista  TiRespFlete = 3
	TiRespFlete_DespachanteAduanas    TiRespFlete = 4
	TiRespFlete_AgenteTransporte      TiRespFlete = 5
)

func (t TiRespFlete) String() string {
	switch t {
	case TiRespFlete_EmisorFactura:
		return "Emisor de la factura"
	case TiRespFlete_PoseedorFacturaBienes:
		return "Poseedor de la factura y bienes"
	case TiRespFlete_EmpresaTransportista:
		return "Empresa transportista"
	case TiRespFlete_DespachanteAduanas:
		return "Despachante de Aduanas"
	case TiRespFlete_AgenteTransporte:
		return "Agente de transporte o intermediario"
	default:
		return fmt.Sprintf("%d", t)
	}
}

// ============================================================================
// TiNatVendedorAF: Naturaleza del Vendedor en Autofactura
// ============================================================================
type TiNatVendedorAF int16

const (
	TiNatVendedorAF_NoContribuyente TiNatVendedorAF = 1
	TiNatVendedorAF_Extranjero      TiNatVendedorAF = 2
)

func (t TiNatVendedorAF) String() string {
	switch t {
	case TiNatVendedorAF_NoContribuyente:
		return "No contribuyente"
	case TiNatVendedorAF_Extranjero:
		return "Extranjero"
	default:
		return fmt.Sprintf("%d", t)
	}
}

// ============================================================================
// TTipDocRec: Tipo de Documento del Receptor
// ============================================================================
type TTipDocRec int16

const (
	TTipDocRec_CedulaParaguaya    TTipDocRec = 1
	TTipDocRec_Pasaporte          TTipDocRec = 2
	TTipDocRec_CedulaExtranjera   TTipDocRec = 3
	TTipDocRec_CarnetResidencia   TTipDocRec = 4
	TTipDocRec_Innominado         TTipDocRec = 5
	TTipDocRec_TarjetaDiplomatica TTipDocRec = 6
	TTipDocRec_Otro               TTipDocRec = 9
)

func (t TTipDocRec) String() string {
	switch t {
	case TTipDocRec_CedulaParaguaya:
		return "Cédula paraguaya"
	case TTipDocRec_Pasaporte:
		return "Pasaporte"
	case TTipDocRec_CedulaExtranjera:
		return "Cédula extranjera"
	case TTipDocRec_CarnetResidencia:
		return "Carnet de residencia"
	case TTipDocRec_Innominado:
		return "Innominado"
	case TTipDocRec_TarjetaDiplomatica:
		return "Tarjeta Diplomática de exoneración fiscal"
	case TTipDocRec_Otro:
		return "Otro"
	default:
		return fmt.Sprintf("%d", t)
	}
}

// ============================================================================
// TDepartamento: Departamento (Paraguay Departments)
// ============================================================================
type TDepartamento int16

const (
	TDepartamento_Capital         TDepartamento = 1
	TDepartamento_Concepcion      TDepartamento = 2
	TDepartamento_SanPedro        TDepartamento = 3
	TDepartamento_Cordillera      TDepartamento = 4
	TDepartamento_Guaira          TDepartamento = 5
	TDepartamento_Caaguazu        TDepartamento = 6
	TDepartamento_Caazapa         TDepartamento = 7
	TDepartamento_Itapua          TDepartamento = 8
	TDepartamento_Misiones        TDepartamento = 9
	TDepartamento_Paraguari       TDepartamento = 10
	TDepartamento_AltoParana      TDepartamento = 11
	TDepartamento_Central         TDepartamento = 12
	TDepartamento_Neembucu        TDepartamento = 13
	TDepartamento_Amambay         TDepartamento = 14
	TDepartamento_Canindeyu       TDepartamento = 15
	TDepartamento_PresidenteHayes TDepartamento = 16
	TDepartamento_Boqueron        TDepartamento = 17
	TDepartamento_AltoParaguay    TDepartamento = 18
)

func (t TDepartamento) String() string {
	switch t {
	case TDepartamento_Capital:
		return "CAPITAL"
	case TDepartamento_Concepcion:
		return "CONCEPCION"
	case TDepartamento_SanPedro:
		return "SAN PEDRO"
	case TDepartamento_Cordillera:
		return "CORDILLERA"
	case TDepartamento_Guaira:
		return "GUAIRA"
	case TDepartamento_Caaguazu:
		return "CAAGUAZU"
	case TDepartamento_Caazapa:
		return "CAAZAPA"
	case TDepartamento_Itapua:
		return "ITAPUA"
	case TDepartamento_Misiones:
		return "MISIONES"
	case TDepartamento_Paraguari:
		return "PARAGUARI"
	case TDepartamento_AltoParana:
		return "ALTO PARANA"
	case TDepartamento_Central:
		return "CENTRAL"
	case TDepartamento_Neembucu:
		return "ÑEEMBUCU"
	case TDepartamento_Amambay:
		return "AMAMBAY"
	case TDepartamento_Canindeyu:
		return "CANINDEYU"
	case TDepartamento_PresidenteHayes:
		return "PRESIDENTE HAYES"
	case TDepartamento_Boqueron:
		return "BOQUERON"
	case TDepartamento_AltoParaguay:
		return "ALTO PARAGUAY"
	default:
		return fmt.Sprintf("Departamento %d", t)
	}
}

// ============================================================================
// TcUniMed: Unidad de Medida (Measurement Unit)
// ============================================================================
type TcUniMed int16

const (
	TcUniMed_Unidad        TcUniMed = 77
	TcUniMed_Hora          TcUniMed = 15
	TcUniMed_Kilogramo     TcUniMed = 71
	TcUniMed_Metros        TcUniMed = 87
	TcUniMed_MetroCuadrado TcUniMed = 79
	TcUniMed_MetroCubico   TcUniMed = 80
	TcUniMed_Litro         TcUniMed = 76
	TcUniMed_Par           TcUniMed = 94
	TcUniMed_Docena        TcUniMed = 56
	TcUniMed_Pieza         TcUniMed = 98
	TcUniMed_Global        TcUniMed = 66
	TcUniMed_Dia           TcUniMed = 54
	TcUniMed_Mes           TcUniMed = 88
	TcUniMed_Año           TcUniMed = 110
	TcUniMed_Servicio      TcUniMed = 117
)

func (u TcUniMed) String() string {
	switch u {
	case TcUniMed_Unidad:
		return "Unidad"
	case TcUniMed_Hora:
		return "Hora"
	case TcUniMed_Kilogramo:
		return "Kilogramo"
	case TcUniMed_Metros:
		return "Metro"
	case TcUniMed_MetroCuadrado:
		return "Metro cuadrado"
	case TcUniMed_MetroCubico:
		return "Metro cúbico"
	case TcUniMed_Litro:
		return "Litro"
	case TcUniMed_Par:
		return "Par"
	case TcUniMed_Docena:
		return "Docena"
	case TcUniMed_Pieza:
		return "Pieza"
	case TcUniMed_Global:
		return "Global"
	case TcUniMed_Dia:
		return "Día"
	case TcUniMed_Mes:
		return "Mes"
	case TcUniMed_Año:
		return "Año"
	case TcUniMed_Servicio:
		return "Servicio"
	default:
		return fmt.Sprintf("Unidad %d", u)
	}
}

func (u TcUniMed) Abreviatura() string {
	switch u {
	case TcUniMed_Unidad:
		return "UNI"
	case TcUniMed_Hora:
		return "h"
	case TcUniMed_Kilogramo:
		return "kg"
	case TcUniMed_Metros:
		return "m"
	case TcUniMed_MetroCuadrado:
		return "m²"
	case TcUniMed_MetroCubico:
		return "m³"
	case TcUniMed_Litro:
		return "L"
	case TcUniMed_Par:
		return "PAR"
	case TcUniMed_Docena:
		return "DOC"
	case TcUniMed_Pieza:
		return "PZA"
	case TcUniMed_Global:
		return "GLO"
	case TcUniMed_Dia:
		return "DIA"
	case TcUniMed_Mes:
		return "MES"
	case TcUniMed_Año:
		return "AÑO"
	case TcUniMed_Servicio:
		return "SER"
	default:
		return "UNI"
	}
}

// ============================================================================
// CMondT: Codigo de Moneda (Currency Code - ISO 4217)
// ============================================================================
type CMondT string

const (
	CMondT_PYG CMondT = "PYG" // Guaraní
	CMondT_USD CMondT = "USD" // US Dollar
	CMondT_EUR CMondT = "EUR" // Euro
	CMondT_BRL CMondT = "BRL" // Brazilian Real
	CMondT_ARS CMondT = "ARS" // Argentine Peso
	CMondT_UYU CMondT = "UYU" // Peso Uruguayo
	CMondT_CLP CMondT = "CLP" // Chilean Peso
	CMondT_BOB CMondT = "BOB" // Boliviano
	CMondT_PEN CMondT = "PEN" // Nuevo Sol
	CMondT_COP CMondT = "COP" // Colombian Peso
	CMondT_GBP CMondT = "GBP" // Pound Sterling
	CMondT_JPY CMondT = "JPY" // Yen
	CMondT_CHF CMondT = "CHF" // Swiss Franc
	CMondT_CAD CMondT = "CAD" // Canadian Dollar
	CMondT_AUD CMondT = "AUD" // Australian Dollar
	CMondT_CNY CMondT = "CNY" // Yuan Renminbi
)

func (c CMondT) String() string { return string(c) }

func (c CMondT) Codigo() string { return string(c) }

// ============================================================================
// PaisType: Codigo de Pais (Country Code - ISO 3166-1 Alpha-3)
// ============================================================================
type PaisType string

const (
	PaisType_PRY PaisType = "PRY" // Paraguay
	PaisType_USA PaisType = "USA" // United States
	PaisType_ARG PaisType = "ARG" // Argentina
	PaisType_BRA PaisType = "BRA" // Brazil
	PaisType_URY PaisType = "URY" // Uruguay
	PaisType_BOL PaisType = "BOL" // Bolivia
	PaisType_CHL PaisType = "CHL" // Chile
	PaisType_PER PaisType = "PER" // Peru
	PaisType_COL PaisType = "COL" // Colombia
	PaisType_ECU PaisType = "ECU" // Ecuador
	PaisType_VEN PaisType = "VEN" // Venezuela
	PaisType_MEX PaisType = "MEX" // Mexico
	PaisType_ESP PaisType = "ESP" // Spain
	PaisType_DEU PaisType = "DEU" // Germany
	PaisType_FRA PaisType = "FRA" // France
	PaisType_ITA PaisType = "ITA" // Italy
	PaisType_GBR PaisType = "GBR" // United Kingdom
	PaisType_CHN PaisType = "CHN" // China
	PaisType_JPN PaisType = "JPN" // Japan
	PaisType_KOR PaisType = "KOR" // South Korea
)

func (p PaisType) String() string { return string(p) }
func (p PaisType) Codigo() string { return string(p) }
func (p PaisType) Nombre() string {
	switch p {
	case PaisType_PRY:
		return "Paraguay"
	case PaisType_USA:
		return "Estados Unidos"
	case PaisType_ARG:
		return "Argentina"
	case PaisType_BRA:
		return "Brasil"
	case PaisType_URY:
		return "Uruguay"
	case PaisType_BOL:
		return "Bolivia"
	case PaisType_CHL:
		return "Chile"
	case PaisType_PER:
		return "Perú"
	case PaisType_COL:
		return "Colombia"
	case PaisType_ECU:
		return "Ecuador"
	case PaisType_VEN:
		return "Venezuela"
	case PaisType_MEX:
		return "México"
	case PaisType_ESP:
		return "España"
	case PaisType_DEU:
		return "Alemania"
	case PaisType_FRA:
		return "Francia"
	case PaisType_ITA:
		return "Italia"
	case PaisType_GBR:
		return "Reino Unido"
	case PaisType_CHN:
		return "China"
	case PaisType_JPN:
		return "Japón"
	case PaisType_KOR:
		return "Corea del Sur"
	default:
		return string(p)
	}
}

// ============================================================================
// TiCarCarga: Caracteristicas de la Carga (Load Characteristics)
// ============================================================================
type TiCarCarga int16

const (
	TiCarCarga_MercaderiasConCadenaFrio TiCarCarga = 1
	TiCarCarga_CargaPeligrosa           TiCarCarga = 2
	TiCarCarga_Otras                    TiCarCarga = 3
)

func (t TiCarCarga) String() string {
	switch t {
	case TiCarCarga_MercaderiasConCadenaFrio:
		return "Mercaderías con cadena de frío"
	case TiCarCarga_CargaPeligrosa:
		return "Carga peligrosa"
	case TiCarCarga_Otras:
		return "Otras"
	default:
		return fmt.Sprintf("%d", t)
	}
}

// ============================================================================
// TiTipDocAso: Tipo de Documento Asociado
// ============================================================================
type TiTipDocAso int16

const (
	TiTipDocAso_Electronico           TiTipDocAso = 1
	TiTipDocAso_Impreso               TiTipDocAso = 2
	TiTipDocAso_ConstanciaElectronica TiTipDocAso = 3
)

func (t TiTipDocAso) String() string {
	switch t {
	case TiTipDocAso_Electronico:
		return "Electrónico"
	case TiTipDocAso_Impreso:
		return "Impreso"
	case TiTipDocAso_ConstanciaElectronica:
		return "Constancia electrónica"
	default:
		return fmt.Sprintf("%d", t)
	}
}

// ============================================================================
// TiTIpoDoc: Tipo de Documento Impreso (Printed Document Type)
// ============================================================================
type TiTIpoDoc int16

const (
	TiTIpoDoc_Factura              TiTIpoDoc = 1
	TiTIpoDoc_NotaCredito          TiTIpoDoc = 2
	TiTIpoDoc_NotaDebito           TiTIpoDoc = 3
	TiTIpoDoc_NotaRemision         TiTIpoDoc = 4
	TiTIpoDoc_ComprobanteRetencion TiTIpoDoc = 5
)

func (t TiTIpoDoc) String() string {
	switch t {
	case TiTIpoDoc_Factura:
		return "Factura"
	case TiTIpoDoc_NotaCredito:
		return "Nota de crédito"
	case TiTIpoDoc_NotaDebito:
		return "Nota de débito"
	case TiTIpoDoc_NotaRemision:
		return "Nota de remisión"
	case TiTIpoDoc_ComprobanteRetencion:
		return "Comprobante de retención"
	default:
		return fmt.Sprintf("%d", t)
	}
}

// ============================================================================
// TdTipCons: Tipo de Constancia
// ============================================================================
type TdTipCons int16

const (
	TdTipCons_ConstanciaNoRetencion      TdTipCons = 1
	TdTipCons_ConstanciaMicroproductores TdTipCons = 2
)

func (t TdTipCons) String() string {
	switch t {
	case TdTipCons_ConstanciaNoRetencion:
		return "Constancia de no retención"
	case TdTipCons_ConstanciaMicroproductores:
		return "Constancia de microproductores"
	default:
		return fmt.Sprintf("%d", t)
	}
}

// ============================================================================
// TTiEvento: Tipo de Evento (Event Type)
// ============================================================================
type TTiEvento int16

const (
	TTiEvento_Cancelacion             TTiEvento = 1
	TTiEvento_Inutilizacion           TTiEvento = 2
	TTiEvento_Endoso                  TTiEvento = 3
	TTiEvento_AcuseDE                 TTiEvento = 10
	TTiEvento_ConformidadDE           TTiEvento = 11
	TTiEvento_DisconformidadDE        TTiEvento = 12
	TTiEvento_DesconocimientoDE       TTiEvento = 13
	TTiEvento_NotificacionRecepcion   TTiEvento = 14
	TTiEvento_Nominacion              TTiEvento = 20
	TTiEvento_ActualizacionTransporte TTiEvento = 21
)

func (t TTiEvento) String() string {
	switch t {
	case TTiEvento_Cancelacion:
		return "Cancelación"
	case TTiEvento_Inutilizacion:
		return "Inutilización"
	case TTiEvento_Endoso:
		return "Endoso"
	case TTiEvento_AcuseDE:
		return "Acuse del DE"
	case TTiEvento_ConformidadDE:
		return "Conformidad del DE"
	case TTiEvento_DisconformidadDE:
		return "Disconformidad del DE"
	case TTiEvento_DesconocimientoDE:
		return "Desconocimiento del DE"
	case TTiEvento_NotificacionRecepcion:
		return "Notificación de recepción"
	case TTiEvento_Nominacion:
		return "Nominación"
	case TTiEvento_ActualizacionTransporte:
		return "Actualización de datos de transporte"
	default:
		return fmt.Sprintf("Evento %d", t)
	}
}

// ============================================================================
// TiTipoTransporte: Tipo de Transporte
// ============================================================================
type TiTipoTransporte int16

const (
	TiTipoTransporte_Propio  TiTipoTransporte = 1
	TiTipoTransporte_Tercero TiTipoTransporte = 2
)

func (t TiTipoTransporte) String() string {
	switch t {
	case TiTipoTransporte_Propio:
		return "Propio"
	case TiTipoTransporte_Tercero:
		return "Tercero"
	default:
		return fmt.Sprintf("%d", t)
	}
}

// ============================================================================
// TiModalidadTransporte: Modalidad de Transporte
// ============================================================================
type TiModalidadTransporte int16

const (
	TiModalidadTransporte_Terrestre  TiModalidadTransporte = 1
	TiModalidadTransporte_Fluvial    TiModalidadTransporte = 2
	TiModalidadTransporte_Aereo      TiModalidadTransporte = 3
	TiModalidadTransporte_Multimodal TiModalidadTransporte = 4
)

func (t TiModalidadTransporte) String() string {
	switch t {
	case TiModalidadTransporte_Terrestre:
		return "Terrestre"
	case TiModalidadTransporte_Fluvial:
		return "Fluvial"
	case TiModalidadTransporte_Aereo:
		return "Aéreo"
	case TiModalidadTransporte_Multimodal:
		return "Multimodal"
	default:
		return fmt.Sprintf("%d", t)
	}
}

// ============================================================================
// TiTipVehiculo: Tipo de Vehículo
// ============================================================================
type TiTipVehiculo int16

const (
	TiTipVehiculo_Camion    TiTipVehiculo = 1
	TiTipVehiculo_Camioneta TiTipVehiculo = 2
	TiTipVehiculo_Furgoneta TiTipVehiculo = 3
	TiTipVehiculo_Barco     TiTipVehiculo = 4
	TiTipVehiculo_Avion     TiTipVehiculo = 5
)

func (t TiTipVehiculo) String() string {
	switch t {
	case TiTipVehiculo_Camion:
		return "Camión"
	case TiTipVehiculo_Camioneta:
		return "Camioneta"
	case TiTipVehiculo_Furgoneta:
		return "Furgoneta"
	case TiTipVehiculo_Barco:
		return "Barco"
	case TiTipVehiculo_Avion:
		return "Avión"
	default:
		return fmt.Sprintf("%d", t)
	}
}

// ============================================================================
// TiTipoCombustible: Tipo de Combustible (para sector automotor)
// ============================================================================
type TiTipoCombustible int16

const (
	TiTipoCombustible_Nafta     TiTipoCombustible = 1
	TiTipoCombustible_Diesel    TiTipoCombustible = 2
	TiTipoCombustible_Gas       TiTipoCombustible = 3
	TiTipoCombustible_Electrico TiTipoCombustible = 4
	TiTipoCombustible_Hibrido   TiTipoCombustible = 5
	TiTipoCombustible_NaftaGas  TiTipoCombustible = 6
	TiTipoCombustible_DieselGas TiTipoCombustible = 7
	TiTipoCombustible_Alcohol   TiTipoCombustible = 8
	TiTipoCombustible_Vapor     TiTipoCombustible = 9
)

func (t TiTipoCombustible) String() string {
	switch t {
	case TiTipoCombustible_Nafta:
		return "Nafta"
	case TiTipoCombustible_Diesel:
		return "Diésel"
	case TiTipoCombustible_Gas:
		return "Gas"
	case TiTipoCombustible_Electrico:
		return "Eléctrico"
	case TiTipoCombustible_Hibrido:
		return "Híbrido"
	case TiTipoCombustible_NaftaGas:
		return "Nafta/Gas"
	case TiTipoCombustible_DieselGas:
		return "Diésel/Gas"
	case TiTipoCombustible_Alcohol:
		return "Alcohol"
	case TiTipoCombustible_Vapor:
		return "Vapor"
	default:
		return fmt.Sprintf("%d", t)
	}
}

// ============================================================================
// TiTipRegimen: Tipo de Régimen
// ============================================================================
type TiTipRegimen int16

const (
	TiTipRegimen_Turismo          TiTipRegimen = 1
	TiTipRegimen_Importador       TiTipRegimen = 2
	TiTipRegimen_Exportador       TiTipRegimen = 3
	TiTipRegimen_Maquila          TiTipRegimen = 4
	TiTipRegimen_Ley60_90         TiTipRegimen = 5
	TiTipRegimen_PequenoProductor TiTipRegimen = 6
	TiTipRegimen_MedianoProductor TiTipRegimen = 7
	TiTipRegimen_Contable         TiTipRegimen = 8
)

func (t TiTipRegimen) String() string {
	switch t {
	case TiTipRegimen_Turismo:
		return "Régimen de Turismo"
	case TiTipRegimen_Importador:
		return "Importador"
	case TiTipRegimen_Exportador:
		return "Exportador"
	case TiTipRegimen_Maquila:
		return "Maquila"
	case TiTipRegimen_Ley60_90:
		return "Ley N° 60/90"
	case TiTipRegimen_PequenoProductor:
		return "Régimen del Pequeño Productor"
	case TiTipRegimen_MedianoProductor:
		return "Régimen del Mediano Productor"
	case TiTipRegimen_Contable:
		return "Régimen Contable"
	default:
		return fmt.Sprintf("%d", t)
	}
}

// ============================================================================
// TiTipoConformidad: Tipo de Conformidad
// ============================================================================
type TiTipoConformidad int16

const (
	TiTipoConformidad_Total   TiTipoConformidad = 1
	TiTipoConformidad_Parcial TiTipoConformidad = 2
)

func (t TiTipoConformidad) String() string {
	switch t {
	case TiTipoConformidad_Total:
		return "Conformidad total"
	case TiTipoConformidad_Parcial:
		return "Conformidad parcial"
	default:
		return fmt.Sprintf("%d", t)
	}
}
