package main

import (
	"app-connector/config"
	"encoding/xml"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

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
	XMLName xml.Name `xml:"GetSRxDataResponse"`
	// XMLName    xml.Name   `xml:"http://www.SR-Suite.com/SRxServerWebService"`
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

type RequestConnectAPI struct {
	BasicAuth BasicAuth
	Headers   map[string]string
	Url       string
}

type RequestGetDataAPI struct {
	BasicAuth BasicAuth
	Headers   map[string]string
	Url       string
}

type BasicAuth struct {
	Username string
	Password string
}

var client *resty.Client

// TODO Config
var reqConnect = RequestConnectAPI{BasicAuth: BasicAuth{Username: "pmaschnc2", Password: "Pmas1a2b3C@EgaT"}, Headers: map[string]string{"Content-Type": "text/xml", "SOAPAction": "http://www.SR-Suite.com/SRxServerWebService/Connect", "Host": "10.40.61.57"}, Url: "http://10.40.61.57/srxwebservice/SRxServerWebService.asmx"}
var reqGetData = RequestConnectAPI{BasicAuth: BasicAuth{Username: "chn-c2", Password: "chn-c2pmas"}, Headers: map[string]string{"Content-Type": "text/xml", "SOAPAction": "http://www.SR-Suite.com/SRxServerWebService/GetSRxData", "Host": "10.40.61.57"}, Url: "http://10.40.61.57/srxwebservice/SRxServerWebService.asmx"}

func main() {
	fmt.Println("Start PMAS-CONNECTOR Application...")
	config.ReadConfig()
	log.Println(config.Config.App)

	// InitClient()
	// ConnectAPI()
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

	// resp, err := client.R().
	// 	EnableTrace().
	// 	Get("https://httpbin.org/get")

	// // Explore response object
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

func ConnectAPI() {
	// Iqony
	// Connect API

	// Body
	bodyConnect := `<?xml version="1.0" encoding="utf-8"?>
	<soap12:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap12="http://schemas.xmlsoap.org/soap/envelope/">
	  <soap12:Body>
		<Connect xmlns="http://www.SR-Suite.com/SRxServerWebService">
		  <serverName>SR.CHN2.Online</serverName>
		  <version>1</version>
		</Connect>
	  </soap12:Body>
	</soap12:Envelope>`

	// POST JSON string
	// No need to set content type, if you have client level setting
	resp, err := client.R().
		// SetBasicAuth("pmaschnc2", "Pmas1a2b3C@EgaT").
		SetBasicAuth(reqConnect.BasicAuth.Username, reqConnect.BasicAuth.Password).
		SetHeaders(reqConnect.Headers).
		SetBody(bodyConnect).
		Post(reqConnect.Url)
	if err != nil {
		fmt.Println("  Error      :", err)
	}
	fmt.Println("Connect API")
	fmt.Println("Response Info:")
	fmt.Println("  Error      :", err)
	fmt.Println("  Status Code:", resp.StatusCode())
	fmt.Println("  Status     :", resp.Status())
	fmt.Println("  Proto      :", resp.Proto())
	fmt.Println("  Time       :", resp.Time())
	fmt.Println("  Received At:", resp.ReceivedAt())
	fmt.Println("  Body       :\n", resp)
	fmt.Println()
}

func GetDataAPI() ([]byte, error) {
	bodyGetData := `<?xml version="1.0" encoding="utf-8"?>
	<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
	  <soap:Body>
		<GetSRxData xmlns="http://www.SR-Suite.com/SRxServerWebService">
		  <tags>
			<SRxTag>
			  <Name>SPC_21MKD12CY021XQ01_ALM</Name>
			  <Index>0</Index>
			  <Archive>1d</Archive>
			  <Unit></Unit>
			  <FormulaOnTheFly>0</FormulaOnTheFly>
			</SRxTag>
			<SRxTag>
			  <Name>SPC_21MBL10CP003XQ01_T_ALM</Name>
			  <Index>0</Index>
			  <Archive>1d</Archive>
			  <Unit></Unit>
			  <FormulaOnTheFly>0</FormulaOnTheFly>
			</SRxTag>
		  </tags>
		  <referenceDate>2023-08-02T01:00:00Z</referenceDate>
		  <timeUnitString>d</timeUnitString>
		  <count>1</count>
		  <intervalTimeUnitString>d</intervalTimeUnitString>
		  <intervalTimeUnitCount>1</intervalTimeUnitCount>
		  <srxTimeKind>LOCAL</srxTimeKind>
		</GetSRxData>
	  </soap:Body>
	</soap:Envelope>`
	// Get data API
	resp, err := client.R().
		SetBasicAuth(reqGetData.BasicAuth.Username, reqGetData.BasicAuth.Password).
		SetHeaders(reqGetData.Headers).
		SetBody(bodyGetData).
		Post(reqGetData.Url)
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

	// Version SOAP 1.2
	// Connect API
	// bodyConnect = `<?xml version="1.0" encoding="utf-8"?>
	// <soap12:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap12="http://schemas.xmlsoap.org/soap/envelope/">
	//   <soap12:Body>
	// 	<Connect xmlns="http://www.SR-Suite.com/SRxServerWebService">
	// 	  <serverName>SR.NBC1.Online</serverName>
	// 	  <version>1</version>
	// 	</Connect>
	//   </soap12:Body>
	// </soap12:Envelope>`
	// // POST JSON string
	// // No need to set content type, if you have client level setting
	// resp, err = client.R().
	// 	SetBasicAuth("pmasnbc1", "Pmas1a2b3C@EgaT").
	// 	SetHeaders(map[string]string{
	// 		"Content-Type": "application/soap+xml; charset=utf-8",
	// 		"SOAPAction":   "http://www.SR-Suite.com/SRxServerWebService/Connect",
	// 		"Host":         "10.40.61.56",
	// 	}).
	// 	SetBody(bodyConnect).
	// 	Post("http://10.40.61.56/srxwebservice/SRxServerWebService.asmx")
	// if err != nil {
	// 	fmt.Println("  Error      :", err)
	// }
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

func ParseXML(data []byte) (Response, error) {
	// Unmarshal
	var resp Response
	err := xml.Unmarshal([]byte(data), &resp)
	if err != nil {
		log.Println(err)
	}
	return resp, err
}
