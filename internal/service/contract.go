package service

import (
	"context"
	"tg-bot-golang/internal/appmodels.go"
)

type msSQL interface {
	GetProductById(good_id int) (int, error)
	SetProperties(id, value, pref string) error
}

// type users interface {
// 	recognizeUser()
// }

type cash interface {
	PutProduct(ctx context.Context, key string, product *appmodels.Product)
	GetProduct(ctx context.Context, key string) (*appmodels.Product, error)
}

type raec interface {
	GetInfoRaec(code int) *appmodels.Product
}
