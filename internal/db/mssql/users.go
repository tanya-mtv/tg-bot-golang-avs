package mssql

import (
	"context"
	"database/sql"
	"fmt"
	"tg-bot-golang/internal/config"
	"tg-bot-golang/internal/logger"
)

type UsersStorage struct {
	log       logger.Logger
	cfg       *config.Config
	sqlClient *sql.DB
}

func NewmsUsersStorage(log logger.Logger, cfg *config.Config, sqlClient *sql.DB) *UsersStorage {
	return &UsersStorage{log: log, cfg: cfg, sqlClient: sqlClient}
}

func (u *UsersStorage) CheckUser(ctx context.Context, name string, tgID string) (string, error) {
	var getname = ""
	query := fmt.Sprintf("IF EXISTS (SELECT TOP(1) 1 FROM TGUsers WHERE Name = '%s') BEGIN UPDATE  TGUsers SET TGID = '%s' OUTPUT INSERTED.Name WHERE Name = '%s' END;", name, tgID, name)
	err := u.sqlClient.QueryRow(query).Scan(&getname)
	if err != nil {
		u.log.Errorf(err.Error())
		return "", err
	}

	return getname, nil
}

func (u *UsersStorage) GetUserByID(ctx context.Context, id string) (string, error) {
	var name = ""
	query := fmt.Sprintf("SELECT  TOP 1 name FROM TGUsers WHERE TGID = '%s' AND IsActive = 1", id)
	err := u.sqlClient.QueryRow(query).Scan(&name)
	if err != nil {
		u.log.Errorf(err.Error())
		return "", err
	}

	return name, nil
}

func (u *UsersStorage) CheckAdmin(ctx context.Context, id string) (bool, error) {
	var isAdm = false
	query := fmt.Sprintf("SELECT  TOP 1 IsAdmin FROM TGUsers WHERE TGID = '%s'", id)
	err := u.sqlClient.QueryRow(query).Scan(&isAdm)
	if err != nil {
		u.log.Errorf(err.Error())
		return isAdm, err
	}

	return isAdm, nil
}

func (u *UsersStorage) CreateUser(ctx context.Context, name string) error {
	query := fmt.Sprintf("IF NOT EXISTS (SELECT TOP(1) 1 FROM TGUsers WHERE Name = '%s') BEGIN INSERT INTO TGUsers (ModifiedDate, Name, IsAdmin, isActive) VALUES(GETDATE(), '%s', 0, 1) END", name, name)
	_, err := u.sqlClient.Exec(query)
	if err != nil {
		u.log.Errorf(err.Error())
		return err
	}

	return nil
}

func (u *UsersStorage) GetAllUsers(ctx context.Context) ([]string, error) {
	var users []string

	query := "SELECT  Name from  TGUsers WHERE IsActive = 1"
	rows, err := u.sqlClient.Query(query)
	if err != nil {
		u.log.Errorf(err.Error())
		return users, err
	}
	for rows.Next() {
		var user string
		err = rows.Scan(&user)
		if err != nil {
			u.log.Errorf(err.Error())
			return users, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (u *UsersStorage) DeleteUser(ctx context.Context, name string) error {
	query := fmt.Sprintf("UPDATE  TGUsers set IsActive = 0 WHERE Name = '%s'", name)
	_, err := u.sqlClient.Exec(query)
	if err != nil {
		u.log.Errorf(err.Error())
		return err
	}

	return nil
}
