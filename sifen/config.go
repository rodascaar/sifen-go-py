package sifen

import (
	"fmt"

	"github.com/rodascaar/sifen-go-py/internal/util"
	"github.com/rodascaar/sifen-go-py/sifen/cache"
)

type TipoAmbiente string

const (
	TipoAmbienteDev  TipoAmbiente = "DEV"
	TipoAmbienteProd TipoAmbiente = "PROD"
)

type TipoCertificadoCliente string

const (
	TipoCertificadoClientePFX TipoCertificadoCliente = "PFX"
)

const (
	// Default URLs
	URL_BASE_DEV         = "https://sifen-test.set.gov.py"
	URL_BASE_PROD        = "https://sifen.set.gov.py"
	URL_CONSULTA_QR_DEV  = "https://ekuatia.set.gov.py/consultas-test/qr?"
	URL_CONSULTA_QR_PROD = "https://ekuatia.set.gov.py/consultas/qr?"

	SDK_CURRENT_VERSION = "0.2.4" // TODO: Update version logic
)

type SifenConfig struct {
	Ambiente      TipoAmbiente
	UrlBase       string
	UrlBaseLocal  string
	UrlConsultaQr string

	PathRecibe       string
	PathRecibeLote   string
	PathEvento       string
	PathConsultaLote string
	PathConsultaRUC  string
	PathConsulta     string

	UsarCertificadoCliente       bool
	TipoCertificadoCliente       TipoCertificadoCliente
	CertificadoCliente           string
	ContrasenaCertificadoCliente string

	IdCSC string
	CSC   string

	HttpConnectTimeout int // Milliseconds
	HttpReadTimeout    int // Milliseconds
	UserAgent          string

	// Configuración de Caché
	CacheConfig cache.CacheConfig
}

func NewSifenConfig() *SifenConfig {
	cfg := &SifenConfig{
		Ambiente:      TipoAmbienteDev,
		UrlBaseLocal:  URL_BASE_DEV,
		UrlConsultaQr: URL_CONSULTA_QR_DEV,

		PathRecibe:       "/de/ws/sync/recibe.wsdl",
		PathRecibeLote:   "/de/ws/async/recibe-lote.wsdl",
		PathEvento:       "/de/ws/eventos/evento.wsdl",
		PathConsultaLote: "/de/ws/consultas/consulta-lote.wsdl",
		PathConsultaRUC:  "/de/ws/consultas/consulta-ruc.wsdl",
		PathConsulta:     "/de/ws/consultas/consulta.wsdl",

		UsarCertificadoCliente: true,
		IdCSC:                  "0002",
		CSC:                    "EFGH0000000000000000000000000000",

		HttpConnectTimeout: 15 * 1000,
		HttpReadTimeout:    45 * 1000,
		UserAgent:          "rshk-jsifenlib/" + SDK_CURRENT_VERSION + " (GoPort)",

		CacheConfig: cache.DefaultCacheConfig(),
	}
	return cfg
}

func (c *SifenConfig) SetAmbiente(env TipoAmbiente) {
	c.Ambiente = env
	switch env {
	case TipoAmbienteDev:
		c.UrlBaseLocal = URL_BASE_DEV
		c.UrlConsultaQr = URL_CONSULTA_QR_DEV
	case TipoAmbienteProd:
		c.UrlBaseLocal = URL_BASE_PROD
		c.UrlConsultaQr = URL_CONSULTA_QR_PROD
	}
}

func (c *SifenConfig) SetIdCSC(id string) {
	c.IdCSC = util.LeftPad(id, '0', 4)
}

func (c *SifenConfig) String() string {
	return fmt.Sprintf("SifenConfig{Ambiente=%s, UrlBase=%s, ...}", c.Ambiente, c.UrlBase)
}
