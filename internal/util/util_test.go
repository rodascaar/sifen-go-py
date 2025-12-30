package util

import (
	"testing"
	"time"
)

func TestLeftPad(t *testing.T) {
	tests := []struct {
		input    string
		padChar  rune
		length   int
		expected string
	}{
		{"123", '0', 5, "00123"},
		{"12345", '0', 5, "12345"},
		{"123456", '0', 5, "123456"},
		{"", '0', 3, "000"},
		{"1", '0', 1, "1"},
	}

	for _, tt := range tests {
		result := LeftPad(tt.input, tt.padChar, tt.length)
		if result != tt.expected {
			t.Errorf("LeftPad(%q, %q, %d) = %q; want %q",
				tt.input, string(tt.padChar), tt.length, result, tt.expected)
		}
	}
}

func TestRightPad(t *testing.T) {
	tests := []struct {
		input    string
		padChar  rune
		length   int
		expected string
	}{
		{"123", '0', 5, "12300"},
		{"12345", '0', 5, "12345"},
		{"", '0', 3, "000"},
	}

	for _, tt := range tests {
		result := RightPad(tt.input, tt.padChar, tt.length)
		if result != tt.expected {
			t.Errorf("RightPad(%q, %q, %d) = %q; want %q",
				tt.input, string(tt.padChar), tt.length, result, tt.expected)
		}
	}
}

func TestCalculateRUCVerifyDigit(t *testing.T) {
	tests := []struct {
		baseRUC  string
		expected int
	}{
		{"80069563", 1},
		{"50062360", 0},
	}

	for _, tt := range tests {
		result := CalculateRUCVerifyDigit(tt.baseRUC)
		if result != tt.expected {
			t.Errorf("CalculateRUCVerifyDigit(%q) = %d; want %d",
				tt.baseRUC, result, tt.expected)
		}
	}
}

func TestValidateRUC(t *testing.T) {
	tests := []struct {
		ruc      string
		valid    bool
		hasError bool
	}{
		{"80069563-1", true, false},
		{"invalid", false, true},
		{"12345-X", false, true},
	}

	for _, tt := range tests {
		valid, err := ValidateRUC(tt.ruc)
		if tt.hasError && err == nil {
			t.Errorf("ValidateRUC(%q) expected error but got none", tt.ruc)
		}
		if !tt.hasError && valid != tt.valid {
			t.Errorf("ValidateRUC(%q) = %v; want %v", tt.ruc, valid, tt.valid)
		}
	}
}

func TestSplitRUC(t *testing.T) {
	tests := []struct {
		ruc      string
		base     string
		dv       string
		hasError bool
	}{
		{"80069563-1", "80069563", "1", false},
		{"invalid", "", "", true},
	}

	for _, tt := range tests {
		base, dv, err := SplitRUC(tt.ruc)
		if tt.hasError && err == nil {
			t.Errorf("SplitRUC(%q) expected error but got none", tt.ruc)
			continue
		}
		if !tt.hasError {
			if base != tt.base || dv != tt.dv {
				t.Errorf("SplitRUC(%q) = (%q, %q); want (%q, %q)",
					tt.ruc, base, dv, tt.base, tt.dv)
			}
		}
	}
}

func TestGenerateCDC(t *testing.T) {
	params := CDCParams{
		TipoDocumento:     1,
		RUC:               "80069563",
		DigitoVerificador: "1",
		Establecimiento:   "001",
		PuntoExpedicion:   "001",
		NumeroDocumento:   "0000001",
		TipoContribuyente: 2,
		FechaEmision:      time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		TipoEmision:       1,
		CodigoSeguridad:   "123456789",
	}

	cdc, err := GenerateCDC(params)
	if err != nil {
		t.Errorf("GenerateCDC() error = %v", err)
		return
	}

	if len(cdc) != 44 {
		t.Errorf("GenerateCDC() CDC length = %d; want 44", len(cdc))
	}
}

func TestGenerateCDCValidation(t *testing.T) {
	// Test missing RUC
	params := CDCParams{
		TipoDocumento: 1,
		RUC:           "",
	}
	_, err := GenerateCDC(params)
	if err == nil {
		t.Error("GenerateCDC() with empty RUC should return error")
	}

	// Test invalid establecimiento
	params.RUC = "80069563"
	params.Establecimiento = "01" // Too short
	_, err = GenerateCDC(params)
	if err == nil {
		t.Error("GenerateCDC() with invalid establecimiento should return error")
	}
}

func TestCalculateTotals(t *testing.T) {
	items := []TotalsInput{
		{PrecioUnitario: 100000, Cantidad: 2, Descuento: 0, TasaIVA: 10},
		{PrecioUnitario: 50000, Cantidad: 1, Descuento: 5000, TasaIVA: 5},
		{PrecioUnitario: 25000, Cantidad: 1, Descuento: 0, EsExento: true},
	}

	result := CalculateTotals(items)

	// Total bruto should be 100000*2 + 50000*1 + 25000*1 = 275000
	expectedBruto := 275000.0
	if result.TotalBruto != expectedBruto {
		t.Errorf("TotalBruto = %.2f; want %.2f", result.TotalBruto, expectedBruto)
	}

	// Total descuento should be 5000
	if result.TotalDescuento != 5000.0 {
		t.Errorf("TotalDescuento = %.2f; want 5000", result.TotalDescuento)
	}

	// Exento should be 25000
	if result.SubtotalExe != 25000.0 {
		t.Errorf("SubtotalExe = %.2f; want 25000", result.SubtotalExe)
	}
}

func TestGenerateSecurityCode(t *testing.T) {
	code := GenerateSecurityCode()
	if len(code) != 9 {
		t.Errorf("GenerateSecurityCode() length = %d; want 9", len(code))
	}
}

func TestGenerateQRURL(t *testing.T) {
	params := QRParams{
		CDC:          "01800695631001001000000612021112917595714694",
		NroTimbrado:  "12345678",
		FechaEmision: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		NroItem:      5,
		MontoTotal:   150000.00,
		CSCid:        "0001",
	}

	baseURL := "https://ekuatia.set.gov.py/consultas/qr?"
	url := GenerateQRURL(baseURL, params)

	if url == "" {
		t.Error("GenerateQRURL() returned empty string")
	}

	// Basic validation that URL contains expected components
	if len(url) < 50 {
		t.Error("GenerateQRURL() URL seems too short")
	}
}
