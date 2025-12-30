package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// LeftPad adds padding characters to the left of a string
func LeftPad(s string, padChar rune, length int) string {
	if len(s) >= length {
		return s
	}
	return strings.Repeat(string(padChar), length-len(s)) + s
}

// RightPad adds padding characters to the right of a string
func RightPad(s string, padChar rune, length int) string {
	if len(s) >= length {
		return s
	}
	return s + strings.Repeat(string(padChar), length-len(s))
}

// ============================================================================
// CDC Generation
// ============================================================================

// CDCParams contains parameters for CDC generation
type CDCParams struct {
	TipoDocumento     int16  // iTiDE (1-8)
	RUC               string // RUC del contribuyente
	DigitoVerificador string // DV del RUC
	Establecimiento   string // Código de establecimiento (3 chars)
	PuntoExpedicion   string // Punto de expedición (3 chars)
	NumeroDocumento   string // Número de documento (7 chars)
	TipoContribuyente int16  // 1=Persona Física, 2=Persona Jurídica
	FechaEmision      time.Time
	TipoEmision       int16  // 1=Normal, 2=Contingencia
	CodigoSeguridad   string // Código de seguridad (9 chars)
}

// GenerateCDC generates the CDC (Código de Control) for a document
// CDC = 44 characters
func GenerateCDC(params CDCParams) (string, error) {
	// Validate inputs
	if params.RUC == "" {
		return "", fmt.Errorf("RUC is required")
	}
	if len(params.Establecimiento) != 3 {
		return "", fmt.Errorf("establecimiento must be 3 characters")
	}
	if len(params.PuntoExpedicion) != 3 {
		return "", fmt.Errorf("punto expedicion must be 3 characters")
	}
	if len(params.NumeroDocumento) != 7 {
		return "", fmt.Errorf("numero documento must be 7 characters")
	}
	if len(params.CodigoSeguridad) != 9 {
		return "", fmt.Errorf("codigo seguridad must be 9 characters")
	}

	// Build CDC components
	// Format: DD+RRRRRRRRRR+DV+EEE+PPP+NNNNNNN+T+E+AAAAMMDD+SSSSSSSSS+V
	// DD = Tipo documento (2 digits)
	// RRRRRRRRRR = RUC (8 digits, padded)
	// DV = Dígito verificador (1 digit)
	// EEE = Establecimiento (3 digits)
	// PPP = Punto expedición (3 digits)
	// NNNNNNN = Número documento (7 digits)
	// T = Tipo contribuyente (1 digit)
	// E = Tipo emisión (1 digit)
	// AAAAMMDD = Fecha emisión (8 digits)
	// SSSSSSSSS = Código seguridad (9 digits)
	// V = Dígito verificador CDC (1 digit, calculated)

	tipoDoc := LeftPad(strconv.Itoa(int(params.TipoDocumento)), '0', 2)
	rucPadded := LeftPad(params.RUC, '0', 8)
	dv := LeftPad(params.DigitoVerificador, '0', 1)
	fecha := params.FechaEmision.Format("20060102")
	tipoContrib := strconv.Itoa(int(params.TipoContribuyente))
	tipoEmision := strconv.Itoa(int(params.TipoEmision))

	// Build CDC without verify digit
	cdcWithoutDV := tipoDoc +
		rucPadded +
		dv +
		params.Establecimiento +
		params.PuntoExpedicion +
		params.NumeroDocumento +
		tipoContrib +
		tipoEmision +
		fecha +
		params.CodigoSeguridad

	// Calculate verify digit
	verifyDigit := CalculateCDCVerifyDigit(cdcWithoutDV)

	return cdcWithoutDV + strconv.Itoa(verifyDigit), nil
}

// CalculateCDCVerifyDigit calculates the verify digit for CDC using Module 11
func CalculateCDCVerifyDigit(cdc string) int {
	// Implementation of Module 11 algorithm for CDC
	// Weights: 2, 3, 4, 5, 6, 7, 2, 3, 4, 5... (repeating)
	weights := []int{2, 3, 4, 5, 6, 7}
	sum := 0
	weightIndex := 0

	// Process from right to left
	for i := len(cdc) - 1; i >= 0; i-- {
		digit, err := strconv.Atoi(string(cdc[i]))
		if err != nil {
			continue // Skip non-numeric characters
		}
		sum += digit * weights[weightIndex%len(weights)]
		weightIndex++
	}

	remainder := sum % 11
	result := 11 - remainder

	if result == 10 {
		return 0
	}
	if result == 11 {
		return 0
	}
	return result
}

// ============================================================================
// RUC Validation
// ============================================================================

