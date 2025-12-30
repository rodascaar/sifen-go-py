package kude

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"html/template"
	"time"

	"github.com/rodascaar/sifen-go-py/sifen/models"
	"github.com/rodascaar/sifen-go-py/sifen/types"
)

// ============================================================================
// KuDE Generator - Generador de Representación Gráfica del DE
// ============================================================================

// KuDEConfig contiene la configuración para generar el KuDE
type KuDEConfig struct {
	BaseURLQR        string // URL base para QR (ej: https://ekuatia.set.gov.py/consultas/qr?)
	CSC              string // Código de Seguridad del Contribuyente
	IdCSC            string // ID del CSC
	LogoEmisorPath   string // Ruta opcional al logo del emisor
	LogoEmisorBase64 string // O logo en Base64
}

// KuDEData contiene todos los datos necesarios para generar el KuDE
type KuDEData struct {
	// Identificación del documento
	CDC             string
	TipoDocumento   string
	NumeroTimbrado  string
	Establecimiento string
	PuntoExpedicion string
	NumeroDocumento string
	FechaEmision    time.Time
	FechaFirma      time.Time

	// Emisor
	RUCEmisor       string
	DVEmisor        string
	NombreEmisor    string
	NombreFantasia  string
	DireccionEmisor string
	TelefonoEmisor  string
	EmailEmisor     string

	// Receptor
	RUCReceptor       string
	DVReceptor        string
	NombreReceptor    string
	DireccionReceptor string
	TelefonoReceptor  string
	EmailReceptor     string

	// Items
	Items []KuDEItem

	// Totales
	SubtotalExentas float64
	SubtotalIVA5    float64
	SubtotalIVA10   float64
	TotalIVA5       float64
	TotalIVA10      float64
	TotalIVA        float64
	TotalGeneral    float64
	Moneda          string
	MontoLetras     string // Monto total en letras

	// Condición de pago
	CondicionPago string
	FormaPago     string

	// Información Adicional
	InformacionInteres string // Campo J003
	ActividadEconomica string // Descripción de actividad económica

	// QR generado
	URLCompleta  string
	QRCodeBase64 string // QR como imagen Base64 (si se genera externamente)
}

// KuDEItem representa un item del documento
type KuDEItem struct {
	Codigo      string
	Descripcion string
	Cantidad    float64
	Unidad      string
	PrecioUnit  float64
	Descuento   float64
	Exenta      float64
	IVA5        float64
	IVA10       float64
}

// KuDEGenerator genera representaciones gráficas del DE
type KuDEGenerator struct {
	config KuDEConfig
}

// NewKuDEGenerator crea un nuevo generador de KuDE
func NewKuDEGenerator(config KuDEConfig) *KuDEGenerator {
	return &KuDEGenerator{config: config}
}

// GenerateQRURL genera la URL completa para el código QR según especificación SIFEN
func (g *KuDEGenerator) GenerateQRURL(data KuDEData) string {
	// Construir cadena para hash
	// Formato: nVersion=150&Id=CDC&dFeEmiDE=fecha&dRucRec=RUC&dTotGralOpe=total&dTotIVA=iva&cItems=items&DigestValue=digest&IdCSC=idcsc&cHashQR=hash

	fechaEmision := data.FechaEmision.Format("2006-01-02T15:04:05")

	// Concatenar para generar hash
	// El hash se calcula sobre los parámetros + CSC
	paramsForHash := fmt.Sprintf("nVersion=150&Id=%s&dFeEmiDE=%s&dRucRec=%s&dTotGralOpe=%.2f&dTotIVA=%.2f&cItems=%d&DigestValue=&IdCSC=%s",
		data.CDC,
		fechaEmision,
		data.RUCReceptor,
		data.TotalGeneral,
		data.TotalIVA,
		len(data.Items),
		g.config.IdCSC,
	)

	// Calcular SHA-256 del string + CSC
	hashInput := paramsForHash + g.config.CSC
	hash := sha256.Sum256([]byte(hashInput))
	hashHex := hex.EncodeToString(hash[:])

	// URL completa
	return fmt.Sprintf("%s%s&cHashQR=%s", g.config.BaseURLQR, paramsForHash, hashHex)
}

