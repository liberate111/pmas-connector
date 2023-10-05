package model

import "encoding/xml"

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
