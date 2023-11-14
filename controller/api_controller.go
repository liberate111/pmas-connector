package controller

import (
	"app-connector/constant"
	"app-connector/logger"
	"app-connector/model"
	"encoding/xml"
	"fmt"
	"log/slog"
	"time"
	_ "time/tzdata"
)

// Connect API
func ConnectAPI(req model.RequestApi) error {
	resp, err := client.R().
		SetBasicAuth(req.BasicAuth.Username, req.BasicAuth.Password).
		SetHeaders(req.Headers).
		SetBody(req.Body).
		Post(req.Url)
	if err != nil || resp.IsError() {
		logger.Error("connect API", "error", err.Error())
		return fmt.Errorf("connect API error %w", err)
	}
	logger.Debug("connect API", slog.Group("response", slog.Int("status", resp.StatusCode()), slog.Duration("response time", resp.Time()), slog.String("response body", resp.String())))
	logger.Info("connect API", "status", constant.SUCCESS, "url", req.Url)
	return nil
}

func GetDataAPI(req model.RequestApi) ([]byte, error) {
	var v []byte
	// Create request body
	t := time.Now().AddDate(0, 0, -2)
	loc := "Asia/Bangkok"
	zoneLoc, err := time.LoadLocation(loc)
	if err != nil {
		return nil, err
	}
	refTime := t.In(zoneLoc).Format(time.RFC3339)

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

	if len(req.Tags) == 0 {
		return v, fmt.Errorf("length of tag is equal to zero")
	}
	var tags model.Tags
	for _, v := range req.Tags {
		tags.SRxTag = append(tags.SRxTag, model.SRxTag{Name: v})
	}
	body.Body.GetSRxData.Tags = tags

	// Encoding XML
	xmlBody, err := xml.Marshal(body)
	if err != nil {
		return v, fmt.Errorf("xml encoding error %w", err)
	}
	logger.Debug("get data API", slog.Group("xml request", slog.String("body", string(xmlBody))))

	// Get data API
	resp, err := client.R().
		SetBasicAuth(req.BasicAuth.Username, req.BasicAuth.Password).
		SetHeaders(req.Headers).
		SetBody(xmlBody).
		Post(req.Url)
	if err != nil || resp.IsError() {
		return v, fmt.Errorf("request error %w", err)
	}
	logger.Debug("get data API", slog.Group("response", slog.Int("status", resp.StatusCode()), slog.Duration("response time", resp.Time()), slog.String("response body", resp.String())))
	logger.Info("get data API", "status", constant.SUCCESS, "url", req.Url)
	return resp.Body(), err
}
