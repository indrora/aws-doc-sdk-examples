package main

import (
	"context"
	"net/http"
	"os"

	"embed"

	"example.aws/gov2/dynamodb/crud-service/api"
	"example.aws/gov2/dynamodb/crud-service/db"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	DEBUG bool
)

//go:embed web/*
var embedDirStatic embed.FS

func main() {
	app := fiber.New()

	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	err := godotenv.Load()
	if err != nil {
		log.Info().AnErr("error", err).Msg("Failed to load dotenv!")
	}

	awsConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic("Failed to get AWS configuration")
	}

	if DEBUG = os.Getenv("DEBUG") != ""; DEBUG {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		log.Info().Msg("Hello, debug world!")
		_ = godotenv.Load("env.debug")
	}

	// Make sure that the middleware sets up our dynamodb connection.

	tableName, exists := os.LookupEnv("DB_TABLENAME")
	if !exists {
		log.Fatal().Msg("No table defined by DB_TABLENAME")
	}

	dbConn := db.GetDDBConnection(tableName, awsConfig)

	db.DB = dbConn

	group := app.Group("/api")

	group.Put("/link", api.CreateLink)
	group.Delete("/link/:id", api.DeleteLink)
	group.Get("/link/:id", api.GetLinkStats)
	group.Get("/link/report", api.GetLinks)

	app.Get("/go/:id", api.DoRedirect)

	app.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(embedDirStatic),
		PathPrefix: "web",
		Browse:     true,
	}))

	// Access file "image.png" under `static/` directory via URL: `http://<server>/static/image.png`.
	// Without `PathPrefix`, you have to access it via URL:
	// `http://<server>/static/static/image.png`.
	app.Use("/static", filesystem.New(filesystem.Config{
		Root:       http.FS(embedDirStatic),
		PathPrefix: "web/static",
		Browse:     true,
	}))

	app.Listen(":8080")
}
