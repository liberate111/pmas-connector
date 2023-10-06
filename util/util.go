package util

import (
	"app-connector/model"
	"encoding/xml"
	"fmt"
	"time"

	go_ora "github.com/sijms/go-ora/v2"
)

type TagStatus string

const NUM_RED TagStatus = "14"
const NUM_YELLOW TagStatus = "13"
const NUM_GREEN TagStatus = "12"
const NUM_BLUE TagStatus = "11"
const NUM_MAGENTA TagStatus = "2"
const NUM_GREY TagStatus = "1"
const STATUS_RED TagStatus = "RED"
const STATUS_YELLOW TagStatus = "YELLOW"
const STATUS_GREEN TagStatus = "GREEN"
const STATUS_MAGENTA TagStatus = "MAGENTA"
const STATUS_BLUE TagStatus = "BLUE"
const STATUS_GREY TagStatus = "GREY"

func ParseXML(data []byte) (model.Response, error) {
	// Unmarshal
	var resp model.Response
	err := xml.Unmarshal([]byte(data), &resp)
	if err != nil {
		return resp, fmt.Errorf("parse xml encode data error %w", err)
	}
	return resp, err
}

func ConvertStatus(num string) (string, error) {
	if num == string(NUM_RED) {
		return string(STATUS_RED), nil
	} else if num == string(NUM_YELLOW) {
		return string(STATUS_YELLOW), nil
	} else if num == string(NUM_GREEN) {
		return string(STATUS_GREEN), nil
	} else if num == string(NUM_MAGENTA) {
		return string(STATUS_MAGENTA), nil
	} else if num == string(NUM_BLUE) {
		return string(STATUS_BLUE), nil
	} else if num == string(NUM_GREY) {
		return string(STATUS_GREY), nil
	}
	return "", fmt.Errorf("error not match any status type: %v", num)
}

func Timestamptz() (go_ora.TimeStampTZ, error) {
	t := time.Now()
	loc := "Asia/Bangkok"
	zoneLoc, err := time.LoadLocation(loc)
	if err != nil {
		return go_ora.TimeStampTZ{}, err
	}
	return go_ora.TimeStampTZ(t.In(zoneLoc)), nil
}
