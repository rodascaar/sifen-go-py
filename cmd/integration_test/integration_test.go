package integration_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/rodascaar/sifen-go-py/sifen"
	// "github.com/rodascaar/sifen-go-py/sifen/events" // Explicitly used in TestEventos if manually building, or implicitly via client methods
	"github.com/rodascaar/sifen-go-py/sifen/models"
	"github.com/rodascaar/sifen-go-py/sifen/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	CERT_FILE     = "../../certificate.p12" // Adjust path as needed
	CERT_PASS     = "password"              // Replace with actual password or env var
	RUC_EMISOR    = "80000001"
	TIMBRADO_TEST = 12560156
	ID_CSC        = "0001"
	CSC           = "ABCD0000000000000000000000000000"
)

var client *sifen.SifenClient

func setup() error {
	// Check if certificate exists, otherwise skip
	if _, err := os.Stat(CERT_FILE); os.IsNotExist(err) {
		return fmt.Errorf("certificate file not found at %s", CERT_FILE)
	}

	config := sifen.NewSifenConfig()
	config.Ambiente = sifen.TipoAmbienteDev
	config.UsarCertificadoCliente = true
	config.CertificadoCliente = CERT_FILE
	config.ContrasenaCertificadoCliente = CERT_PASS
	config.IdCSC = ID_CSC
	config.CSC = CSC
	config.HttpConnectTimeout = 30000 // 30s
	config.HttpReadTimeout = 30000    // 30s

	var err error
	client, err = sifen.NewSifenClient(config)
	return err
}

func TestMain(m *testing.M) {
	// Don't enforce setup here to allow individual tests to skip if setup fails
	code := m.Run()
	if client != nil {
		client.Close()
	}
	os.Exit(code)
}

func getClient(t *testing.T) *sifen.SifenClient {
	if client == nil {
		err := setup()
		if err != nil {
			t.Skipf("Skipping integration test: %v (Please upload %s to run tests)", err, CERT_FILE)
		}
	}
	return client
}

func createDummyDE(id string) *models.DocumentoElectronico {
	de := models.NewDE(id)

	// Basic Fields for validation
	de.DE.GOpeDE = models.TgOpeDE{
		ITipEmi:    types.TTipEmi_Normal,
		DDesTipEmi: "Normal",
		DCodSeg:    "12345", // Random security code
		DInfoEmi:   "Factura de prueba",
	}

	de.DE.GTimb = models.TgTimb{
		ITiDE:    types.TTiDE_FacturaElectronica,
		DDesTiDE: "Factura Electronica",
		DNumTim:  TIMBRADO_TEST,
		DEst:     "001",
		DPunExp:  "001",
		DNumDoc:  "0000001", // Should be somewhat dynamic in real tests
		DFeIniT:  "2024-01-01",
	}

	de.DE.GDatGralOpe = models.TdDatGralOpe{
		DFeEmiDE: time.Now().Format("2006-01-02T15:04:05"),
		GEmis: models.TgEmis{
			DRucEm:     RUC_EMISOR,
			DDVEmi:     "7",
			ITipCont:   types.TiTipCont_PersonaJuridica, // Fixed Type
			DNomEmi:    "Empresa de Prueba S.A.",
			CDepEmi:    types.TDepartamento_Capital,
			DDesDepEmi: "CAPITAL",
			CDisEmi:    1,
			DDesDisEmi: "ASUNCION",
			CCiuEmi:    1,
			DDesCiuEmi: "ASUNCION",
			DDirEmi:    "Calle Falsa 123",
			DNumCas:    "123", // Fixed Type (string)
			DTelEmi:    "021123456",
			DEmailE:    "test@empresa.com.py",
			GActEcoList: []models.TgActEco{
				{CActEco: "47119", DDesActEco: "COMERCIO AL POR MENOR"},
			},
		},
		GDatRec: models.TgDatRec{
			INatRec:    types.TiNatRec_NoContribuyente, // Fixed Type
			ITiOpe:     types.TiTiOpe_B2C,              // Fixed Type
			CPaisRec:   types.PaisType_PRY,             // Fixed Type
			DDesPaisRe: "Paraguay",
			ITipIDRec:  new(int16), // Optional, pointer
			// Use DNumIDRec directly if appropriate, or adjust based on specific logic
			DNumIDRec: "1234567",
			DNomRec:   "Consumidor Final",
		},
	}
	// Initialize pointer values
	*de.DE.GDatGralOpe.GDatRec.ITipIDRec = int16(types.TTipDocRec_CedulaParaguaya) // Fixed Type and Cast

	return de
}

func TestRecepcionDE_Sync(t *testing.T) {
	c := getClient(t)

	// Create a unique ID for this run (mocking CDC structure or random ID for now)
	id := fmt.Sprintf("01800%035d", time.Now().UnixNano())
	de := createDummyDE(id)

	resp, err := c.RecepcionDE(de)
	require.NoError(t, err)
	assert.NotNil(t, resp)

	t.Logf("Response Status: %s", resp.DCodRes)
	t.Logf("Response Msg: %s", resp.DMsgRes)

	// In test environment without valid RUC/Timbrado, we expect rejection but valid HTTP/SOAP response
	if resp.DCodRes == "0500" {
		t.Log("Server Error (expected if service is flaky or bad request structure)")
	}
}

func TestRecepcionLoteDE_Async(t *testing.T) {
	c := getClient(t)

	var docs []*models.DocumentoElectronico
	for i := 0; i < 3; i++ {
		id := fmt.Sprintf("01800%035d", time.Now().UnixNano()+int64(i))
		docs = append(docs, createDummyDE(id))
	}

	resp, err := c.RecepcionLoteDE(docs)
	require.NoError(t, err)
	assert.NotNil(t, resp)

	t.Logf("Batch Response Status: %s", resp.DCodRes)
	t.Logf("Ticket: %s", resp.DProtConsLot) // Fixed Field Name
}

func TestConsultas(t *testing.T) {
	c := getClient(t)

	t.Run("ConsultaRUC", func(t *testing.T) {
		resp, err := c.ConsultaRUC(RUC_EMISOR)
		require.NoError(t, err)
		t.Logf("RUC Status: %s - %s", resp.DCodRes, resp.DMsgRes)
	})

	t.Run("ConsultaDE", func(t *testing.T) {
		// Needs a valid CDC. We can test with a dummy one and expect "Not Found"
		cdc := "01800000000000000000000000000000000000000000"
		resp, err := c.ConsultaDE(cdc)
		require.NoError(t, err)
		t.Logf("DE Status: %s - %s", resp.DCodRes, resp.DMsgRes)
	})
}

func TestEventos(t *testing.T) {
	c := getClient(t)

	// Test Cancellation Event
	// Requires a valid CDC to cancel. In this dummy test, we use a made-up CDC.
	// We expect a rejection due to invalid CDC or "DE not found", but the SOAP call should succeed.
	cdc := "01800000000000000000000000000000000000000000"
	motivo := "Prueba de cancelación desde integración"

	resp, err := c.CancelarDE(cdc, motivo)
	require.NoError(t, err)
	assert.NotNil(t, resp)

	t.Logf("Evaluated Event: %s", resp.DCodRes)
	t.Logf("Message: %s", resp.DMsgRes)
}
