package main

import (
	"app-connector/config"
	"app-connector/logger"
	"encoding/xml"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/go-resty/resty/v2"
)

// Response for get data API
type Response struct {
	XMLName  xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	SoapBody SOAPBodyResponse
}

type SOAPBodyResponse struct {
	XMLName xml.Name `xml:"Body"`
	Resp    ResponseBody
	// FaultDetails *Fault
}

type ResponseBody struct {
	XMLName    xml.Name     `xml:"GetSRxDataResponse"`
	ResultData []ResultData `xml:"GetSRxDataResult>SRxResultData"`
}

type ResultData struct {
	XMLName xml.Name `xml:"SRxResultData"`
	TagData TagData  `xml:"Tag"`
	Data    Data     `xml:"Data"`
	State   State    `xml:"State"`
}

type TagData struct {
	XMLName         xml.Name `xml:"Tag"`
	Name            string   `xml:"Name"`
	Archive         string   `xml:"Archive"`
	Unit            string   `xml:"Unit"`
	FormulaOnTheFly string   `xml:"FormulaOnTheFly"`
}

type State struct {
	XMLName      xml.Name `xml:"State"`
	IsValid      string   `xml:"IsValid"`
	ErrorMessage string   `xml:"ErrorMessage"`
	ErrorCode    string   `xml:"ErrorCode"`
}

type Data struct {
	XMLName      xml.Name `xml:"Data"`
	TimeDataItem []TimeDataItem
}

type TimeDataItem struct {
	XMLName xml.Name `xml:"TimeDataItem"`
	Value   string   `xml:"Value"`
	Time    string   `xml:"Time"`
}

// Request for get data API
type SOAPEnvelope struct {
	XMLName xml.Name `xml:"soap:Envelope"`
	Xsi     string   `xml:"xmlns:xsi,attr"`
	Xsd     string   `xml:"xmlns:xsd,attr"`
	Soap    string   `xml:"xmlns:soap,attr"`
	// Header  *SOAPHeader `xml:",omitempty"`
	Body SOAPBody
}

type SOAPBody struct {
	XMLName xml.Name `xml:"soap:Body"`
	// Fault   *SOAPFault  `xml:",omitempty"`
	// Content interface{} `xml:",omitempty"`
	GetSRxData GetSRxData `xml:"http://www.SR-Suite.com/SRxServerWebService GetSRxData"`
}

type GetSRxData struct {
	Tags                   Tags   `xml:"tags"`
	ReferenceDate          string `xml:"referenceDate"`
	TimeUnitString         string `xml:"timeUnitString"`
	Count                  string `xml:"count"`
	IntervalTimeUnitString string `xml:"intervalTimeUnitString"`
	IntervalTimeUnitCount  string `xml:"intervalTimeUnitCount"`
	SrxTimeKind            string `xml:"srxTimeKind"`
	Xmlns                  string `xml:"_xmlns"`
}

type Tags struct {
	SRxTag []SRxTag `xml:"SRxTag"`
}

type SRxTag struct {
	Name string `xml:"Name"`
	// Index           string `xml:"Index"`
	// Archive         string `xml:"Archive"`
	// Unit            string `xml:"Unit"`
	// FormulaOnTheFly string `xml:"FormulaOnTheFly"`
}

var client *resty.Client

func init() {
	config.ReadConfig()
	logger.InitLog()
	InitClient()
}

func main() {
	ConnectAPI()
	// respBody, err := GetDataAPI()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// tagData, err := ParseXML(respBody)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Printf("Data: %+v\n", tagData)

}

func InitClient() {
	// Create a Resty Client
	client = resty.New()
}

func ConnectAPI() {
	// Connect API
	resp, err := client.R().
		SetBasicAuth(config.Config.Api.Connect.BasicAuth.Username, config.Config.Api.Connect.BasicAuth.Password).
		SetHeaders(config.Config.Api.Connect.Headers).
		SetBody(config.Config.Api.Connect.Body).
		Post(config.Config.Api.Connect.Url)
	if err != nil {
		logger.Logger.Error("connect API", "error", err)
	}
	logger.Logger.Debug("connect API", slog.Group("response", slog.Int("status", resp.StatusCode()), slog.Duration("response time", resp.Time()), slog.String("response body", resp.String())))
	// fmt.Println("Connect API")
	// fmt.Println("Response Info:")
	// fmt.Println("  Error      :", err)
	// fmt.Println("  Status Code:", resp.StatusCode())
	// fmt.Println("  Status     :", resp.Status())
	// fmt.Println("  Proto      :", resp.Proto())
	// fmt.Println("  Time       :", resp.Time())
	// fmt.Println("  Received At:", resp.ReceivedAt())
	// fmt.Println("  Body       :\n", resp)
	// fmt.Println()
}

func GetDataAPI() ([]byte, error) {
	// Create request body
	refTime := time.Now().AddDate(0, 0, -1).Format(time.RFC3339)
	body := SOAPEnvelope{
		Body: SOAPBody{
			GetSRxData: GetSRxData{
				ReferenceDate:          refTime,
				TimeUnitString:         "d",
				Count:                  "1",
				IntervalTimeUnitString: "d",
				IntervalTimeUnitCount:  "1",
				SrxTimeKind:            "LOCAL",
			},
		},
		Xsi:  "http://www.w3.org/2001/XMLSchema-instance",
		Xsd:  "http://www.w3.org/2001/XMLSchema",
		Soap: "http://schemas.xmlsoap.org/soap/envelope/",
	}

	if len(config.Config.Api.GetData.Tags) == 0 {
		var v []byte
		return v, fmt.Errorf("length of tag is equal to zero")
	}
	var tags Tags
	for _, v := range config.Config.Api.GetData.Tags {
		tags.SRxTag = append(tags.SRxTag, SRxTag{Name: v})
	}
	body.Body.GetSRxData.Tags = tags
	log.Printf("reqBody: %+v\n", body)

	// Encoding XML
	xmlBody, err := xml.Marshal(body)
	if err != nil {
		log.Fatalf("Error XML encoding, %s", err)
	}
	log.Printf("xmlBody: %+v\n", string(xmlBody))

	// Get data API
	resp, err := client.R().
		SetBasicAuth(config.Config.Api.GetData.BasicAuth.Username, config.Config.Api.GetData.BasicAuth.Password).
		SetHeaders(config.Config.Api.GetData.Headers).
		SetBody(xmlBody).
		Post(config.Config.Api.GetData.Url)
	if err != nil {
		fmt.Println("  Error      :", err)
	}

	fmt.Println("Get data API")
	fmt.Println("Response Info:")
	fmt.Println("  Error      :", err)
	fmt.Println("  Status Code:", resp.StatusCode())
	fmt.Println("  Status     :", resp.Status())
	fmt.Println("  Proto      :", resp.Proto())
	fmt.Println("  Time       :", resp.Time())
	fmt.Println("  Received At:", resp.ReceivedAt())
	fmt.Println("  Body       :\n", resp)
	fmt.Println()
	fmt.Println("=====================================================")

	return resp.Body(), err
}

func ParseXML(data []byte) (Response, error) {
	// Unmarshal
	var resp Response
	err := xml.Unmarshal([]byte(data), &resp)
	if err != nil {
		log.Println(err)
	}
	return resp, err
}
