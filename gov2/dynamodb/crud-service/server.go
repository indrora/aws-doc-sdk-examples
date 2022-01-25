package main

import (
	"context"
	"os"

	"example.aws/gov2/dynamodb/crud-service/api"
	"example.aws/gov2/dynamodb/crud-service/db"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	DEBUG bool
)

func main() {
	app := fiber.New()

	awsConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic("Failed to get AWS configuration")
	}

	err = godotenv.Load()

	if err != nil {
		log.Info().AnErr("error", err).Msg("Failed to load dotenv!")
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
	db.SetupDBConnection(dbConn)

	_ = api.GetApi(app)

	app.Listen(":8080")
}
