package cmd

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4/middleware"

	"github.com/ChimeraCoder/anaconda"

	"github.com/mpppk/sutaba-server/internal/option"

	"github.com/labstack/echo/v4"
	"github.com/spf13/afero"

	"github.com/spf13/cobra"
)

type CRCRequest struct {
	CRCToken string `json:"crc_token" query:"crc_token"`
}

type CRCResponse struct {
	ResponseToken string `json:"response_token"`
}

type TweetCreateEvents struct {
	ForUserId         string           `json:"for_user_id"`
	TweetCreateEvents []anaconda.Tweet `json:"tweet_create_events"`
}

func bodyDumpHandler(c echo.Context, reqBody, resBody []byte) {
	fmt.Printf("Request Body: %v\n", string(reqBody))
	fmt.Printf("Response Body: %v\n", string(resBody))
}

func newServerCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start server",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := option.NewServerCmdConfigFromViper()
			if err != nil {
				return err
			}

			e := echo.New()
			e.Use(middleware.BodyDump(bodyDumpHandler))

			endpoint := "/twitter/aaa"
			e.GET(endpoint, func(c echo.Context) error {
				req := new(CRCRequest)
				if err = c.Bind(req); err != nil {
					return err
				}
				response := &CRCResponse{ResponseToken: CreateCRCToken(req.CRCToken, conf.TwitterConsumerSecret)}
				return c.JSON(http.StatusOK, response)
			})

			e.POST(endpoint, func(c echo.Context) error {
				events := new(TweetCreateEvents)
				if err = c.Bind(events); err != nil {
					return err
				}
				fmt.Printf("%#v\n", events)
				if events.TweetCreateEvents == nil {
					return c.NoContent(http.StatusNoContent)
				}

				return c.NoContent(http.StatusNoContent)
			})

			port := "1323"
			envPort := os.Getenv("PORT")
			if envPort != "" {
				port = envPort
			}
			e.Logger.Fatal(e.Start(":" + port))
			return nil
		},
	}
	return cmd, nil
}

func CreateCRCToken(crcToken, consumerSecret string) string {
	mac := hmac.New(sha256.New, []byte(consumerSecret))
	mac.Write([]byte(crcToken))
	return "sha256=" + base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func init() {
	cmdGenerators = append(cmdGenerators, newServerCmd)
}
