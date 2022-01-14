package api

import "github.com/gofiber/fiber/v2"

func GetApi(app *fiber.App) fiber.Group {
	group := app.Group("/api", nil)

	group.Put("/link", CreateLink)
	group.Delete("/link", DeleteLink)
	group.Get("/link", GetLinkStats)
	group.Get("/recent", GetRecentLinks)
}
