package twitter

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GenerateCRCTestHandler(twitterConsumerSecret string) func(c echo.Context) error {
	return func(c echo.Context) error {
		req := new(CRCRequest)
		if err := c.Bind(req); err != nil {
			return err
		}
		response := &CRCResponse{ResponseToken: CreateCRCToken(req.CRCToken, twitterConsumerSecret)}
		return c.JSON(http.StatusOK, response)
	}
}
