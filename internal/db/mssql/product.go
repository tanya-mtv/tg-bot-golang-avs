package mssql

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"tg-bot-golang/internal/config"
	"tg-bot-golang/internal/logger"
)

type MSSQLStorage struct {
	log       logger.Logger
	cfg       *config.Config
	sqlClient *sql.DB
}

func NewmsSQLStorage(log logger.Logger, cfg *config.Config, sqlClient *sql.DB) *MSSQLStorage {
	return &MSSQLStorage{log: log, cfg: cfg, sqlClient: sqlClient}
}

func (r *MSSQLStorage) GetProductById(goodID int) (int, error) {
	var raecID int

	query := fmt.Sprintf("SELECT prod.CodeRAEK FROM dbo.Product  prod WHERE prod.code = %d", goodID)
	err := r.sqlClient.QueryRow(query).Scan(&raecID)
	if err != nil {
		r.log.Errorf(err.Error())
		return 0, err
	}

	return raecID, nil
}

func (r *MSSQLStorage) SetProperties(id, value, pref string) error {
	var query string
	var queryGetCode string
	var productID int

	tx, err := r.sqlClient.Begin()
	if err != nil {
		return err
	}
	defer func() {
		err = tx.Rollback()
		if err != nil {
			return
		}
	}()

	switch {
	case strings.Contains(pref, "BarCode"):
		unit := getUnit(pref)

		code := getCode(id)
		query = fmt.Sprintf(`insert into TGUnit (ModifiedDate, TGStatus, Owner, Barcode, Unit)
                                VALUES (GETDATE(), 1, '%s','%s', (SELECT id FROM dbo.UnitQualifier WHERE code = '%s'))`, code, value, unit)

		queryGetCode = fmt.Sprintf(`SELECT id FROM Product WHERE code = '%s'`, code)
	case strings.Contains(pref, "BtnWheight"):

		intvalue, err := strconv.Atoi(value)

		if err != nil {
			r.log.Errorf("Can't convert weight value to int")
			return err
		}
		query = fmt.Sprintf(`insert into TGUnit (ModifiedDate, TGStatus, Owner, Weight, Unit)
                                VALUES (GETDATE(), 1, '%s','%d', (SELECT id FROM dbo.UnitQualifier WHERE code = '796'))`, id, intvalue)
		queryGetCode = fmt.Sprintf("SELECT id FROM Product WHERE code = %s", id)

	case strings.Contains(pref, "BtnVolume"):
		intvalue, err := strconv.Atoi(value)
		if err != nil {
			r.log.Errorf("Can't convert volume value to int")
			return err
		}
		query = fmt.Sprintf(`insert into TGUnit (ModifiedDate, TGStatus, Owner, Volume, Unit)
                                VALUES (GETDATE(), 1, '%s','%d', (SELECT id FROM dbo.UnitQualifier WHERE code = '796'))`, id, intvalue)
		queryGetCode = fmt.Sprintf("SELECT id FROM Product  WHERE code = %s", id)
	}

	stmtCode, err := tx.Prepare(queryGetCode)

	err = stmtCode.QueryRow().Scan(&productID)
	if err != nil {
		r.log.Errorf("Can't select id: " + err.Error())
		return err
	}
	defer stmtCode.Close()

	stmt, err := tx.Prepare(query)
	_, err = stmt.Exec()

	if err != nil {
		r.log.Errorf("Can't insert data: " + err.Error())
		return err
	}
	defer stmt.Close()

	if productID != 0 {
		err = tx.Commit()
		if err != nil {
			return err
		}
	} else {
		return errors.New("UnnownProduct")
	}

	return nil
}

func getUnit(pref string) string {
	var unit string
	switch {
	case strings.Contains(pref, "796/"):
		unit = "796"
	case strings.Contains(pref, "006/"):
		unit = "006"
	case strings.Contains(pref, "123/"):
		unit = "123"
	case strings.Contains(pref, "778/"):
		unit = "778"

	}

	return unit
}

func getCode(pref string) string {
	var code string
	switch {
	case strings.Contains(pref, "796/"):
		code = strings.Replace(pref, "796/", "", -1)
	case strings.Contains(pref, "006/"):
		code = strings.Replace(pref, "006/", "", -1)
	case strings.Contains(pref, "123/"):
		code = strings.Replace(pref, "123/", "", -1)
	case strings.Contains(pref, "778/"):
		code = strings.Replace(pref, "778/", "", -1)
	}

	return code
}
