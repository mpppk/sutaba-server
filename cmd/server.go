package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/mpppk/sutaba-server/pkg/sutaba"

	"github.com/mpppk/sutaba-server/pkg/twitter"

	"github.com/labstack/echo/v4/middleware"

	"github.com/mpppk/sutaba-server/internal/option"

	"github.com/labstack/echo/v4"
	"github.com/spf13/afero"

	"github.com/spf13/cobra"
)

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
			e.GET(endpoint, twitter.GenerateCRCTestHandler(conf.TwitterConsumerSecret))

			e.POST(endpoint, sutaba.GeneratePredictHandler(conf))

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

func init() {
	location := "Asia/Tokyo"
	loc, err := time.LoadLocation(location)
	if err != nil {
		loc = time.FixedZone(location, 9*60*60)
	}
	time.Local = loc

	cmdGenerators = append(cmdGenerators, newServerCmd)
}
