package controller

import (
	"app-connector/config"
	"app-connector/constant"
	"app-connector/logger"
	"app-connector/model"
	"app-connector/util"
	"database/sql/driver"
	"fmt"
	"time"

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
			return err
		}
		err = ConnOracle.Open()
		if err != nil {
			return err
		}
		logger.Info("event: connect_db, status:", constant.SUCCESS, ", driver:", Driver)
		return nil
	} else if Driver == SQLSERVER_DB {
		dsn := dbConfig.Sqlserver.Url
		ConnSqlserver, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
		if err != nil {
			return err
		}
		logger.Info("event: connect_db, status:", constant.SUCCESS, ", driver:", Driver)
		return nil
	} else {
		return fmt.Errorf("driver database not support: %v", Driver)
	}
}

func CloseDB() {
	if Driver == ORACLE_DB {
		err := ConnOracle.Close()
		if err != nil {
			logger.Error("event: disconnect_db, status: error, msg:", err.Error())
		}
	}
}

func CreateStmt(table string) {
	if Driver == ORACLE_DB {
		// stmt
		sql := fmt.Sprintf(`UPDATE %s SET Status = :1, SFrom = :2, "TimeStamp" = :3 WHERE ALM_Tag = :4`, table)
		stmt = go_ora.NewStmt(sql, ConnOracle)
	}
}

func CloseStmt() {
	if Driver == ORACLE_DB {
		_ = stmt.Close()
	}
}

func UpdateStatus(status, sform, tag string, t time.Time, table string) error {
	if Driver == ORACLE_DB {
		tsz := util.Timestamptz(t)
		_, err := stmt.Exec([]driver.Value{status, sform, tsz, tag})
		if err != nil {
			return err
		}
	} else if Driver == SQLSERVER_DB {
		whcl := "ALM_Tag = ?"
		data := map[string]interface{}{"Status": status, "SFrom": sform, "TimeStamp": util.FormatDatetime(t)}
		tx := ConnSqlserver.Table(table).Where(whcl, tag).Updates(data)
		if tx.Error != nil {
			return tx.Error
		}
	} else {
		return fmt.Errorf("driver database not support: %v", Driver)
	}
	return nil
}

func FindByTag(tag string, table string) (string, error) {
	if Driver == ORACLE_DB {
		return "", fmt.Errorf("driver database dose not implement: %v", Driver)
	} else if Driver == SQLSERVER_DB {
		var data model.KpiTag
		whcl := "ALM_Tag = ?"
		tx := ConnSqlserver.Table(table).Where(whcl, tag).First(&data)
		if tx.Error != nil {
			return "", tx.Error
		}
		return data.Status, nil
	} else {
		return "", fmt.Errorf("driver database not support: %v", Driver)
	}
}
