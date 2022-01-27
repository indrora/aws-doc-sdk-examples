package api

import (
	"example.aws/gov2/dynamodb/crud-service/db"
	"github.com/gofiber/fiber/v2"
)

type CreateRequest struct {
	Url   string `json:"url"`
	Email string `json:"email"`
}
type CreateResponse struct {
	Success   bool   `json:"success"`
	Error     string `json:"error"`
	Id        string `json:"id,omitempty"`
	DeleteKey string `json:"deleteKey,omitempty"`
}

func CreateLink(c *fiber.Ctx) error {
	request := CreateRequest{}
	link := db.Link{}
	err := error(nil)

	if err = c.BodyParser(request); err != nil {
		goto failure
	}
	if link, err = db.CreateLink(request.Url, request.Email); err != nil {
		goto failure
	}

	// add to the database

	if err = (db.DB).Add(link); err != nil {
		goto failure
	}

	return c.JSON(link)

failure:
	c.Status(400)
	return c.JSON(CreateResponse{
		Success: false,
		Error:   err.Error(),
	})

}
