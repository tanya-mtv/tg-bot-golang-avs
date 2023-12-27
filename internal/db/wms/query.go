package wms

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"tg-bot-golang/internal/appmodels.go"
	"tg-bot-golang/internal/config"
	"tg-bot-golang/internal/logger"
)

type AxelotStorage struct {
	log       logger.Logger
	cfg       *config.Config
	sqlClient *sql.DB
}

func NewmsAxelotStorage(log logger.Logger, cfg *config.Config, sqlClient *sql.DB) *AxelotStorage {
	return &AxelotStorage{log: log, cfg: cfg, sqlClient: sqlClient}
}

func (a *AxelotStorage) GetRemains(key string) ([]appmodels.Remains, error) {
	var remainslist []appmodels.Remains

	code, err := strconv.Atoi(key)
	if err != nil {
		code = 0
	}

	query := fmt.Sprintf(`SELECT DISTINCT
        T2._Description name,
        T2._Code code,
        T3._fld1136 cell,
        T5._Description eh,
        T1._Fld2966 count
        FROM _InfoRg2958 T1 WITH(NOLOCK)
            LEFT OUTER JOIN _Reference65 T2 WITH(NOLOCK)
                ON (T1._Fld2959RRef = T2._IDRRef)
                    LEFT OUTER JOIN _Reference85 T3 WITH(NOLOCK)
                        ON (T1._Fld2965RRef = T3._IDRRef)
                    LEFT OUTER JOIN _Reference68 T4 WITH(NOLOCK)
                        ON (T1._Fld2964RRef = T4._IDRRef)
                    LEFT OUTER JOIN _Reference58 T5 WITH(NOLOCK)
                        ON (T1._Fld2960RRef = T5._IDRRef)
                    LEFT JOIN _InfoRg4119 T6 WITH(NOLOCK)
                    ON ((T1._Fld2959RRef = T6._Fld4121RRef) AND (T1._Fld2960RRef = T6._Fld4122RRef))
        WHERE (T1._Fld2966 <> 0.0) AND (T3._fld1136 = '%s' OR T6._Fld4120 = '%s' OR T2._code = %09d)`, key, key, code)

	rows, err := a.sqlClient.Query(query)
	if err != nil {
		a.log.Errorf(err.Error())
		return remainslist, err
	}
	defer rows.Close()

	for rows.Next() {
		str := appmodels.Remains{}
		err := rows.Scan(&str.Name, &str.Code, &str.Cell, &str.EH, &str.Count)
		if err != nil {
			a.log.Errorf(err.Error())
		}
		remainslist = append(remainslist, str)
	}

	return remainslist, nil
}

func (a *AxelotStorage) GetOrder(id string) ([]appmodels.Order, error) {
	order := make([]appmodels.Order, 0)
	query := fmt.Sprintf(
		`SELECT
            T5._Description executor,
            T4._Description zone
        FROM _InfoRg3490 T1 WITH(NOLOCK)
        LEFT OUTER JOIN _Document102 T2 WITH(NOLOCK)
            ON (T1._Fld3504RRef = T2._IDRRef)
        LEFT OUTER JOIN _Document111 T3 WITH(NOLOCK)
        LEFT OUTER JOIN _Reference74 T4 WITH(NOLOCK)
            ON (T3._Fld1644RRef = T4._IDRRef)
            ON (T1._Fld3491RRef = T3._IDRRef)
        LEFT OUTER JOIN _Reference73 T5 WITH(NOLOCK)
            ON (T1._Fld3507RRef = T5._IDRRef)
        WHERE (T2._Fld4283 = '%s')
        GROUP BY T5._Description, T4._Description`, id)

	rows, err := a.sqlClient.Query(query)
	if err != nil {
		a.log.Errorf(err.Error())
		return order, err
	}
	defer rows.Close()

	for rows.Next() {
		str := appmodels.Order{}
		err := rows.Scan(&str.Executor, &str.Zone)
		if err != nil {
			log.Fatal(err)
		}
		order = append(order, str)
	}
	return order, nil

}
