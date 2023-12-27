package server

import (
	"context"
	"tg-bot-golang/internal/appmodels.go"
)

type productservice interface {
	Handle(ctx context.Context, code int) (*appmodels.Product, error)
	SetProperties(id, value, pref string) error
}

type userservice interface {
	CreateUser(ctx context.Context, name string) error
	DeleteUser(ctx context.Context, name string) error
	CheckUser(ctx context.Context, name string, tgID string) (string, error)
	GetUserByID(ctx context.Context, name string) (string, error)
	CheckAdmin(ctx context.Context, id string) (bool, error)
	GetAllUsers(ctx context.Context) ([]string, error)
}

type point interface {
	Add(ctx context.Context, user, value string)
	Get(ctx context.Context, user string) string
	Del(ctx context.Context, user string)
}

type axelot interface {
	GetRemains(key string) ([]appmodels.Remains, error)
	GetOrder(id string) ([]appmodels.Order, error)
}
