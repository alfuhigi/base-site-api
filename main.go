package main

import (
	"github.com/joho/godotenv"
	"os"

	"base-site-api/config"

	"github.com/gofiber/fiber"
	log "github.com/sirupsen/logrus"
)

func setupEnv() {
	var logLevel log.Level
	var err error

	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)

	if os.Getenv("GO_ENV") == "testing" {
		err = godotenv.Load(".test.env")
	} else {
		err = godotenv.Load()
	}

	if err != nil {
		log.Fatalf("Fatal while loading env: %s", err)
	}

	if os.Getenv("GO_ENV") == "development" {
		logLevel = log.DebugLevel
	} else {
		logLevel = log.InfoLevel
	}
	log.SetLevel(logLevel)
}

func main() {
	setupEnv()

	c, err := config.New()

	if err != nil {
		log.Fatal(err)
	}

	NewApp(c)
}

func NewApp(c *config.Config) *fiber.App {
	// SETUP APP
	app := fiber.New(&fiber.Settings{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
	})

	configureGlobalMiddleware(app)

	configureAPIRoutes(app, c)

	startServer(app, c)

	return app
}

func startServer(app *fiber.App, c *config.Config) {

	err := app.Listen(c.Constants.ADDRESS)

	if err != nil {
		log.Fatal(err.Error())
	}
}
