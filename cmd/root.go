package cmd

import (
	"fmt"
	"os"

	"github.com/mpppk/sutaba-server/internal/option"
	"github.com/spf13/afero"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

func newToggleFlag() *option.BoolFlag {
	return &option.BoolFlag{
		Flag: &option.Flag{
			Name:  "toggle",
			Usage: "Do nothing",
		},
		Value: false,
	}
}

func NewRootCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "sutaba-server",
		Short: "sutaba-server",
	}

	configFlag := &option.StringFlag{
		Flag: &option.Flag{
			Name:         "config",
			IsPersistent: true,
			Usage:        "config file (default is $HOME/.sutaba-server.yaml)",
		},
	}

	if err := option.RegisterStringFlag(cmd, configFlag); err != nil {
		return nil, err
	}

	var subCmds []*cobra.Command
	for _, cmdGen := range cmdGenerators {
		subCmd, err := cmdGen(fs)
		if err != nil {
			return nil, err
		}
		subCmds = append(subCmds, subCmd)
	}
	cmd.AddCommand(subCmds...)

	return cmd, nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd, err := NewRootCmd(afero.NewOsFs())
	if err != nil {
		panic(err)
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".sutaba-server" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".sutaba-server")
	}

	viper.AutomaticEnv() // read in environment variables that match
	if err := viper.BindEnv("TWITTER_CONSUMER_KEY"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := viper.BindEnv("TWITTER_CONSUMER_SECRET"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := viper.BindEnv("TWITTER_ACCESS_TOKEN"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := viper.BindEnv("TWITTER_ACCESS_TOKEN_SECRET"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
