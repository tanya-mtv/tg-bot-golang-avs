package mssql

import (
	"database/sql"
	"fmt"
	"log"
	"tg-bot-golang/internal/config"

	_ "github.com/alexbrainman/odbc"
)

func NewMSSQLDB(cfg *config.ConfigMSSQL) *sql.DB {

	// connstring := fmt.Sprintf("driver={%s};SERVER=%s;UID=%s;PWD=%s;DATABASE=%s", cfg.DriverName, cfg.Server, cfg.User, cfg.Password, cfg.DSN)
	connstring := fmt.Sprintf("driver={%s};SERVER=%s;UID=%s;PWD=%s;DATABASE=%s", cfg.DriverName, cfg.Server, cfg.User, cfg.Password, cfg.DSN)
	// connstring := "driver={ODBC Driver 13 for SQL Server};SERVER=sql-cl02.avselectro.ru\\sec;UID=1c_telegram;PWD=1La0cVG136ofQ;DATABASE=AVSIntegrationGate_development"
	db, err := sql.Open(cfg.Driver, connstring)

	if err != nil {
		log.Fatal("Error creating connection pool: " + err.Error())
	}

	return db
}
