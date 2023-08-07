package main

import (
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

// <Name>SPC_1CT0208AU02_ALM</Name>
// <Index>0</Index>
// <Archive>1d</Archive>
// <Unit>-</Unit>
// <FormulaOnTheFly>false</FormulaOnTheFly>
type TagData struct {
	XMLName         xml.Name `xml:"Tag"`
	Name            string   `xml:"Name"`
	Archive         string   `xml:"Archive"`
	Unit            string   `xml:"Unit"`
	FormulaOnTheFly string   `xml:"FormulaOnTheFly"`
}

// <State>
//
//	<IsValid>true</IsValid>
//	<ErrorMessage />
//	<ErrorCode>None</ErrorCode>
//
// </State>
type State struct {
	XMLName      xml.Name `xml:"State"`
	IsValid      string   `xml:"IsValid"`
	ErrorMessage string   `xml:"ErrorMessage"`
	ErrorCode    string   `xml:"ErrorCode"`
}

// <TimeDataItem>
// <Value>12</Value>
// <Time>2023-04-08T00:00:00+07:00</Time>
// </TimeDataItem>
type Data struct {
	XMLName      xml.Name `xml:"Data"`
	TimeDataItem []TimeDataItem
}

type TimeDataItem struct {
	XMLName xml.Name `xml:"TimeDataItem"`
	Value   string   `xml:"Value"`
	Time    string   `xml:"Time"`
}

func main() {
	fmt.Println("FUUUUUUUUUUUUUUUUUUUU")
	// Create a Resty Client
	client := resty.New()

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

	// Iqony
	// Connect API
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
		SetBasicAuth("pmaschnc2", "Pmas1a2b3C@EgaT").
		SetHeaders(map[string]string{
			"Content-Type": "text/xml",
			"SOAPAction":   "http://www.SR-Suite.com/SRxServerWebService/Connect",
			"Host":         "10.40.61.57",
		}).
		SetBody(bodyConnect).
		Post("http://10.40.61.57/srxwebservice/SRxServerWebService.asmx")
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
	resp, err = client.R().
		SetBasicAuth("chn-c2", "chn-c2pmas").
		SetHeaders(map[string]string{
			"Content-Type": "text/xml",
			"SOAPAction":   "http://www.SR-Suite.com/SRxServerWebService/GetSRxData",
			"Host":         "10.40.61.57",
		}).
		SetBody(bodyGetData).
		Post("http://10.40.61.57/srxwebservice/SRxServerWebService.asmx")
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

	// Unmarshal
	var respData Response
	if err := xml.Unmarshal([]byte(resp.Body()), &respData); err != nil {
		log.Fatal(err)
	}
	log.Printf("Data: %+v\n", respData)
	// log.Printf("Tag Name: %v\n", respData.SoapBody.Resp.ResultData.TagData.Name)
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
