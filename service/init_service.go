package service

import (
	"app-connector/config"
	"app-connector/constant"
	"app-connector/controller"
	"app-connector/logger"
	"app-connector/model"
	"app-connector/util"
)

func InitTable() {
	// loop for all site here

	// manual
	reqCon := config.Config.Api.Connect
	reqGet := config.Config.Api.GetData
	initBySite(reqCon, reqGet)
}

func initBySite(reqCon, reqGet model.RequestApi) {
	err := controller.ConnectAPI(reqCon)
	if err != nil {
		logger.Logger.Error("connect API", "error", err.Error())
		return
	}
	respBody, err := controller.GetDataAPI(reqGet)
	if err != nil {
		logger.Logger.Error("get data API", "error", err.Error())
	}

	tagData, err := util.ParseXML(respBody)
	if err != nil {
		logger.Logger.Error("parse xml data", "error", err.Error())
	}
	// logger.Logger.Debug("result data", "data", tagData)

	err = initStatus(tagData)
	if err != nil {
		logger.Logger.Error("update status", "error", err.Error())
		return
	}
	logger.Logger.Info("site", "status", constant.SUCCESS)
}

func initStatus(r model.Response) error {
	controller.CreateStmt()
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
		// if v.Data.TimeDataItem[0].Value == v.Data.TimeDataItem[1].Value {
		// 	logger.Logger.Debug("update to db", "tag", v.TagData.Name, "check status", "status not change")
		// 	continue
		// }

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
		err := controller.UpdateStatus(status, sform, v.TagData.Name, tsz)
		if err != nil {
			logger.Logger.Error("update to db", "error", err.Error())
		}
	}
	logger.Logger.Info("update to db", "status", constant.SUCCESS)
	controller.CloseStmt()
	return nil
}
