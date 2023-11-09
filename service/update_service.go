package service

import (
	"app-connector/constant"
	"app-connector/controller"
	"app-connector/logger"
	"app-connector/model"
	"app-connector/util"
)

func UpdateBySite(site model.SiteConfig) {
	logger.Info("event: update_status_tag, status:", constant.START, ", table:", site.TableName, ", site:", site.Name)
	reqCon := site.Api.Connect
	reqGet := site.Api.GetData
	table := site.TableName
	err := controller.ConnectAPI(reqCon)
	if err != nil {
		logger.Error("event: connect_API, status: error, msg:", err.Error(), ", site", site.Name)
		return
	}
	respBody, err := controller.GetDataAPI(reqGet)
	if err != nil {
		logger.Error("event: get_data_API, status: error, msg:", err.Error(), ", site", site.Name)
		return
	}

	tagData, err := util.ParseXML(respBody)
	if err != nil {
		logger.Error("event: parse_xml_data, status: error, msg:", err.Error(), ", site", site.Name)
		return
	}

	err = updateStatus(tagData, table, site.Name)
	if err != nil {
		logger.Error("event: update_status_tag, status: error, msg:", err.Error(), ", site", site.Name)
		return
	}
	logger.Info("event: update_status_tag, status:", constant.SUCCESS, ", table:", site.TableName, ", site:", site.Name)
}

func updateStatus(r model.Response, table string, site string) error {
	controller.CreateStmt(table)
	var sform, status string
	tsz, err := util.Timestampt()
	if err != nil {
		return err
	}
	// tags loop
	for _, v := range r.SoapBody.Resp.ResultData {
		if v.State.IsValid == constant.STATE_FALSE {
			logger.Error("event: update_status_tag_validate_response, status: error, msg: response data from API state is not valid,", "tag:", v.TagData.Name, ", site:", site)
			continue
		}
		if len(v.Data.TimeDataItem) != 2 {
			logger.Error("event: update_status_tag_validate_response, status: error, msg: length of data is not equal to 2,", "tag:", v.TagData.Name, ", site:", site)
			continue
		}
		if v.Data.TimeDataItem[0].Value == string(util.STATUS_NAN) || v.Data.TimeDataItem[1].Value == string(util.STATUS_NAN) {
			logger.Error("event: update_status_tag_validate_response, status: error, msg: status of tag is NaN,", "tag:", v.TagData.Name, ", site:", site)
			continue
		}
		if v.Data.TimeDataItem[0].Value == v.Data.TimeDataItem[1].Value {
			logger.Debug("event: update_status_tag_check_status_api, msg: status is not change from API,", "tag:", v.TagData.Name, ", site:", site)
			lastStatus, err := controller.FindByTag(v.TagData.Name, table)
			if err != nil {
				logger.Error("event: update_status_tag_check_status_table, status: error, msg:", err.Error(), ", tag:", v.TagData.Name, ", site:", site)
				continue
			}
			if v.Data.TimeDataItem[1].Value == lastStatus {
				logger.Debug("event: update_status_tag_check_status_table, msg: status is not change from table,", "tag:", v.TagData.Name, ", site:", site)
				continue
			}
			v.Data.TimeDataItem[0].Value = lastStatus
			logger.Debug("event: update_status_tag_check_status_table, msg: status is change check from table,", "tag:", v.TagData.Name, ", site:", site)
		}
		// update
		sform, err = util.ConvertStatus(v.Data.TimeDataItem[0].Value)
		if err != nil {
			logger.Error("event: update_status_tag_validate_status, status: error, msg:", err.Error(), ", tag:", v.TagData.Name, ", site:", site)
			continue
		}
		status, err = util.ConvertStatus(v.Data.TimeDataItem[1].Value)
		if err != nil {
			logger.Error("event: update_status_tag_validate_status, status: error, msg:", err.Error(), ", tag:", v.TagData.Name, ", site:", site)
			continue
		}
		err := controller.UpdateStatus(status, sform, v.TagData.Name, tsz, table)
		if err != nil {
			logger.Error("event: update_status_tag, status: error, msg:", err.Error(), ", tag:", v.TagData.Name, ", site:", site)
		}
	}
	controller.CloseStmt()
	return nil
}
