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
	logger.Info("initial table", "table", site.TableName, "status", constant.START, "site", site.Name)
	reqCon := site.Api.Connect
	reqGet := site.Api.GetData
	table := site.TableName
	err := controller.ConnectAPI(reqCon)
	if err != nil {
		logger.Error("connect API", "error", err.Error(), "site", site.Name)
		return
	}
	respBody, err := controller.GetDataAPI(reqGet)
	if err != nil {
		logger.Error("get data API", "error", err.Error(), "site", site.Name)
		return
	}

	tagData, err := util.ParseXML(respBody)
	if err != nil {
		logger.Error("parse xml data", "error", err.Error(), "site", site.Name)
		return
	}

	err = initStatus(tagData, table)
	if err != nil {
		logger.Error("initial table", "error", err.Error(), "site", site.Name)
		return
	}
	logger.Info("initial table", "table", site.TableName, "status", constant.SUCCESS, "site", site.Name)
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
			logger.Error("initial table", "error", "response data from API state is not valid", "tag", v.TagData.Name)
			continue
		}
		if len(v.Data.TimeDataItem) != 2 {
			logger.Error("initial table", "error", "length of data is not equal to 2", "tag", v.TagData.Name)
			continue
		}
		if v.Data.TimeDataItem[0].Value == string(util.STATUS_NAN) || v.Data.TimeDataItem[1].Value == string(util.STATUS_NAN) {
			logger.Error("initial table", "error", "tag", v.TagData.Name, "status of tag", "NaN")
			continue
		}

		// update
		sform, err = util.ConvertStatus(v.Data.TimeDataItem[0].Value)
		if err != nil {
			logger.Error("initial table", "error", err.Error())
			continue
		}
		status, err = util.ConvertStatus(v.Data.TimeDataItem[1].Value)
		if err != nil {
			logger.Error("initial table", "error", err.Error())
			continue
		}
		err := controller.UpdateStatus(status, sform, v.TagData.Name, tsz, table)
		if err != nil {
			logger.Error("initial table", "error", err.Error())
		}
	}
	controller.CloseStmt()
	return nil
}
