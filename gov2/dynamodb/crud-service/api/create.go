package api

import (
	"example.aws/gov2/dynamodb/crud-service/db"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
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
	request := new(CreateRequest)
	link := db.Link{}
	err := error(nil)

	c.Accepts("json", "text")
	if err = c.BodyParser(request); err != nil {
		log.Err(err).Msg("Failed to handle body!")
		goto failure
	}
	log.Debug().Str("link", request.Url).Str("email", request.Email).Msg("Going to create a link...")

	if link, err = db.CreateLink(request.Url, request.Email); err != nil {
		log.Err(err).Msg("Failed to create link.")
		goto failure
	}

	// add to the database

	if err = (db.DB).Add(link); err != nil {
		log.Err(err).Msg("Failed to add link to database")
		goto failure
	}

	c.Status(200)
	return c.JSON(CreateResponse{
		Success:   true,
		Error:     "no error",
		Id:        link.Id,
		DeleteKey: link.DeleteKey,
	})

failure:
	c.Status(400)
	return c.JSON(CreateResponse{
		Success: false,
		Error:   "internal error",
	})

}
