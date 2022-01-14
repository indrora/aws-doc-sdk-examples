package api

import (
	"github.com/gofiber/fiber/v2"
)

type DeleteRequest struct {
	Id        string `json:"id"`
	DeleteKey string `json:"deleteKey"`
}

type DeleteResponse struct {
	Id     string `json:"id"`
	Result bool   `json:"success"`
}

func DeleteKey(c *fiber.Ctx) error {
	request := DeleteRequest{}
	if err := c.BodyParser(request); err != nil {
		return err
	}

	return c.JSON(DeleteResponse{
		Id:     request.Id,
		Result: true,
	})
}