// GenerateFromDE genera KuDEData desde un DocumentoElectronico
func (g *KuDEGenerator) GenerateFromDE(de *models.DocumentoElectronico) KuDEData {
	data := KuDEData{
		CDC:             de.DE.Id,
		TipoDocumento:   de.DE.GTimb.ITiDE.String(),
		NumeroTimbrado:  fmt.Sprintf("%d", de.DE.GTimb.DNumTim),
		Establecimiento: de.DE.GTimb.DEst,
		PuntoExpedicion: de.DE.GTimb.DPunExp,
		NumeroDocumento: de.DE.GTimb.DNumDoc,
		FechaEmision:    time.Now(), // This should ideally be de.DE.GDatGralOpe.DFeEmiDE if parsed correctly

		// Emisor
		RUCEmisor:       de.DE.GDatGralOpe.GEmis.DRucEm,
		DVEmisor:        de.DE.GDatGralOpe.GEmis.DDVEmi,
		NombreEmisor:    de.DE.GDatGralOpe.GEmis.DNomEmi,
		NombreFantasia:  de.DE.GDatGralOpe.GEmis.DNomFanEmi,
		DireccionEmisor: de.DE.GDatGralOpe.GEmis.DDirEmi,
		TelefonoEmisor:  de.DE.GDatGralOpe.GEmis.DTelEmi,
		EmailEmisor:     de.DE.GDatGralOpe.GEmis.DEmailE,

		// Receptor
		RUCReceptor:       de.DE.GDatGralOpe.GDatRec.DRucRec,
		NombreReceptor:    de.DE.GDatGralOpe.GDatRec.DNomRec,
		DireccionReceptor: de.DE.GDatGralOpe.GDatRec.DDirRec,
		TelefonoReceptor:  de.DE.GDatGralOpe.GDatRec.DTelRec,
		EmailReceptor:     de.DE.GDatGralOpe.GDatRec.DEmailRec,

		Moneda: string(de.DE.GDatGralOpe.GOpeCom.CMoneOpe),

		// Información de interés (si existe en los campos adicionales del DE, aquí se deja vacío por defecto)
		InformacionInteres: "",
	}

	// Actividad Económica
	if len(de.DE.GDatGralOpe.GEmis.GActEcoList) > 0 {
		data.ActividadEconomica = de.DE.GDatGralOpe.GEmis.GActEcoList[0].DDesActEco
	}

	// Items
	for _, item := range de.DE.GDtipDE.GCamItemList {
		kudeItem := KuDEItem{
			Codigo:      item.DCodInt,
			Descripcion: item.DDesProSer,
			Cantidad:    item.DCantProSer,
			Unidad:      item.DDesUniMed,
			PrecioUnit:  item.GValorItem.DPUniProSer,
		}

		// Calcular por afectación IVA
		if item.GCamIVA != nil {
			switch item.GCamIVA.IAfecIVA {
			case types.TiAfecIVA_Exento, types.TiAfecIVA_Exonerado:
				kudeItem.Exenta = item.GValorItem.GValorRestaItem.DTotOpeItem
			default:
				if item.GCamIVA.DTasaIVA == 5 {
					kudeItem.IVA5 = item.GValorItem.GValorRestaItem.DTotOpeItem
				} else {
					kudeItem.IVA10 = item.GValorItem.GValorRestaItem.DTotOpeItem
				}
			}
		}

		data.Items = append(data.Items, kudeItem)
	}

	// Totales
	if de.DE.GTotSub != nil {
		data.SubtotalExentas = de.DE.GTotSub.DSubExe
		data.SubtotalIVA5 = de.DE.GTotSub.DSub5
		data.SubtotalIVA10 = de.DE.GTotSub.DSub10
		data.TotalIVA5 = de.DE.GTotSub.DLiqTotIVA5
		data.TotalIVA10 = de.DE.GTotSub.DLiqTotIVA10
		data.TotalIVA = de.DE.GTotSub.DTotIVA
		data.TotalGeneral = de.DE.GTotSub.DTotGralOpe
	}

	// Monto en Letras
	data.MontoLetras = numeroALetras(data.TotalGeneral, data.Moneda)

	// Generar URL QR
	data.URLCompleta = g.GenerateQRURL(data)

	return data
}

