package main

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	err := godotenv.Load()
	if err != nil {
		log.Info().AnErr("error", err).Msg("Failed to load dotenv!")
	}

	awsConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic("Failed to get AWS configuration")
	}

	tableName, exists := os.LookupEnv("DB_TABLENAME")

	if !exists {
		log.Fatal().Msg("Could not get table name from environment variable DB_TABLENAME!")
	}

	ddb := dynamodb.NewFromConfig(awsConfig)
	log.Info().Str("desiredTableName", tableName).Msg("Checking to see if the table specified exists")

	ddbListTablesResult, err := ddb.ListTables(context.Background(), &dynamodb.ListTablesInput{})

	if err != nil {
		log.Err(err).Msg("Failed to get list of tables... Is the configuration correct?")
		os.Exit(-1)
	}

	if !contains(ddbListTablesResult.TableNames, tableName) {
		log.Info().Msg("Making the table")
		ddbCreate, err := ddb.CreateTable(context.Background(),
			&dynamodb.CreateTableInput{
				AttributeDefinitions: []ddbtypes.AttributeDefinition{
					{AttributeName: aws.String("Id"), AttributeType: ddbtypes.ScalarAttributeTypeS},
					{AttributeName: aws.String("Uri"), AttributeType: ddbtypes.ScalarAttributeTypeS},
					{AttributeName: aws.String("Email"), AttributeType: ddbtypes.ScalarAttributeTypeS},
					{AttributeName: aws.String("DeleteKey"), AttributeType: ddbtypes.ScalarAttributeTypeS},
					{AttributeName: aws.String("createdOn"), AttributeType: ddbtypes.ScalarAttributeTypeS},
					{AttributeName: aws.String("NumHits"), AttributeType: ddbtypes.ScalarAttributeTypeN},
				},
				KeySchema: []ddbtypes.KeySchemaElement{
					{AttributeName: aws.String("Id"), KeyType: ddbtypes.KeyTypeHash},
				},
				TableName: &tableName,
			})
		if err != nil {
			log.Err(err).Msg("Failed to create table...")
		} else {
			log.Info().Str("tableArn", *ddbCreate.TableDescription.TableArn).Msg("Created table")
		}
	}

	emailaddr, exists := os.LookupEnv("REPORT_EMAIL_ADDR")

	if !exists {
		log.Panic().Msg("You must supply an email address as REPORT_EMAIL_ADDR in env")
	}

	sesClient := sesv2.NewFromConfig(awsConfig)

	// see if the identity has already been verified

	getEmailIdentResponse, err := sesClient.GetEmailIdentity(context.Background(), &sesv2.GetEmailIdentityInput{
		EmailIdentity: &emailaddr,
	})

	if err != nil {
		log.Info().Msg("Didn't find it. ")
	} else if getEmailIdentResponse.VerifiedForSendingStatus == false {
		log.Info().Msg("Check your email, or remove the identity and try again.")
	}

	sesClient.CreateEmailIdentity(context.Background(), &sesv2.CreateEmailIdentityInput{
		EmailIdentity: &emailaddr,
	})

	log.Info().Msg("Check your email for the identity verification email")

}

func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
