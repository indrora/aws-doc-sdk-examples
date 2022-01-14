package api

import (
	"example.aws/gov2/dynamodb/crud-service/db"
	"github.com/gofiber/fiber/v2"
)

type CreateRequest struct {
	Url string `json:"url"`
}
type CreateResponse struct {
	Id        string `json:"id"`
	DeleteKey string `json:"deleteKey`
}

func CreateLink(c *fiber.Ctx) error {
	request := CreateRequest{}

	if err := c.BodyParser(request); err != nil {
		return err
	}

	link, err := db.CreateLink(request.Url)
	if err != nil {
		return err
	}

	// add to the database

	return c.JSON(*link)

}
