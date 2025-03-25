package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rtmelsov/GopherMart/internal/config"
	"github.com/rtmelsov/GopherMart/internal/models"
	"io"
	"log"
	"net/http"
)

type Order struct {
	Order string `json:"order"`
}

func GetAccrual(conf config.ConfigI, num string) (*models.Accrual, *models.Error) {
	var order models.Accrual

	resp, err := http.Get(fmt.Sprintf("%s/api/orders/%s", conf.GetEnvVariables().AccrualSystemAddress, num))
	if err != nil {
		return nil, conf.ErrorHandler("EXTERNAL: error to get resp from accrual", err.Error(), http.StatusInternalServerError)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close request body: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, conf.ErrorHandler("EXTERNAL", err.Error(), http.StatusInternalServerError)
	}

	err = json.Unmarshal(body, &order)
	if err != nil {
		return nil, conf.ErrorHandler("EXTERNAL", err.Error(), http.StatusInternalServerError)
	}

	if order.Accrual == 0 {
		order.Status = "PROCESSING"
	}

	return &order, nil
}

func PostAccrual(conf config.ConfigI, num string) *models.Error {
	var order = Order{
		Order: num,
	}
	reqBody, err := json.Marshal(order)
	if err != nil {
		return conf.ErrorHandler("EXTERNAL", err.Error(), http.StatusBadRequest)
	}

	resp, err := http.Post(fmt.Sprintf("%s/api/orders", conf.GetEnvVariables().AccrualSystemAddress), "application/json", bytes.NewBuffer(reqBody))

	if err != nil {
		return conf.ErrorHandler("EXTERNAL", err.Error(), http.StatusUnprocessableEntity)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close request body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusConflict && resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return conf.ErrorHandler("EXTERNAL", "error while checking the order number", http.StatusUnprocessableEntity)
	}

	return nil
}
