package api

import (
	"example.aws/gov2/dynamodb/crud-service/db"
	"github.com/gofiber/fiber/v2"
)

type LinkStatsResponse struct {
	ResponseBase
	Id   string `json:"id"`
	Hits uint64 `json:"hits"`
}

func GetLinkStats(c *fiber.Ctx) error {

	id := c.Params("id", "")

	if id == "" {
		return SendJSONResponse(c, 400, ResponseBase{
			Error:  "no id specified",
			Result: false,
		})
	}

	// Get the link

	link := (*db.DB).Get(id)

	if link == nil {
		return SendJSONResponse(c, 404, ResponseBase{
			Error:  "not found",
			Result: false,
		})
	}

	return SendJSONResponse(c, 200, LinkStatsResponse{
		ResponseBase: ResponseBase{
			Error:  "no error",
			Result: true,
		},
		Id:   id,
		Hits: link.NumHits,
	})

}
