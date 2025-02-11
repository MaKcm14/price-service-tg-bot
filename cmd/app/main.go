package main

import "github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/app"

func main() {
	service := app.NewService()
	service.Run()
}
