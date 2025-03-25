package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rtmelsov/GopherMart/internal/config"
	"github.com/rtmelsov/GopherMart/internal/db"
	"github.com/rtmelsov/GopherMart/internal/handlers"
	"github.com/rtmelsov/GopherMart/internal/middleware"
	"github.com/rtmelsov/GopherMart/internal/repository"
	"github.com/rtmelsov/GopherMart/internal/services"
)

func main() {
	// получили все перемененные окружения
	c, localError := config.InitConfig()
	if localError != nil {
		panic(localError)
	}

	// init gin router
	g := gin.Default()

	// подключаемся к базе данных
	d, err := db.GetDB(c)
	if err != nil {
		panic(err)
	}

	// передаем методы базы данных в репозитории чтобы мы могли работать с ним через него (изолируем данные базы данных кроме методов)
	r := repository.GetRepository(c, d)

	// передаем методы репозитории в сервисы чтобы к ним обращаться через них, а не напрямую (изолируем данные репозитории кроме методов)
	s := services.NewService(c, r)

	// передаем в ручки все методы сервиса (изолируем данные сервиса кроме методов)
	h := handlers.NewHandler(c, s)

	// middleware
	g.Use(gin.Logger())

	// добавили в начале /api/user
	clientAPI := g.Group(c.GetEnvVariables().RootURL)

	clientAPI.POST("/register", h.Register)
	clientAPI.POST("/login", h.Login)

	protected := g.Group(c.GetEnvVariables().RootURL)
	protected.Use(middleware.Auth(c, d))
	{
		protected.POST("/orders", h.PostOrders)
		protected.GET("/orders", h.GetOrders)

		protected.GET("/balance", h.GetBalance)

		protected.POST("/balance/withdraw", h.PostBalanceWithdraw)
		protected.GET("/withdrawals", h.GetWithdrawals)
	}

	if err := g.Run(c.GetEnvVariables().RunAddress); err != nil {
		return
	}
}
