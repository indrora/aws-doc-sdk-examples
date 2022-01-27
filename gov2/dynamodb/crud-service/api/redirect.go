package api

import (
	"example.aws/gov2/dynamodb/crud-service/db"
	"github.com/gofiber/fiber/v2"
)

func DoRedirect(c *fiber.Ctx) error {
	id := c.Params("id", "")

	if id == "" {
		return c.Redirect("/")
	} else {
		// Find it
		link := db.DB.Get(id)
		if link != nil {
			db.DB.Increment(id)
			return c.Redirect(link.Uri)
		} else {
			return c.Next()
		}
	}
}
