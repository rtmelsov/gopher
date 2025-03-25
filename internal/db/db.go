package db

import (
	"github.com/rtmelsov/GopherMart/internal/config"
	"github.com/rtmelsov/GopherMart/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
)

type DB struct {
	db   *gorm.DB
	conf config.ConfigI
	mu   sync.RWMutex
}

type DBI interface {
	ErrorHandler(err string, code int) *models.Error

	Register(value *models.DBUser) (*models.DBUser, *models.Error)
	Login(value *models.DBUser) (*models.DBUser, *models.Error)
	GetUser(userID uint) (*models.DBUser, *models.Error)

	PostOrders(order *models.DBOrder) *models.Error
	GetOrders(id *uint) (*[]models.DBOrder, *models.Error)
	GetOrder(orderNumber string) (*models.DBOrder, *models.Error)

	GetBalance(id *uint) (*models.Balance, *models.Error)

	GetWithdrawal(orderNumber string) (*models.DBWithdrawal, *models.Error)

	GetWithdrawals(id *uint) (*[]models.DBWithdrawal, *models.Error)
	PostOrderWithDraw(withdrawal *models.DBWithdrawal) *models.Error
}

func (db *DB) ErrorHandler(err string, code int) *models.Error {
	return db.conf.ErrorHandler("DB", err, code)
}

func GetDB(conf config.ConfigI) (DBI, error) {
	//dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(conf.GetEnvVariables().DataBaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.Migrator().DropTable(&models.DBUser{}, &models.DBOrder{}, &models.DBWithdrawal{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&models.DBUser{}, &models.DBOrder{}, &models.DBWithdrawal{})
	if err != nil {
		return nil, err
	}

	return &DB{
		db:   db,
		conf: conf,
		mu:   sync.RWMutex{},
	}, nil
}