// GenerateHTML genera el KuDE en formato HTML
func (g *KuDEGenerator) GenerateHTML(data KuDEData) (string, error) {
	tmpl, err := template.New("kude").Funcs(template.FuncMap{
		"formatMoney": formatMoney,
		"formatDate":  formatDate,
		"formatCDC":   formatCDC,
	}).Parse(kudeHTMLTemplate)

	if err != nil {
		return "", fmt.Errorf("error al parsear template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("error al ejecutar template: %w", err)
	}

	return buf.String(), nil
}

// GenerateHTMLBase64 genera el KuDE HTML codificado en Base64 (útil para embeber)
func (g *KuDEGenerator) GenerateHTMLBase64(data KuDEData) (string, error) {
	html, err := g.GenerateHTML(data)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString([]byte(html)), nil
}

// Helpers
func formatMoney(amount float64) string {
	// Formato manual con separador de miles
	// Go fmt no soporta %V para miles nativamente de forma estándar en todas las versiones
	s := fmt.Sprintf("%.0f", amount)
	if len(s) <= 3 {
		return s
	}
	var result string
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result += "."
		}
		result += string(c)
	}
	return result
}

func formatDate(t time.Time) string {
	return t.Format("02/01/2006 15:04:05")
}

func formatCDC(cdc string) string {
	// Formatear en grupos de 4 con espacios
	var result string
	for i, char := range cdc {
		if i > 0 && i%4 == 0 {
			result += " "
		}
		result += string(char)
	}
	return result
}

// ============================================================================
// Template HTML del KuDE
// ============================================================================

