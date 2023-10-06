package controller

import (
	"app-connector/config"
	"app-connector/constant"
	"app-connector/logger"
	"app-connector/model"
	"encoding/xml"
	"fmt"
	"log/slog"
	"time"
)

// Connect API
func ConnectAPI() {
	resp, err := client.R().
		SetBasicAuth(config.Config.Api.Connect.BasicAuth.Username, config.Config.Api.Connect.BasicAuth.Password).
		SetHeaders(config.Config.Api.Connect.Headers).
		SetBody(config.Config.Api.Connect.Body).
		Post(config.Config.Api.Connect.Url)
	if err != nil {
		logger.Logger.Error("connect API", "error", err.Error())
	}
	logger.Logger.Debug("connect API", slog.Group("response", slog.Int("status", resp.StatusCode()), slog.Duration("response time", resp.Time()), slog.String("response body", resp.String())))
	logger.Logger.Info("connect API", "status", constant.SUCCESS)
}

func GetDataAPI() ([]byte, error) {
	var v []byte
	// Create request body
	refTime := time.Now().AddDate(0, 0, -1).Format(time.RFC3339)
	body := model.SOAPEnvelope{
		Body: model.SOAPBody{
			GetSRxData: model.GetSRxData{
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
	var tags model.Tags
	for _, v := range config.Config.Api.GetData.Tags {
		tags.SRxTag = append(tags.SRxTag, model.SRxTag{Name: v})
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
	logger.Logger.Info("get data API", "status", constant.SUCCESS)
	return resp.Body(), err
}
