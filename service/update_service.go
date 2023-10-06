package service

import (
	"app-connector/config"
	"app-connector/constant"
	"app-connector/controller"
	"app-connector/logger"
	"app-connector/model"
	"app-connector/util"
	"database/sql/driver"
	"fmt"
)

func UpdateBySite() {
	controller.ConnectAPI()
	respBody, err := controller.GetDataAPI()
	if err != nil {
		logger.Logger.Error("get data API", "error", err.Error())
	}

	tagData, err := util.ParseXML(respBody)
	if err != nil {
		logger.Logger.Error("parse xml data", "error", err.Error())
	}
	logger.Logger.Debug("result data", "data", tagData)

	err = updateStatus(tagData)
	if err != nil {
		logger.Logger.Error("check status", "error", err.Error())
	}
	logger.Logger.Info("site", "status", constant.SUCCESS)
}

func updateStatus(r model.Response) error {
	// stmt
	sql := fmt.Sprintf(`UPDATE %s SET STATUS = :1, SFROM = :2, "TimeStamp" = :3 WHERE KPI_TAG = :4`, config.Config.DB.Oracle.TableName)
	updStmt, err := controller.ConnOracle.Prepare(sql)
	if err != nil {
		logger.Logger.Error("update to db", "error", err.Error())
	}

	defer func() {
		_ = updStmt.Close()
	}()

	var sform, status string
	tsz, err := util.Timestamptz()
	if err != nil {
		logger.Logger.Error("update to db", "error", err.Error())
	}
	// tags loop
	for _, v := range r.SoapBody.Resp.ResultData {
		if len(v.Data.TimeDataItem) != 2 {
			logger.Logger.Error("update to db", "error", "length of Data is not equal to 2", "tag", v.TagData.Name)
			continue
		}
		if v.Data.TimeDataItem[0].Value == v.Data.TimeDataItem[1].Value {
			logger.Logger.Debug("update to db", "tag", v.TagData.Name, "check status", "status not change")
			continue
		}

		// update
		sform, err = util.ConvertStatus(v.Data.TimeDataItem[0].Value)
		if err != nil {
			logger.Logger.Error("update to db", "error", err.Error())
			continue
		}
		status, err = util.ConvertStatus(v.Data.TimeDataItem[1].Value)
		if err != nil {
			logger.Logger.Error("update to db", "error", err.Error())
			continue
		}

		_, err = updStmt.Exec([]driver.Value{status, sform, tsz, v.TagData.Name})
		if err != nil {
			logger.Logger.Error("update to db", "error", err.Error())
		}
	}
	logger.Logger.Info("update to db", "status", constant.SUCCESS)
	return nil
}