const kudeHTMLTemplate = `<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>KuDE - {{.TipoDocumento}}</title>
    <style>
        @page {
            size: A4; /* Formato estándar definido */
            margin: 10mm;
            @bottom-right {
                content: "Página " counter(page) " de " counter(pages);
                font-family: Arial, sans-serif;
                font-size: 8px;
            }
        }
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: 'Segoe UI', Arial, sans-serif;
            font-size: 11px; /* Letra ligeramente más pequeña para ajustarse mejor */
            line-height: 1.3;
            color: #000; /* Color negro puro para impresión */
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            background: #fff;
        }
        .kude-container {
            border: 2px solid #000;
            border-radius: 5px;
            padding: 10px;
            background: white;
        }
        
        /* Header Grid */
        .header-grid {
            display: grid;
            grid-template-columns: 3fr 2fr;
            gap: 15px;
            border-bottom: 2px solid #000;
            padding-bottom: 10px;
            margin-bottom: 10px;
        }
        
        /* Emisor Info */
        .emisor-info {
            text-align: left;
        }
        .emisor-logo {
            max-width: 200px;
            max-height: 80px;
            margin-bottom: 10px;
        }
        .emisor-name {
            font-size: 16px;
            font-weight: bold;
            color: #000;
            text-transform: uppercase;
        }
        .emisor-fantasia {
            font-size: 14px;
            font-weight: bold;
            font-style: italic;
        }
        .emisor-details {
            margin-top: 5px;
            font-size: 10px;
        }
        
        /* Boxed Info (RUC/Timbrado) */
        .doc-info-box {
            border: 1px solid #000;
            border-radius: 4px;
            padding: 10px;
            text-align: left;
            position: relative;
        }
        .doc-title {
            text-align: center;
            font-weight: bold;
            font-size: 14px;
            margin: 10px 0;
            text-transform: uppercase;
        }
        .doc-number {
            text-align: center;
            font-size: 16px;
            font-weight: bold;
            margin-top: 5px;
        }
        
        /* Info Sections Boxed */
        .box-section {
            border: 1px solid #000;
            border-radius: 4px;
            padding: 8px;
            margin-bottom: 10px;
        }
        
        .info-grid {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 10px;
        }
        .info-row {
            display: flex;
            margin-bottom: 2px;
        }
        .label {
            font-weight: bold;
            margin-right: 5px;
            min-width: 80px;
        }
        
        /* Table */
        table {
            width: 100%;
            border-collapse: collapse;
            font-size: 10px;
            margin-bottom: 10px;
        }
        th, td {
            border: 1px solid #000;
            padding: 5px;
        }
        th {
            background-color: #f0f0f0;
            font-weight: bold;
            text-align: center;
        }
        .col-right { text-align: right; }
        .col-center { text-align: center; }
        
        /* Totals */
        .totals-section {
            border: 1px solid #000;
            border-radius: 4px;
            margin-bottom: 10px;
        }
        .total-row {
            display: flex;
            justify-content: space-between;
            padding: 5px 10px;
            border-bottom: 1px solid #ddd;
        }
        .total-row:last-child {
            border-bottom: none;
        }
        .total-row.main-total {
            background-color: #e0e0e0;
            font-weight: bold;
            font-size: 12px;
            border-top: 1px solid #000;
        }
        
        /* Letras */
        .monto-letras-row {
            border-top: 1px solid #000;
            padding: 5px 10px;
            font-weight: bold;
            font-size: 11px;
            background-color: #f9f9f9;
        }
        
        /* IVA */
        .iva-row {
            border-top: 1px solid #000;
            display: flex;
            justify-content: space-between;
            padding: 5px 10px;
            font-size: 10px;
            background-color: #f0f0f0;
        }

        /* Footer / QR */
        .footer-grid {
            display: flex;
            border: 1px solid #000;
            border-radius: 4px;
            padding: 10px;
            align-items: center;
        }
        .qr-container {
            margin-right: 20px;
            min-width: 100px;
            text-align: center;
        }
        .qr-container img {
            width: 100px;
            height: 100px;
        }
        .footer-text {
            flex: 1;
            font-size: 10px;
        }
        .cdc-text {
            font-family: monospace;
            font-weight: bold;
            font-size: 11px;
            margin-top: 5px;
            letter-spacing: 1px;
        }
        
        .powered-by {
            text-align: right;
            font-size: 8px;
            color: #666;
            margin-top: 5px;
        }

        @media print {
            body { 
                padding: 0; 
                background: white; 
            }
            .kude-container { border: 2px solid #000; }
        }
    </style>
</head>
<body>
    <div class="kude-container">
        <!-- Header -->
        <div class="header-grid">
            <div class="emisor-info">
                {{if .NombreFantasia}}
                    <div class="emisor-fantasia">{{.NombreFantasia}}</div>
                {{end}}
                <div class="emisor-name">{{.NombreEmisor}}</div>
                <div class="emisor-details">
                    <div>{{.ActividadEconomica}}</div>
                    <div><strong>Dirección:</strong> {{.DireccionEmisor}}</div>
                    <div><strong>Teléfono:</strong> {{.TelefonoEmisor}}</div>
                    <div><strong>Email:</strong> {{.EmailEmisor}}</div>
                </div>
            </div>
            
            <div class="doc-info-box">
                <div><strong>RUC:</strong> {{.RUCEmisor}}-{{.DVEmisor}}</div>
                <div><strong>Timbrado Nº:</strong> {{.NumeroTimbrado}}</div>
                <div class="doc-title">Factura Electrónica</div>
                <div class="doc-number">Nro: {{.Establecimiento}}-{{.PuntoExpedicion}}-{{.NumeroDocumento}}</div>
            </div>
        </div>

        <!-- Info Principal -->
        <div class="box-section">
            <div class="info-grid">
                <div>
                    <div class="info-row"><span class="label">Fecha Emisión:</span> {{formatDate .FechaEmision}}</div>
                    <div class="info-row"><span class="label">Condición:</span> {{.CondicionPago}}</div>
                    <div class="info-row"><span class="label">Moneda:</span> {{.Moneda}}</div>
                </div>
                <div>
                    <div class="info-row"><span class="label">Vencimiento:</span> -</div>
                </div>
            </div>
        </div>

        <!-- Cliente -->
        <div class="box-section">
            <div class="info-grid">
                <div>
                    <div class="info-row"><span class="label">Nombre / Razón Social:</span> {{.NombreReceptor}}</div>
                    <div class="info-row"><span class="label">RUC / CI:</span> {{.RUCReceptor}}{{if .DVReceptor}}-{{.DVReceptor}}{{end}}</div>
                    <div class="info-row"><span class="label">Dirección:</span> {{.DireccionReceptor}}</div>
                </div>
                <div>
                    <div class="info-row"><span class="label">Teléfono:</span> {{.TelefonoReceptor}}</div>
                    <div class="info-row"><span class="label">Email:</span> {{.EmailReceptor}}</div>
                </div>
            </div>
        </div>

        <!-- Items Table -->
        <table>
            <thead>
                <tr>
                    <th>Código</th>
                    <th>Descripción</th>
                    <th>Cant.</th>
                    <th>Unid.</th>
                    <th>Precio Unit.</th>
                    <th>Exenta</th>
                    <th>5%</th>
                    <th>10%</th>
                </tr>
            </thead>
            <tbody>
                {{range .Items}}
                <tr>
                    <td>{{.Código}}</td>
                    <td>{{.Descripcion}}</td>
                    <td class="col-center">{{.Cantidad}}</td>
                    <td class="col-center">{{.Unidad}}</td>
                    <td class="col-right">{{formatMoney .PrecioUnit}}</td>
                    <td class="col-right">{{if .Exenta}}{{formatMoney .Exenta}}{{end}}</td>
                    <td class="col-right">{{if .IVA5}}{{formatMoney .IVA5}}{{end}}</td>
                    <td class="col-right">{{if .IVA10}}{{formatMoney .IVA10}}{{end}}</td>
                </tr>
                {{end}}
            </tbody>
        </table>

        <!-- Totals -->
        <div class="totals-section">
            <div class="total-row">
                <span>SUBTOTALES:</span>
                <span class="col-right" style="width: 100px;">{{formatMoney .SubtotalExentas}}</span>
                <span class="col-right" style="width: 100px;">{{formatMoney .SubtotalIVA5}}</span>
                <span class="col-right" style="width: 100px;">{{formatMoney .SubtotalIVA10}}</span>
            </div>
            
            <div class="monto-letras-row">
                TOTAL A PAGAR: {{.MontoLetras}}
            </div>
            
            <div class="total-row main-total">
                <span>TOTAL GENERAL ({{.Moneda}}):</span>
                <span>{{formatMoney .TotalGeneral}}</span>
            </div>
            
            <div class="iva-row">
                <span><strong>LIQUIDACIÓN DEL IVA:</strong></span>
                <span>(5%) {{formatMoney .TotalIVA5}}</span>
                <span>(10%) {{formatMoney .TotalIVA10}}</span>
                <span><strong>TOTAL IVA:</strong> {{formatMoney .TotalIVA}}</span>
            </div>
        </div>
        
        {{if .InformacionInteres}}
        <div class="box-section">
            <strong>Información de Interés:</strong> {{.InformacionInteres}}
        </div>
        {{end}}

        <!-- Footer / QR -->
        <div class="footer-grid">
            <div class="qr-container">
                {{if .QRCodeBase64}}
                    <img src="data:image/png;base64,{{.QRCodeBase64}}" alt="QR Code" />
                {{else}}
                    <!-- Placeholder -->
                    <div style="width:100px;height:100px;border:1px solid #ccc;display:flex;align-items:center;justify-content:center;">QR</div>
                {{end}}
            </div>
            <div class="footer-text">
                <div style="margin-bottom: 5px;">
                    Consulte la validez de esta Factura Electrónica con el número de CDC impreso abajo en:<br>
                    <a href="https://ekuatia.set.gov.py/consultas" style="color: #000; text-decoration: none;">https://ekuatia.set.gov.py/consultas</a>
                </div>
                <div><strong>CDC (Código de Control):</strong></div>
                <div class="cdc-text">{{formatCDC .CDC}}</div>
                <div style="margin-top: 5px; font-size: 9px;">
                    ESTE DOCUMENTO ES UNA REPRESENTACIÓN GRÁFICA DE UN DOCUMENTO ELECTRÓNICO (XML).<br>
                    Si su documento presenta algún error, podrá solicitar la modificación dentro de las 72hs.
                </div>
            </div>
        </div>
        
        <div class="powered-by">
            Sistema de Facturación Electrónica SIFEN-Go
        </div>
    </div>
</body>
</html>`

