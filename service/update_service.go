package service

import (
	"app-connector/constant"
	"app-connector/controller"
	"app-connector/logger"
	"app-connector/model"
	"app-connector/util"
)

func UpdateBySite(reqCon, reqGet model.RequestApi, table string) {
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

	err = updateStatus(tagData, table)
	if err != nil {
		logger.Logger.Error("update status", "error", err.Error())
		return
	}
	logger.Logger.Info("site", "status", constant.SUCCESS)
}

func updateStatus(r model.Response, table string) error {
	controller.CreateStmt(table)
	var sform, status string
	tsz, err := util.Timestampt()
	if err != nil {
		logger.Logger.Error("update to db", "error", err.Error())
	}
	// tags loop
	for _, v := range r.SoapBody.Resp.ResultData {
		if v.State.IsValid == constant.STATE_FALSE {
			logger.Logger.Error("update to db", "error", "response data from API state is not valid", "tag", v.TagData.Name)
			continue
		}
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
			logger.Logger.Error("update to db", "error", err.Error(), "tag", v.TagData.Name)
			continue
		}
		status, err = util.ConvertStatus(v.Data.TimeDataItem[1].Value)
		if err != nil {
			logger.Logger.Error("update to db", "error", err.Error(), "tag", v.TagData.Name)
			continue
		}
		err := controller.UpdateStatus(status, sform, v.TagData.Name, tsz, table)
		if err != nil {
			logger.Logger.Error("update to db", "error", err.Error(), "tag", v.TagData.Name)
		}
	}
	logger.Logger.Info("update to db", "status", constant.SUCCESS, "tag", "all tags")
	controller.CloseStmt()
	return nil
}
