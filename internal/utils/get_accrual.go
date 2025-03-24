package utils

import (
	"encoding/json"
	"fmt"
	"github.com/rtmelsov/GopherMart/internal/config"
	"github.com/rtmelsov/GopherMart/internal/models"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
)

func GetAccrual(conf config.ConfigI, num string) (*models.Accrual, *models.Error) {
	var order models.Accrual
	resp, err := http.Get(fmt.Sprintf("%s/api/orders/%s", conf.GetEnvVariables().AccrualSystemAddress, num))
	if err != nil {
		conf.GetLogger().Error("error to get resp from accrual", zap.Error(err))
		return nil, &models.Error{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		}
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close request body: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	conf.GetLogger().Info("body from accrual", zap.String("body", string(body)))
	err = json.Unmarshal(body, &order)
	if err != nil {
		return nil, &models.Error{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		}
	}

	if order.Accrual == 0 {
		order.Status = "PROCESSING"
	}

	return &order, nil
}
