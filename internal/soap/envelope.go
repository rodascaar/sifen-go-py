package soap

import "encoding/xml"

type Envelope struct {
	XMLName xml.Name `xml:"http://www.w3.org/2003/05/soap-envelope Envelope"`
	Header  *Header  `xml:"Header,omitempty"`
	Body    Body
}

type Header struct {
	XMLName xml.Name `xml:"Header"`
	// Add header content if needed
}

type Body struct {
	XMLName xml.Name `xml:"Body"`
	Content interface{}
}

func NewEnvelope(content interface{}) *Envelope {
	return &Envelope{
		Body: Body{
			Content: content,
		},
	}
}
