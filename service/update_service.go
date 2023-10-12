package service

import (
	"app-connector/constant"
	"app-connector/controller"
	"app-connector/logger"
	"app-connector/model"
	"app-connector/util"
)

func UpdateBySite(site model.SiteConfig) {
	logger.Logger.Info("update", "status", constant.START, "site", site.Name, "table", site.TableName)
	reqCon := site.Api.Connect
	reqGet := site.Api.GetData
	table := site.TableName
	err := controller.ConnectAPI(reqCon)
	if err != nil {
		logger.Logger.Error("connect API", "error", err.Error(), "site", site.Name)
		return
	}
	respBody, err := controller.GetDataAPI(reqGet)
	if err != nil {
		logger.Logger.Error("get data API", "error", err.Error(), "site", site.Name)
	}

	tagData, err := util.ParseXML(respBody)
	if err != nil {
		logger.Logger.Error("parse xml data", "error", err.Error(), "site", site.Name)
	}
	// logger.Logger.Debug("result data", "data", tagData)

	err = updateStatus(tagData, table)
	if err != nil {
		logger.Logger.Error("update status", "error", err.Error(), "site", site.Name)
		return
	}
	logger.Logger.Info("update", "status", constant.SUCCESS, "site", site.Name, "table", site.TableName)
}

func updateStatus(r model.Response, table string) error {
	controller.CreateStmt(table)
	var sform, status string
	tsz, err := util.Timestampt()
	if err != nil {
		return err
	}
	// tags loop
	for _, v := range r.SoapBody.Resp.ResultData {
		if v.State.IsValid == constant.STATE_FALSE {
			logger.Logger.Error("update", "error", "response data from API state is not valid", "tag", v.TagData.Name)
			continue
		}
		if len(v.Data.TimeDataItem) != 2 {
			logger.Logger.Error("update", "error", "length of Data is not equal to 2", "tag", v.TagData.Name)
			continue
		}
		if v.Data.TimeDataItem[0].Value == v.Data.TimeDataItem[1].Value {
			logger.Logger.Debug("update", "tag", v.TagData.Name, "check status", "status not change")
			continue
		}

		// update
		sform, err = util.ConvertStatus(v.Data.TimeDataItem[0].Value)
		if err != nil {
			logger.Logger.Error("update", "error", err.Error(), "tag", v.TagData.Name)
			continue
		}
		status, err = util.ConvertStatus(v.Data.TimeDataItem[1].Value)
		if err != nil {
			logger.Logger.Error("update", "error", err.Error(), "tag", v.TagData.Name)
			continue
		}
		err := controller.UpdateStatus(status, sform, v.TagData.Name, tsz, table)
		if err != nil {
			logger.Logger.Error("update", "error", err.Error(), "tag", v.TagData.Name)
		}
	}
	logger.Logger.Info("update", "status", constant.SUCCESS, "tag", "all tags")
	controller.CloseStmt()
	return nil
}
