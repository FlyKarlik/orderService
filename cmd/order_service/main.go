package main

import (
	"github.com/FlyKarlik/orderService/config"
	"github.com/FlyKarlik/orderService/internal/app/order_service"
	"github.com/FlyKarlik/orderService/pkg/logger"
	"github.com/FlyKarlik/orderService/pkg/validate"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	if err := validate.Validate(cfg); err != nil {
		panic(err)
	}

	logger, err := logger.New(cfg.OrderService.LogLevel)
	if err != nil {
		panic(err)
	}

	orderService := order_service.New(cfg, logger)
	if err := orderService.Start(); err != nil {
		panic(err)
	}
}