// ValidateRUC validates a Paraguayan RUC
func ValidateRUC(ruc string) (bool, error) {
	// RUC format: XXXXXXXX-V where V is the verify digit
	parts := strings.Split(ruc, "-")
	if len(parts) != 2 {
		return false, fmt.Errorf("RUC must contain hyphen separator")
	}

	baseRUC := parts[0]
	providedDV := parts[1]

	// Validate base RUC is numeric
	if _, err := strconv.Atoi(baseRUC); err != nil {
		return false, fmt.Errorf("RUC base must be numeric")
	}

	// Calculate verify digit
	calculatedDV := CalculateRUCVerifyDigit(baseRUC)

	expectedDV, err := strconv.Atoi(providedDV)
	if err != nil {
		return false, fmt.Errorf("verify digit must be numeric")
	}

	return calculatedDV == expectedDV, nil
}

// CalculateRUCVerifyDigit calculates the verify digit for a RUC using Module 11
func CalculateRUCVerifyDigit(baseRUC string) int {
	// Module 11 with base 11
	weights := []int{2, 3, 4, 5, 6, 7, 8, 9, 10}
	sum := 0
	weightIndex := 0

	// Process from right to left
	for i := len(baseRUC) - 1; i >= 0; i-- {
		digit, err := strconv.Atoi(string(baseRUC[i]))
		if err != nil {
			continue
		}
		if weightIndex < len(weights) {
			sum += digit * weights[weightIndex]
		}
		weightIndex++
	}

	remainder := sum % 11
	result := 11 - remainder

	if result == 10 {
		return 0
	}
	if result == 11 {
		return 0
	}
	return result
}

// SplitRUC splits a RUC into base and verify digit
func SplitRUC(ruc string) (base, dv string, err error) {
	parts := strings.Split(ruc, "-")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid RUC format")
	}
	return parts[0], parts[1], nil
}

// ============================================================================
// Code Generation
// ============================================================================

// GenerateSecurityCode generates a random security code
func GenerateSecurityCode() string {
	now := time.Now()
	return fmt.Sprintf("%09d", now.UnixNano()%1000000000)
}

// ============================================================================
// Totals Calculation
// ============================================================================

// TotalsInput contains item data for totals calculation
type TotalsInput struct {
	PrecioUnitario float64
	Cantidad       float64
	Descuento      float64
	TasaIVA        float64 // 0, 5, or 10
	EsExento       bool
	EsExonerado    bool
}

// TotalsResult contains calculated totals
type TotalsResult struct {
	SubtotalExe    float64
	SubtotalExo    float64
	Subtotal5      float64
	Subtotal10     float64
	TotalBruto     float64
	TotalDescuento float64
	TotalNeto      float64
	BaseGravada5   float64
	BaseGravada10  float64
	IVA5           float64
	IVA10          float64
	TotalIVA       float64
}

// CalculateTotals calculates document totals from items
func CalculateTotals(items []TotalsInput) TotalsResult {
	result := TotalsResult{}

	for _, item := range items {
		bruto := item.PrecioUnitario * item.Cantidad
		neto := bruto - item.Descuento

		result.TotalBruto += bruto
		result.TotalDescuento += item.Descuento
		result.TotalNeto += neto

		if item.EsExento {
			result.SubtotalExe += neto
		} else if item.EsExonerado {
			result.SubtotalExo += neto
		} else {
			switch item.TasaIVA {
			case 5:
				result.Subtotal5 += neto
				result.BaseGravada5 += neto / 1.05
				result.IVA5 += neto - (neto / 1.05)
			case 10:
				result.Subtotal10 += neto
				result.BaseGravada10 += neto / 1.10
				result.IVA10 += neto - (neto / 1.10)
			}
		}
	}

	result.TotalIVA = result.IVA5 + result.IVA10
	return result
}

// ============================================================================
// QR Code Generation
// ============================================================================

// QRParams contains parameters for QR code URL generation
type QRParams struct {
	CDC          string
	NroTimbrado  string
	FechaEmision time.Time
	NroItem      int
	FechaFirma   time.Time
	MontoTotal   float64
	CSCid        string
}

// GenerateQRURL generates the URL for SIFEN QR code
func GenerateQRURL(baseURL string, params QRParams) string {
	// QR format: baseURL?nVersion=150&Id=CDC&dFeEmiDE=fecha...
	fecha := params.FechaEmision.Format("2006-01-02T15:04:05")

	return fmt.Sprintf("%snVersion=150&Id=%s&dFeEmiDE=%s&dRucRec=&dTotGralOpe=%.2f&dTotIVA=0&cItems=%d&DigesValue=&IdCSC=%s",
		baseURL,
		params.CDC,
		fecha,
		params.MontoTotal,
		params.NroItem,
		params.CSCid,
	)
}
