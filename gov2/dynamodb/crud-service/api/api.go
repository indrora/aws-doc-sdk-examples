package api

import "github.com/gofiber/fiber/v2"

type ResponseBase struct {
	Result bool   `json:"success"`
	Error  string `json:"error,omitempty"`
}

func SendJSONResponse(c *fiber.Ctx, code int, response interface{}) error {
	c.Status(code)
	return c.JSON(response)
}
func SendJSONError(c *fiber.Ctx, err error) error {
	return SendJSONResponse(c, 500, ResponseBase{
		Result: false,
		Error:  err.Error(),
	})
}
