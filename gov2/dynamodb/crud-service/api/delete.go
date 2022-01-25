package api

import (
	"errors"

	"example.aws/gov2/dynamodb/crud-service/db"
	"github.com/gofiber/fiber/v2"
)

type DeleteRequest struct {
	DeleteKey string `json:"deleteKey"`
}

type DeleteResponse struct {
	Id string `json:"id,omitempty"`
	ResponseBase
}

func DeleteLink(c *fiber.Ctx) error {

	id := c.Params("id", "")
	if id == "" {
		return SendJSONError(c, errors.New("no id specified"))
	}

	request := DeleteRequest{}

	if err := c.BodyParser(request); err != nil {
		return SendJSONError(c, err)
	}
	if request.DeleteKey == "" {
		return SendJSONError(c, errors.New("no delete key specified"))
	}

	link := (*db.DB).Get(id)

	if link == nil {
		return SendJSONResponse(c, 404, DeleteResponse{
			Id: id,
			ResponseBase: ResponseBase{
				Result: false,
				Error:  "Unknown link ID",
			},
		})
	}

	if link.DeleteKey != request.DeleteKey {
		return SendJSONResponse(c, 400, DeleteResponse{
			Id: id,
			ResponseBase: ResponseBase{Result: false,
				Error: "Bad delete key"},
		})
	}

	if (*db.DB).Delete(id) {
		return SendJSONResponse(c, 200, DeleteResponse{
			Id: id,
			ResponseBase: ResponseBase{
				Result: true,
				Error:  "no error"},
		})
	} else {
		return SendJSONResponse(c, 500, ResponseBase{
			Result: false,
			Error:  "internal error",
		})
	}
}
