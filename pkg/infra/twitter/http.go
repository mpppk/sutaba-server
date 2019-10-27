package twitter

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mpppk/anacondaaaa"
)

func GenerateCRCTestHandler(twitterConsumerSecret string) func(c echo.Context) error {
	return func(c echo.Context) error {
		req := new(anacondaaaa.CRCRequest)
		if err := c.Bind(req); err != nil {
			return err
		}
		response := &anacondaaaa.CRCResponse{ResponseToken: CreateCRCToken(req.CRCToken, twitterConsumerSecret)}
		return c.JSON(http.StatusOK, response)
	}
}
