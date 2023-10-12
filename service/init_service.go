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
	site := config.Config.Site
	for _, v := range site {
		initBySite(v)
	}
}

func initBySite(site model.SiteConfig) {
	logger.Logger.Info("initial table", "status", constant.START, "site", site.Name, "table", site.TableName)
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

	err = initStatus(tagData, table)
	if err != nil {
		logger.Logger.Error("initial table", "error", err.Error(), "site", site.Name)
		return
	}
	logger.Logger.Info("initial table", "status", constant.SUCCESS, "site", site.Name, "table", site.TableName)
}

func initStatus(r model.Response, table string) error {
	controller.CreateStmt(table)
	var sform, status string
	tsz, err := util.Timestampt()
	if err != nil {
		return err
	}
	// tags loop
	for _, v := range r.SoapBody.Resp.ResultData {
		if v.State.IsValid == constant.STATE_FALSE {
			logger.Logger.Error("initial table", "error", "response data from API state is not valid", "tag", v.TagData.Name)
			continue
		}
		if len(v.Data.TimeDataItem) != 2 {
			logger.Logger.Error("initial table", "error", "length of Data is not equal to 2", "tag", v.TagData.Name)
			continue
		}

		// update
		sform, err = util.ConvertStatus(v.Data.TimeDataItem[0].Value)
		if err != nil {
			logger.Logger.Error("initial table", "error", err.Error())
			continue
		}
		status, err = util.ConvertStatus(v.Data.TimeDataItem[1].Value)
		if err != nil {
			logger.Logger.Error("initial table", "error", err.Error())
			continue
		}
		err := controller.UpdateStatus(status, sform, v.TagData.Name, tsz, table)
		if err != nil {
			logger.Logger.Error("initial table", "error", err.Error())
		}
	}
	controller.CloseStmt()
	return nil
}
