package otp

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/coretrix/hitrix/service/component/generator"
)

const SMSOTPProviderMada = "Mada"

type Mada struct {
	username         string
	password         string
	url              string
	sourceName       string
	otpLength        int
	phonePrefixes    []string
	generatorService generator.IGenerator
}

func NewMadaSMSOTPProvider(
	username string,
	password string,
	url string,
	sourceName string,
	otpLength int,
	phonePrefixes []string,
	generatorService generator.IGenerator,
) *Mada {
	return &Mada{
		username:         username,
		password:         password,
		url:              url,
		sourceName:       sourceName,
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
	v := &soapRQ{
		XMLNsXSI:     "http://www.w3.org/2001/XMLSchema-instance",
		XMLNsXSD:     "http://www.w3.org/2001/XMLSchema",
		XMLNsSoapEnv: "http://schemas.xmlsoap.org/soap/envelope/",
		XMLNsSOS:     "http://www.openmindnetworks.com/SoS",
		Body: &soapBody{
			SOS: &sos{
				SoapEnv: "http://schemas.xmlsoap.org/soap/encoding/",
				SMRequest: &smRequest{
					Source: &source{
						TON:  "5",
						NPI:  "1",
						ADDR: m.sourceName,
					},
					Destination: &destination{
						TON:  "1",
						NPI:  "1",
						ADDR: strings.TrimPrefix(recipientPhoneNumber, "+"),
					},
					ShortMessage: &shortMessage{
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

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s", m.url), bytes.NewBuffer(payload))
	if err != nil {
		return string(payload), "", err
	}

	req.Header.Set("Content-Type", "text/xml")
	req.Header.Set("Authorization", "Basic "+basicAuth(m.username, m.password))

	response, err := client.Do(req)
	if err != nil {
		return string(payload), "", err
	}
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return string(payload), string(bodyBytes), err
	}

	resp := &soapResp{}
	if err := xml.Unmarshal(bodyBytes, resp); err != nil {
		return string(payload), string(bodyBytes), err
	}

	if err := validateSmResponse(resp.Body.SOS.SMResponse.CommandStatus); err != nil {
		return string(payload), string(bodyBytes), err
	}

	return string(payload), string(bodyBytes), nil
}

func basicAuth(username, password string) string {
	return base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
}

type soapRQ struct {
	XMLName      xml.Name `xml:"soapenv:Envelope"`
	XMLNsXSI     string   `xml:"xmlns:xsi,attr"`
	XMLNsXSD     string   `xml:"xmlns:xsd,attr"`
	XMLNsSoapEnv string   `xml:"xmlns:soapenv,attr"`
	XMLNsSOS     string   `xml:"xmlns:sos,attr"`
	Header       *soapHeader
	Body         *soapBody
}

type soapHeader struct {
	XMLName xml.Name `xml:"soapenv:Header"`
}

type soapBody struct {
	XMLName xml.Name `xml:"soapenv:Body"`
	SOS     *sos
}

type sos struct {
	XMLName   xml.Name `xml:"sos:SubmitSM"`
	SoapEnv   string   `xml:"soapenv:encodingStyle,attr"`
	SMRequest *smRequest
}

type smRequest struct {
	XMLName            xml.Name `xml:"smRequest"`
	Source             *source
	Destination        *destination
	ShortMessage       *shortMessage
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

type soapResp struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    *soapRespBody
}

type soapRespBody struct {
	XMLName xml.Name `xml:"Body"`
	SOS     *sosResp
}

type sosResp struct {
	XMLName    xml.Name `xml:"SubmitSMResponse"`
	SMResponse *smResponse
}

type smResponse struct {
	XMLName       xml.Name `xml:"smResponse"`
	CommandStatus string   `xml:"commandStatus"`
}

func validateSmResponse(commandStatus string) error {
	if commandStatus == "0" {
		return nil
	}
	if errorMsg, ok := errorCodes[commandStatus]; ok {
		return fmt.Errorf("otp Mada[%s]: %s", commandStatus, errorMsg)
	}

	return fmt.Errorf("otp Mada[%s]: expected commandStatus code 0, but got (unrecognized) %s", commandStatus, commandStatus)
}

// No online reference
// Received this from technical team MyMada via email
var errorCodes = map[string]string{
	"0":    "No Error",
	"1":    "Message too long",
	"2":    "Command length is invalid",
	"3":    "Command ID is invalid or not supported",
	"4":    "Incorrect bind status for given command",
	"5":    "Already bound",
	"6":    "Invalid Priority Flag",
	"7":    "Invalid registered delivery flag",
	"8":    "System error",
	"10":   "Invalid source address",
	"11":   "Invalid destination address",
	"12":   "Message ID is invalid",
	"13":   "Bind failed",
	"14":   "Invalid password",
	"15":   "Invalid System ID",
	"17":   "Canceling message failed",
	"19":   "Message recplacement failed",
	"20":   "Message queue full",
	"21":   "Invalid service type",
	"51":   "Invalid number of destinations",
	"52":   "Invalid distribution list name",
	"64":   "Invalid destination flag",
	"66":   "Invalid submit with replace request",
	"67":   "Invalid esm class set",
	"68":   "Invalid submit to ditribution list",
	"69":   "Submitting message has failed",
	"72":   "Invalid source address type of number ( TON )",
	"73":   "Invalid source address numbering plan ( NPI )",
	"80":   "Invalid destination address type of number ( TON )",
	"81":   "Invalid destination address numbering plan ( NPI )",
	"83":   "Invalid system type",
	"84":   "Invalid replace_if_present flag",
	"85":   "Invalid number of messages",
	"88":   "Throttling error",
	"97":   "Invalid scheduled delivery time",
	"98":   "Invalid Validty Period value",
	"99":   "Predefined message not found",
	"100":  "ESME Receiver temporary error",
	"101":  "ESME Receiver permanent error",
	"102":  "ESME Receiver reject message error",
	"103":  "Message query request failed",
	"192":  "Error in the optional part of the PDU body",
	"193":  "TLV not allowed",
	"194":  "Invalid parameter length",
	"195":  "Expected TLV missing",
	"196":  "Invalid TLV value",
	"254":  "Transaction delivery failure",
	"255":  "Unknown error",
	"256":  "ESME not authorized to use specified servicetype",
	"257":  "ESME prohibited from using specified operation",
	"258":  "Specified servicetype is unavailable",
	"259":  "Specified servicetype is denied",
	"260":  "Invalid data coding scheme",
	"261":  "Invalid source address subunit",
	"262":  "Invalid destination address subunit",
	"1035": "Insufficient credits to send message",
	"1036": "Destination address blocked by the ActiveXperts SMPP Demo Server",
}
