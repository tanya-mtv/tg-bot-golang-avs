package httpclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"tg-bot-golang/internal/appmodels.go"
	"tg-bot-golang/internal/config"
	"tg-bot-golang/internal/logger"
)

const urlRaec = "http://catalog.raec.su/api/product/"

type RaecAPI struct {
	log        logger.Logger
	cfg        *config.Config
	httpClient *http.Client
}

func NewRaecRepository(log logger.Logger, cfg *config.Config) *RaecAPI {
	httpCl := &http.Client{}
	return &RaecAPI{log: log, cfg: cfg, httpClient: httpCl}
}

func (r *RaecAPI) GetInfoRaec(code int) *appmodels.Product {
	produrl := fmt.Sprintf("%s%s", urlRaec, strconv.Itoa(code))
	req, err := http.NewRequest("GET", produrl, nil)
	if err != nil {
		r.log.Errorf("Can't get raec request", err)
	}

	req.Header = http.Header{
		r.cfg.RaecKey: {r.cfg.RaecValue},
	}

	res, err := r.httpClient.Do(req)
	if err != nil {
		r.log.Errorf("Can't get raec data", err)
	}
	defer res.Body.Close()

	product := &appmodels.Product{}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		r.log.Errorf("can;t read request body", err)
	}

	err = json.Unmarshal([]byte(body), &product)

	if err != nil {
		r.log.Errorf("Can't unmarshal raec request", err)
	}
	return product
}
