package main

import (
	"app-connector/config"
	"app-connector/logger"
	"encoding/xml"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/robfig/cron/v3"
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
	cronjob()
}

func gracefulShutdown(done chan bool) {
	// Gracefully Shutdown
	// Make channel listen for signals from OS
	go func() {
		c := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
		signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		<-c
		logger.Logger.Info("Application shutdown...")
		done <- true
	}()
}

func cronjob() {
	localTime, _ := time.LoadLocation("Asia/Bangkok")
	c := cron.New(cron.WithLocation(localTime))

	// chn-c2
	// c.AddFunc("@midnight", chn_c2)
	// for test
	c.AddFunc("@every 30s", chn_c2)
	c.Start()
	done := make(chan bool, 1)
	gracefulShutdown(done)
	<-done
	defer c.Stop()
}

func chn_c2() {
	ConnectAPI()
	respBody, err := GetDataAPI()
	if err != nil {
		logger.Logger.Error("get data API", "error", err.Error())
	}

	tagData, err := ParseXML(respBody)
	if err != nil {
		logger.Logger.Error("parse xml data", "error", err.Error())
	}
	logger.Logger.Debug("result data", "data", tagData)

	// business logic
	err = checkStatus(tagData)
	if err != nil {
		logger.Logger.Error("check status", "error", err.Error())
	}
}

func InitClient() {
	// Create a Resty Client
	client = resty.New()
	client.DisableWarn = false
	client.
		// Set retry count to non zero to enable retries
		SetRetryCount(3).
		// You can override initial retry wait time.
		// Default is 100 milliseconds.
		SetRetryWaitTime(5 * time.Second).
		// MaxWaitTime can be overridden as well.
		// Default is 2 seconds.
		SetRetryMaxWaitTime(20 * time.Second).
		// SetRetryAfter sets callback to calculate wait time between retries.
		// Default (nil) implies exponential backoff with jitter
		SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
			return 0, errors.New("quota exceeded")
		})
}

func ConnectAPI() {
	// Connect API
	resp, err := client.R().
		SetBasicAuth(config.Config.Api.Connect.BasicAuth.Username, config.Config.Api.Connect.BasicAuth.Password).
		SetHeaders(config.Config.Api.Connect.Headers).
		SetBody(config.Config.Api.Connect.Body).
		Post(config.Config.Api.Connect.Url)
	if err != nil {
		logger.Logger.Error("connect API", "error", err.Error())
	}
	logger.Logger.Debug("connect API", slog.Group("response", slog.Int("status", resp.StatusCode()), slog.Duration("response time", resp.Time()), slog.String("response body", resp.String())))
}

func GetDataAPI() ([]byte, error) {
	var v []byte
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
		return v, fmt.Errorf("length of tag is equal to zero")
	}
	var tags Tags
	for _, v := range config.Config.Api.GetData.Tags {
		tags.SRxTag = append(tags.SRxTag, SRxTag{Name: v})
	}
	body.Body.GetSRxData.Tags = tags
	logger.Logger.Debug("get data API", slog.Group("request body", body))

	// Encoding XML
	xmlBody, err := xml.Marshal(body)
	if err != nil {
		return v, fmt.Errorf("xml encoding error %w", err)
	}
	logger.Logger.Debug("get data API", "xml request body", string(xmlBody))

	// Get data API
	resp, err := client.R().
		SetBasicAuth(config.Config.Api.GetData.BasicAuth.Username, config.Config.Api.GetData.BasicAuth.Password).
		SetHeaders(config.Config.Api.GetData.Headers).
		SetBody(xmlBody).
		Post(config.Config.Api.GetData.Url)
	if err != nil {
		return v, fmt.Errorf("request error %w", err)
	}
	logger.Logger.Debug("get data API", slog.Group("response", slog.Int("status", resp.StatusCode()), slog.Duration("response time", resp.Time()), slog.String("response body", resp.String())))
	return resp.Body(), err
}

func ParseXML(data []byte) (Response, error) {
	// Unmarshal
	var resp Response
	err := xml.Unmarshal([]byte(data), &resp)
	if err != nil {
		return resp, fmt.Errorf("parse xml encode data error %w", err)
	}
	return resp, err
}

func checkStatus(r Response) error {
	// tags loop
	for _, v := range r.SoapBody.Resp.ResultData {
		if len(v.Data.TimeDataItem) != 2 {
			logger.Logger.Error("check status", "error", "status value greater than 2", "tag", v.TagData.Name)
			continue
		}
		if v.Data.TimeDataItem[0].Value == v.Data.TimeDataItem[1].Value {
			logger.Logger.Debug("status not change", "tag", v.TagData.Name)
			continue
		}
		// status change
		// update to DB
		logger.Logger.Debug("update to db")
	}
	return nil
}
