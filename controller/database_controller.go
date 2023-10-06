package controller

import (
	"app-connector/config"
	"app-connector/constant"
	"app-connector/logger"
	"app-connector/model"
	"database/sql/driver"
	"fmt"

	go_ora "github.com/sijms/go-ora/v2"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

const (
	ORACLE_DB    string = "oracle"
	SQLSERVER_DB string = "sqlserver"
)

var (
	ConnOracle    *go_ora.Connection
	ConnSqlserver *gorm.DB
	Driver        string
	dbConfig      model.DatabaseDriver
	stmt          *go_ora.Stmt
)

func ConnectDB() error {
	var err error
	Driver = config.Config.App.Db
	dbConfig = config.Config.DB
	if Driver == ORACLE_DB {
		connStr := go_ora.BuildUrl(dbConfig.Oracle.Url, dbConfig.Oracle.Port, dbConfig.Oracle.ServiceName, dbConfig.Oracle.Username, dbConfig.Oracle.Password, nil)
		ConnOracle, err = go_ora.NewConnection(connStr)
		if err != nil {
			logger.Logger.Error("connect to db", "error", err.Error())
			return err
		}
		err = ConnOracle.Open()
		if err != nil {
			logger.Logger.Error("connect to db", "error", err.Error())
			return err
		}
		logger.Logger.Info("connect to db", "status", constant.SUCCESS)
		return nil
	} else if Driver == SQLSERVER_DB {
		// TODO : connect db
		dsn := "sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm"
		ConnSqlserver, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
		return nil
	} else {
		return fmt.Errorf("driver database not support: %v", Driver)
	}
}

func CloseDB() {
	if Driver == ORACLE_DB {
		err := ConnOracle.Close()
		if err != nil {
			logger.Logger.Error("close db connection", "error", err.Error())
		}
	}
	if Driver == SQLSERVER_DB {
		// TODO : disconnect db

	}
}

func CreateStmt() error {
	if Driver == ORACLE_DB {
		// stmt
		sql := fmt.Sprintf(`UPDATE %s SET STATUS = :1, SFROM = :2, "TimeStamp" = :3 WHERE KPI_TAG = :4`, dbConfig.Oracle.TableName)
		// updStmt, err := ConnOracle.Prepare(sql)
		// if err != nil {
		// 	logger.Logger.Error("update to db", "error", err.Error())
		// }
		stmt = go_ora.NewStmt(sql, ConnOracle)
		return nil
	} else if Driver == SQLSERVER_DB {

	} else {
		return fmt.Errorf("driver database not support: %v", Driver)
	}
	return nil
}

func CloseStmt() {
	if Driver == ORACLE_DB {
		_ = stmt.Close()
	}
}

func UpdateStatus(status, sform, tag string, tsz go_ora.TimeStampTZ) error {
	if Driver == ORACLE_DB {
		_, err := stmt.Exec([]driver.Value{status, sform, tsz, tag})
		if err != nil {
			logger.Logger.Error("update to db", "error", err.Error())
		}
	} else if Driver == SQLSERVER_DB {

	} else {
		return fmt.Errorf("driver database not support: %v", Driver)
	}
	return nil
}
