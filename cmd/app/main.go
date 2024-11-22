package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"go.uber.org/zap"

	_ "github.com/xoticdsign/auf-citaty/docs"
	"github.com/xoticdsign/auf-citaty/internal/cache"
	"github.com/xoticdsign/auf-citaty/internal/database"
	"github.com/xoticdsign/auf-citaty/internal/middleware"
	"github.com/xoticdsign/auf-citaty/internal/routes"
	"github.com/xoticdsign/auf-citaty/utils/errhandling"
	"github.com/xoticdsign/auf-citaty/utils/logging"
)

// General description
//
//	@title						Auf Citaty API
//	@version					1.0.0
//	@description				TODO
//	@contact.name				xoti$
//	@contact.url				https://t.me/xoticdsign
//	@contact.email				xoticdollarsign@outlook.com
//	@license.name				MIT
//	@license.url				https://mit-license.org/
//	@host						127.0.0.1:8080
//	@BasePath					/
//	@produce					json
//	@schemes					http
//
//	@securitydefinitions.apikey	KeyAuth
//	@in							query
//	@name						auf-citaty-key
func main() {
	godotenv.Load()

	err := cache.RunRedis()
	if err != nil {
		log.Fatal(err)
	}

	err = logging.RunZap()
	if err != nil {
		log.Fatal(err)
	}

	err = database.RunGORM()
	if err != nil {
		log.Fatal(err)
	}

	appName := os.Getenv("APP_NAME")
	addr := os.Getenv("SERVER_ADDRESS")

	app := fiber.New(fiber.Config{
		ServerHeader:  appName,
		StrictRouting: true,
		CaseSensitive: true,
		ReadTimeout:   time.Second * 20,
		WriteTimeout:  time.Second * 20,
		ErrorHandler:  errhandling.ErrorHandler,
		AppName:       appName,
	})

	middleware.GetMiddleware(app)
	routes.GetRoutes(app)

	logging.Logger.Info(
		"Сервер запущен",
		zap.String("Address", addr),
	)

	err = app.Listen(addr)
	if err != nil {
		log.Fatal(err)
	}
}
