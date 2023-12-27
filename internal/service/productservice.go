package service

import (
	"context"
	"strconv"
	"tg-bot-golang/internal/appmodels.go"
	"tg-bot-golang/internal/config"
	"tg-bot-golang/internal/logger"
)

type ProductService struct {
	log       logger.Logger
	cfg       *config.Config
	sqlRepo   msSQL
	redisRepo cash
	raecRepo  raec
}

func NewProductService(log logger.Logger, cfg *config.Config, sqlRepo msSQL, redisRepo cash, raecRepo raec) *ProductService {
	return &ProductService{
		log:       log,
		cfg:       cfg,
		sqlRepo:   sqlRepo,
		redisRepo: redisRepo,
		raecRepo:  raecRepo,
	}
}

func (q *ProductService) Handle(ctx context.Context, code int) (*appmodels.Product, error) {
	product := &appmodels.Product{}

	if product, err := q.redisRepo.GetProduct(ctx, strconv.Itoa(code)); err == nil && product != nil {
		q.log.Infof("Get Product %+v with key %s", product, code)
		return product, nil
	}

	good_id, err := q.sqlRepo.GetProductById(code)

	if err != nil {
		q.log.Errorf("Can't get good_id from avs base")
		return &appmodels.Product{}, err
	}

	product = q.raecRepo.GetInfoRaec(good_id)

	if product == (&appmodels.Product{}) {
		q.log.Infof("Can't find product with key %s", product, code)
		return &appmodels.Product{}, err
	}

	q.redisRepo.PutProduct(ctx, strconv.Itoa(code), product)
	return product, nil
}

func (q *ProductService) SetProperties(code, value, pref string) error {
	return q.sqlRepo.SetProperties(code, value, pref)
}
