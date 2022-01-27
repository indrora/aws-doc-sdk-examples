package api

import "github.com/gofiber/fiber/v2"

func GetApi(app *fiber.App) fiber.Router {
	group := app.Group("/api", nil)

	group.Put("/link", CreateLink)
	group.Delete("/link/:id", DeleteLink)
	group.Get("/link/:id", GetLinkStats)
	group.Put("/link/report", GetLinks)

	return group
}

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
