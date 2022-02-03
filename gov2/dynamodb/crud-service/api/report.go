package api

import (
	"bytes"
	"context"
	"errors"
	"os"
	"text/template"

	_ "embed"

	"example.aws/gov2/dynamodb/crud-service/db"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

var emailBodyTemplate *template.Template
var emailFromAddress string

//go:embed email_template.txt
var email_template string

func init() {

	var err error

	emailFromAddress = ""
	emailBodyTemplate, err = template.New("emailBody").Parse(email_template)

	if err != nil {
		panic("bad email template syntax.")
	}

}

type GetLinksRequest struct {
	Email string `json:"email"`
}

const INVALIDEMAIL = "###INVALIDEMAIL@@INVALID.INVALID##"

func GetLinks(c *fiber.Ctx) error {

	if emailFromAddress == "" {
		emailFromAddress = os.Getenv("REPORT_EMAIL_ADDR")
	}

	conf, err := config.LoadDefaultConfig(context.Background())

	if err != nil {
		return SendJSONResponse(c, 500, ResponseBase{
			Error:  "internal error",
			Result: false,
		})
	}

	requestbody := new(GetLinksRequest)

	err = c.BodyParser(requestbody)
	if err != nil {
		log.Fatal().AnErr("err", err).Msg("Failed to parse body")
		return err
	}

	email := requestbody.Email

	log.Debug().Str("email", requestbody.Email).Msg("Sending report")

	if email == INVALIDEMAIL {
		return SendJSONError(c, errors.New("no or invalid email specified"))
	}

	links := (db.DB).ListByEmail(email)

	for _, link := range links {
		log.Debug().Str("id", link.Id).
			Str("email", link.Email).
			Str("deleteKey", link.DeleteKey).Msg("link debug")
	}

	sesClient := ses.NewFromConfig(conf)

	emailBodyBuffer := new(bytes.Buffer)

	emailBodyTemplate.Execute(emailBodyBuffer, map[string]interface{}{
		"baseurl": c.BaseURL(),
		"links":   links,
	})

	emailBodyString := emailBodyBuffer.String()
	log.Debug().Str("emailbody", emailBodyString).Msg("debug")

	_, err = sesClient.SendEmail(context.Background(), &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{
				email,
			},
		},
		Message: &types.Message{
			Body: &types.Body{
				Text: &types.Content{
					Data:    &emailBodyString,
					Charset: aws.String("utf-8"),
				},
			},
			Subject: &types.Content{
				Data:    aws.String("Your shortened links"),
				Charset: aws.String("utf-8"),
			},
		},
		Source: &emailFromAddress,
	})

	if err != nil {
		log.Err(err).Msg("Failed to send email")
		return SendJSONError(c, errors.New("failed to send email"))
	}

	return SendJSONResponse(c, 200, ResponseBase{
		Error:  "No error",
		Result: true,
	})

}
