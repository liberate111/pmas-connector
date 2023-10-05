package controller

import (
	"app-connector/config"
	"app-connector/constant"
	"app-connector/logger"
	"fmt"

	go_ora "github.com/sijms/go-ora/v2"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var connOracle *go_ora.Connection
var connSqlserver *gorm.DB

const (
	ORACLE_DB    string = "oracle"
	SQLSERVER_DB string = "sqlserver"
)

func ConnectDB() error {
	var err error
	driver := config.Config.App.DB
	dbConfig := config.Config.DB
	if driver == ORACLE_DB {
		connStr := go_ora.BuildUrl(dbConfig.Oracle.Url, dbConfig.Oracle.Port, dbConfig.Oracle.ServiceName, dbConfig.Oracle.Username, dbConfig.Oracle.Password, nil)
		connOracle, err = go_ora.NewConnection(connStr)
		if err != nil {
			logger.Logger.Error("connect to db", "error", err.Error())
			return err
		}
		err = connOracle.Open()
		if err != nil {
			logger.Logger.Error("connect to db", "error", err.Error())
			return err
		}
		logger.Logger.Info("connect to db", "status", constant.SUCCESS)
		return nil
	} else if driver == SQLSERVER_DB {
		// TODO : connect db
		dsn := "sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm"
		connSqlserver, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
		return nil
	} else {
		return fmt.Errorf("driver database not support: %v", driver)
	}
}

func CloseDB() {
	driver := config.Config.App.DB
	if driver == ORACLE_DB {
		err := connOracle.Close()
		if err != nil {
			logger.Logger.Error("close db connection", "error", err.Error())
		}
	}
	if driver == SQLSERVER_DB {
		// TODO : disconnect db

	}
}