// GenerateQRCodeSVG genera un placeholder SVG del QR
func GenerateQRCodeSVG(url string) string {
	return `<!-- SVG Placeholder -->`
}

// ============================================================================
// Numero A Letras (Simplificado para Guaraníes)
// ============================================================================

func numeroALetras(n float64, moneda string) string {
	// Implementación básica para demo. En producción usar librería robusta.
	// Asumimos enteros para PYG mayormente.
	entero := int64(n)
	letras := convertirNumero(entero)

	monedaNombre := "GUARANIES"
	if moneda == "USD" {
		monedaNombre = "DOLARES"
	}

	return fmt.Sprintf("%s %s", letras, monedaNombre)
}

func convertirNumero(n int64) string {
	if n == 0 {
		return "CERO"
	}
	if n < 0 {
		return "MENOS " + convertirNumero(-n)
	}

	unidades := []string{"", "UN", "DOS", "TRES", "CUATRO", "CINCO", "SEIS", "SIETE", "OCHO", "NUEVE"}
	decenas := []string{"", "DIEZ", "VEINTE", "TREINTA", "CUARENTA", "CINCUENTA", "SESENTA", "SETENTA", "OCHENTA", "NOVENTA"}
	dieces := []string{"DIEZ", "ONCE", "DOCE", "TRECE", "CATORCE", "QUINCE", "DIECISEIS", "DIECISIETE", "DIECIOCHO", "DIECINUEVE"}
	centenas := []string{"", "CIENTO", "DOSCIENTOS", "TRESCIENTOS", "CUATROCIENTOS", "QUINIENTOS", "SEISCIENTOS", "SETECIENTOS", "OCHOCIENTOS", "NOVECIENTOS"}

	if n < 10 {
		return unidades[n]
	}
	if n < 20 {
		return dieces[n-10]
	}
	if n < 100 {
		d := n / 10
		u := n % 10
		if u == 0 {
			return decenas[d]
		}
		if d == 2 {
			return "VEINTI" + unidades[u]
		}
		return decenas[d] + " Y " + unidades[u]
	}
	if n < 1000 {
		c := n / 100
		r := n % 100
		if n == 100 {
			return "CIEN"
		}
		return centenas[c] + " " + convertirNumero(r)
	}
	if n < 1000000 {
		m := n / 1000
		r := n % 1000
		if m == 1 {
			return "MIL " + convertirNumero(r)
		}
		return convertirNumero(m) + " MIL " + convertirNumero(r)
	}
	if n < 1000000000000 { // Billones
		m := n / 1000000
		r := n % 1000000
		if m == 1 {
			return "UN MILLON " + convertirNumero(r)
		}
		return convertirNumero(m) + " MILLONES " + convertirNumero(r)
	}

	return fmt.Sprintf("%d", n) // Fallback para números muy grandes
}
