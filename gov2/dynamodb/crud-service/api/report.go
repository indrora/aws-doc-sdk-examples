package api

import (
	"bytes"
	"context"
	"errors"
	"text/template"

	"example.aws/gov2/dynamodb/crud-service/db"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/gofiber/fiber/v2"
)

var emailTemplate *template.Template

func init() {

	//go:embed email_template.txt
	var email_template string
	var err error

	emailTemplate, err = template.New("emailBody").Parse(email_template)

	if err != nil {
		panic("bad email template syntax.")
	}
}

const INVALIDEMAIL = "###INVALIDEMAIL@@INVALID.INVALID##"

func GetLinks(c *fiber.Ctx) error {

	conf, err := config.LoadDefaultConfig(context.Background())

	if err != nil {
		return SendJSONResponse(c, 500, ResponseBase{
			Error:  "internal error",
			Result: false,
		})
	}

	email := c.Params("email", INVALIDEMAIL)
	if email == INVALIDEMAIL {
		return SendJSONError(c, errors.New("No or invalid email specified"))
	}

	sesClient := ses.NewFromConfig(conf)

	emailBuffer := new(bytes.Buffer)

	emailTemplate.Execute(emailBuffer, map[string]interface{}{
		"baseurl": c.BaseURL(),
		"links":   []db.Link{},
	})

	_, err = sesClient.SendEmail(context.Background(), &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{},
		},
		Message:              &types.Message{},
		Source:               new(string),
		ConfigurationSetName: new(string),
		ReplyToAddresses:     []string{},
		ReturnPath:           new(string),
		ReturnPathArn:        new(string),
		SourceArn:            new(string),
		Tags:                 []types.MessageTag{},
	})

	if err != nil {
		return SendJSONError(c, errors.New("Failed to send email"))
	}

	return SendJSONResponse(c, 200, ResponseBase{
		Error:  "No error",
		Result: true,
	})

}
