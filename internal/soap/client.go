package soap

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/pem"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/pkcs12"
)

// ClientConfig holds the necessary configuration for the SOAP client
type ClientConfig struct {
	TimeoutMs          int
	UseClientCert      bool
	ClientCertPath     string
	ClientCertPassword string
	UserAgent          string
}

type Client struct {
	httpClient *http.Client
	config     *ClientConfig
	cert       tls.Certificate
}

func (c *Client) GetCertificate() tls.Certificate {
	return c.cert
}

func NewClient(cfg *ClientConfig) (*Client, error) {
	client := &http.Client{
		Timeout: time.Duration(cfg.TimeoutMs) * time.Millisecond,
	}

	var parsedCert tls.Certificate

	if cfg.UseClientCert {
		certData, err := loadCertificate(cfg.ClientCertPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load certificate: %w", err)
		}

		// Use pkcs12.ToPEM to get blocks, then parse X509KeyPair
		blocks, err := pkcs12.ToPEM(certData, cfg.ClientCertPassword)
		if err != nil {
			return nil, fmt.Errorf("failed to decode pkcs12: %w", err)
		}

		var pemData []byte
		for _, b := range blocks {
			pemData = append(pemData, pem.EncodeToMemory(b)...)
		}

		parsedCert, err = tls.X509KeyPair(pemData, pemData)
		if err != nil {
			return nil, fmt.Errorf("failed to create x509 key pair: %w", err)
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{parsedCert},
			MinVersion:   tls.VersionTLS12,
			MaxVersion:   tls.VersionTLS12,
		}

		transport := &http.Transport{
			TLSClientConfig: tlsConfig,
		}
		client.Transport = transport
	}

	return &Client{
		httpClient: client,
		config:     cfg,
		cert:       parsedCert,
	}, nil
}

func loadCertificate(pathOrBase64 string) ([]byte, error) {
	// Try reading as file
	if _, err := os.Stat(pathOrBase64); err == nil {
		return ioutil.ReadFile(pathOrBase64)
	}

	// Try decoding as base64
	return base64.StdEncoding.DecodeString(pathOrBase64)
}

func (c *Client) Send(url string, payload interface{}) ([]byte, error) {
	// Construct SOAP Envelope
	envelope := NewEnvelope(payload)

	reqBody, err := xml.Marshal(envelope)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal soap envelope: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/xml; charset=utf-8")
	userAgent := c.config.UserAgent
	if userAgent == "" {
		userAgent = "rshk-jsifenlib-go"
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return body, fmt.Errorf("sifen returned status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
