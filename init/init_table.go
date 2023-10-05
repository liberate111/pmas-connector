package init

func InitTable() {
	// // API
	// ConnectAPI()
	// respBody, err := GetDataAPI()
	// if err != nil {
	// 	logger.Logger.Error("get data API", "error", err.Error())
	// }

	// tagData, err := ParseXML(respBody)
	// if err != nil {
	// 	logger.Logger.Error("parse xml data", "error", err.Error())
	// }
	// logger.Logger.Debug("result data", "data", tagData)

	// // DB
	// sql := fmt.Sprintf(`UPDATE %s SET STATUS = :1, SFROM = :2, "TimeStamp" = :3 WHERE KPI_TAG = :4`, config.Config.DB.Oracle.TableName)
	// updStmt, err := connOracle.Prepare(sql)
	// if err != nil {
	// 	logger.Logger.Error("initial data in table", "error", err.Error())
	// }

	// defer func() {
	// 	_ = updStmt.Close()
	// }()

	// var sform, status string
	// var tsz time.Time
	// tsz, err = timestamptz()
	// if err != nil {
	// 	logger.Logger.Error("initial data in table", "error", err.Error())
	// }

	// for _, v := range tagData.SoapBody.Resp.ResultData {
	// 	sform, err = convertStatus(v.Data.TimeDataItem[0].Value)
	// 	if err != nil {
	// 		logger.Logger.Error("initial data in table", "error", err.Error())
	// 		continue
	// 	}
	// 	status, err = convertStatus(v.Data.TimeDataItem[1].Value)
	// 	if err != nil {
	// 		logger.Logger.Error("initial data in table", "error", err.Error())
	// 		continue
	// 	}

	// 	_, err = updStmt.Exec([]driver.Value{status, sform, go_ora.TimeStampTZ(tsz), v.TagData.Name})
	// 	if err != nil {
	// 		logger.Logger.Error("initial data in table", "error", err.Error())
	// 	}
	// 	logger.Logger.Info("initail table", "status", SUCCESS)
	// }
}
