package signature

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/xml"
	"fmt"

	"github.com/beevik/etree"
	"github.com/ucarion/c14n"
)

// Signer handles XML Digital Signature logic
type Signer struct {
	Cert tls.Certificate
}

func NewSigner(cert tls.Certificate) *Signer {
	return &Signer{Cert: cert}
}

func (s *Signer) Sign(xmlBytes []byte, elementID string) ([]byte, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(xmlBytes); err != nil {
		return nil, fmt.Errorf("failed to parse xml: %w", err)
	}

	// 1. Find the element to sign
	// XPath is tricky, we'll traverse.
	// We assume ID attribute is named "Id" (case sensitive for Sifen).
	elem := doc.FindElement(fmt.Sprintf("//*[@Id='%s']", elementID))
	if elem == nil {
		return nil, fmt.Errorf("element with Id='%s' not found", elementID)
	}

	// 2. Canonicalize the element (Exclusive C14N)
	docTemp := etree.NewDocument()
	docTemp.SetRoot(elem.Copy())
	xmlBytesTemp, _ := docTemp.WriteToBytes()

	decoder := xml.NewDecoder(bytes.NewReader(xmlBytesTemp))
	canonicalBytes, err := c14n.Canonicalize(decoder)
	if err != nil {
		return nil, fmt.Errorf("c14n failed: %w", err)
	}

	// 3. Calculate Digest
	hasher := sha256.New()
	hasher.Write(canonicalBytes)
	digest := hasher.Sum(nil)
	digestBase64 := base64.StdEncoding.EncodeToString(digest)

	// 4. Construct SignedInfo
	signedInfo := etree.NewElement("SignedInfo")
	signedInfo.CreateAttr("xmlns", "http://www.w3.org/2000/09/xmldsig#")

	c14nMethod := signedInfo.CreateElement("CanonicalizationMethod")
	c14nMethod.CreateAttr("Algorithm", "http://www.w3.org/2001/10/xml-exc-c14n#")

	sigMethod := signedInfo.CreateElement("SignatureMethod")
	sigMethod.CreateAttr("Algorithm", "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256")

	ref := signedInfo.CreateElement("Reference")
	ref.CreateAttr("URI", "#"+elementID)

	transforms := ref.CreateElement("Transforms")
	trans1 := transforms.CreateElement("Transform")
	trans1.CreateAttr("Algorithm", "http://www.w3.org/2000/09/xmldsig#enveloped-signature")
	trans2 := transforms.CreateElement("Transform")
	trans2.CreateAttr("Algorithm", "http://www.w3.org/2001/10/xml-exc-c14n#")

	digestMethod := ref.CreateElement("DigestMethod")
	digestMethod.CreateAttr("Algorithm", "http://www.w3.org/2001/04/xmlenc#sha256")

	digestVal := ref.CreateElement("DigestValue")
	digestVal.SetText(digestBase64)

	// 5. Canonicalize SignedInfo
	docSignedInfo := etree.NewDocument()
	docSignedInfo.SetRoot(signedInfo.Copy())
	siBytes, _ := docSignedInfo.WriteToBytes()

	siDecoder := xml.NewDecoder(bytes.NewReader(siBytes))
	canonicalSiBytes, err := c14n.Canonicalize(siDecoder)
	if err != nil {
		return nil, fmt.Errorf("c14n SignedInfo failed: %w", err)
	}

	// 6. Sign SignedInfo digest
	// We need RSA Private Key
	privKey, ok := s.Cert.PrivateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not RSA")
	}

	// Hash the canonical SignedInfo
	siHasher := sha256.New()
	siHasher.Write(canonicalSiBytes)
	siDigest := siHasher.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, siDigest)
	if err != nil {
		return nil, fmt.Errorf("rsa sign failed: %w", err)
	}
	signatureBase64 := base64.StdEncoding.EncodeToString(signature)

	// 7. Construct Signature Element
	signatureElem := etree.NewElement("Signature")
	signatureElem.CreateAttr("xmlns", "http://www.w3.org/2000/09/xmldsig#")

	signatureElem.AddChild(signedInfo) // The one we computed digest for

	sigValElem := signatureElem.CreateElement("SignatureValue")
	sigValElem.SetText(signatureBase64)

	keyInfoElem := signatureElem.CreateElement("KeyInfo")
	x509DataElem := keyInfoElem.CreateElement("X509Data")
	x509CertElem := x509DataElem.CreateElement("X509Certificate")

	// Get first leaf cert
	leafCert := s.Cert.Certificate[0]
	x509CertElem.SetText(base64.StdEncoding.EncodeToString(leafCert))

	// 8. Append Signature to Root (or doc)
	elem.AddChild(signatureElem)

	// Return the full XML
	return doc.WriteToBytes()
}
