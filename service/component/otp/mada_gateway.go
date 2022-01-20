package otp

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/coretrix/hitrix/service/component/generator"
)

const SMSOTPProviderMada = "Mada"

type Mada struct {
	username         string
	password         string
	url              string
	otpLength        int
	phonePrefixes    []string
	generatorService generator.IGenerator
}

func NewMadaSMSOTPProvider(username, password, url string, otpLength int, phonePrefixes []string, generatorService generator.IGenerator) *Mada {
	return &Mada{
		username:         username,
		password:         password,
		url:              url,
		otpLength:        otpLength,
		phonePrefixes:    phonePrefixes,
		generatorService: generatorService,
	}
}

func (m *Mada) GetName() string {
	return SMSOTPProviderMada
}

func (m *Mada) GetCode() string {
	var code int64
	if m.otpLength == 0 {
		code = m.generatorService.GenerateRandomRangeNumber(10000, 99999)
	} else {
		min := int64(math.Pow(10, float64(m.otpLength-1)))
		max := int64(math.Pow(10, float64(m.otpLength))) - 1
		code = m.generatorService.GenerateRandomRangeNumber(min, max)
	}

	return strconv.FormatInt(code, 10)
}

func (m *Mada) GetPhonePrefixes() []string {
	return m.phonePrefixes
}

func (m *Mada) SendOTP(phone *Phone, code string) (string, string, error) {
	return m.soapCall(phone.Number, code)
}

func (m *Mada) Call(_ *Phone, _ string, _ string) (string, string, error) {
	// not implemented
	return "", "", nil
}

func (m *Mada) VerifyOTP(_ *Phone, code, generatedCode string) (string, string, bool, bool, error) {
	return "", "", true, code == generatedCode, nil
}

func (m *Mada) soapCall(recipientPhoneNumber, otp string) (string, string, error) {
	v := soapRQ{
		XMLNsXSI:     "http://www.w3.org/2001/XMLSchema-instance",
		XMLNsXSD:     "http://www.w3.org/2001/XMLSchema",
		XMLNsSoapEnv: "http://schemas.xmlsoap.org/soap/envelope/",
		Body: soapBody{
			SOS: sos{
				SoapEnv: "http://schemas.xmlsoap.org/soap/encoding/",
				SMRequest: smRequest{
					Source: source{
						TON:  "5",
						NPI:  "1",
						ADDR: "Riverstream",
					},
					Destination: destination{
						TON:  "1",
						NPI:  "1",
						ADDR: recipientPhoneNumber,
					},
					ShortMessage: shortMessage{
						StringData: otp,
					},
					RegisteredDelivery: "1",
				},
			},
		},
	}

	payload, err := xml.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", "", err
	}

	timeout := 30 * time.Second
	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s:%s@%s", m.username, m.password, m.url), bytes.NewBuffer(payload))
	if err != nil {
		return string(payload), "", err
	}

	req.Header.Set("Accept", "text/xml, multipart/related")
	req.Header.Set("SOAPAction", "http://schemas.xmlsoap.org/soap/envelope/")
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")

	response, err := client.Do(req)
	if err != nil {
		return string(payload), "", err
	}
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return string(payload), string(bodyBytes), err
	}

	return string(payload), string(bodyBytes), nil
}

type soapRQ struct {
	XMLName      xml.Name `xml:"soapenv:Envelope"`
	XMLNsXSI     string   `xml:"xmlns:xsi,attr"`
	XMLNsXSD     string   `xml:"xmlns:xsd,attr"`
	XMLNsSoapEnv string   `xml:"xmlns:soapenv,attr"`
	Header       soapHeader
	Body         soapBody
}

type soapHeader struct {
	XMLName xml.Name `xml:"soapenv:Header"`
}

type soapBody struct {
	XMLName xml.Name `xml:"soapenv:Body"`
	SOS     sos
}

type sos struct {
	XMLName   xml.Name `xml:"sos:SubmitSM"`
	SoapEnv   string   `xml:"soapenv:encodingStyle,attr"`
	SMRequest smRequest
}

type smRequest struct {
	XMLName            xml.Name `xml:"smRequest"`
	Source             source
	Destination        destination
	ShortMessage       shortMessage
	RegisteredDelivery string `xml:"registeredDelivery"`
}

type source struct {
	XMLName xml.Name `xml:"source"`
	TON     string   `xml:"ton"`
	NPI     string   `xml:"npi"`
	ADDR    string   `xml:"addr"`
}

type destination struct {
	XMLName xml.Name `xml:"destination"`
	TON     string   `xml:"ton"`
	NPI     string   `xml:"npi"`
	ADDR    string   `xml:"addr"`
}

type shortMessage struct {
	XMLName    xml.Name `xml:"shortMessage"`
	StringData string   `xml:"stringData"`
}
