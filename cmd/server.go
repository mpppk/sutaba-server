package cmd

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"os"

	"github.com/mpppk/cli-template/internal/option"

	"github.com/labstack/echo/v4"
	"github.com/spf13/afero"

	"github.com/spf13/cobra"
)

type CRCRequest struct {
	CRCToken string `json:"crc_token"`
}

type CRCResponse struct {
	ResponseToken string `json:"response_token"`
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
			e.GET("/twitter/aaa", func(c echo.Context) error {
				req := new(CRCRequest)
				response := &CRCResponse{ResponseToken: CreateCRCToken(req.CRCToken, conf.ConsumerSecret)}
				return c.JSON(http.StatusOK, response)
			})
			port := "5000"
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
