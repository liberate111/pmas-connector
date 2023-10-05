package service

func updateBySite() {
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

	// err = updateStatus(tagData)
	// if err != nil {
	// 	logger.Logger.Error("check status", "error", err.Error())
	// }
	// logger.Logger.Info("site", "status", SUCCESS)
}

// func updateStatus(r Response) error {
// // stmt
// sql := fmt.Sprintf(`UPDATE %s SET STATUS = :1, SFROM = :2, "TimeStamp" = :3 WHERE KPI_TAG = :4`, config.Config.DB.Oracle.TableName)
// updStmt, err := connOracle.Prepare(sql)
// if err != nil {
// 	logger.Logger.Error("update to db", "error", err.Error())
// }

// defer func() {
// 	_ = updStmt.Close()
// }()

// var sform, status string
// var tsz time.Time
// tsz, err = timestamptz()
// if err != nil {
// 	logger.Logger.Error("update to db", "error", err.Error())
// }
// // tags loop
// for _, v := range r.SoapBody.Resp.ResultData {
// 	if len(v.Data.TimeDataItem) != 2 {
// 		logger.Logger.Error("update to db", "error", "length of Data is not equal to 2", "tag", v.TagData.Name)
// 		continue
// 	}
// 	if v.Data.TimeDataItem[0].Value == v.Data.TimeDataItem[1].Value {
// 		logger.Logger.Debug("update to db", "tag", v.TagData.Name, "check status", "status not change")
// 		continue
// 	}

// 	// update
// 	sform, err = convertStatus(v.Data.TimeDataItem[0].Value)
// 	if err != nil {
// 		logger.Logger.Error("update to db", "error", err.Error())
// 		continue
// 	}
// 	status, err = convertStatus(v.Data.TimeDataItem[1].Value)
// 	if err != nil {
// 		logger.Logger.Error("update to db", "error", err.Error())
// 		continue
// 	}

// 	_, err = updStmt.Exec([]driver.Value{status, sform, go_ora.TimeStampTZ(tsz), v.TagData.Name})
// 	if err != nil {
// 		logger.Logger.Error("update to db", "error", err.Error())
// 	}
// }
// logger.Logger.Info("update to db", "status", SUCCESS)
// return nil
// }
