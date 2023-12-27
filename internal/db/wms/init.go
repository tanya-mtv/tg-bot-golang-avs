package wms

import (
	"database/sql"
	"fmt"
	"log"
	"tg-bot-golang/internal/config"

	_ "github.com/alexbrainman/odbc"
)

func NewAxelotConn(cfg *config.ConfigAxelot) *sql.DB {

	connstring := fmt.Sprintf("driver={%s};SERVER=%s;UID=%s;PWD=%s;DATABASE=%s", cfg.DriverName, cfg.Server, cfg.User, cfg.Password, cfg.DSN)

	db, err := sql.Open(cfg.Driver, connstring)

	if err != nil {
		log.Fatal("Error creating connection pool: " + err.Error())
	}

	return db
}
