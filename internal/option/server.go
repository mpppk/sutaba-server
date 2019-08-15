package option

// ServerCmdConfig is config for server command
type ServerCmdConfig struct {
	TwitterConsumerKey       string `mapstructure:"TWITTER_CONSUMER_KEY"`
	TwitterConsumerSecret    string `mapstructure:"TWITTER_CONSUMER_SECRET"`
	TwitterAccessToken       string `mapstructure:"TWITTER_ACCESS_TOKEN"`
	TwitterAccessTokenSecret string `mapstructure:"TWITTER_ACCESS_TOKEN_SECRET"`
}

// NewServerCmdConfigFromViper generate config for server command from viper
func NewServerCmdConfigFromViper() (*ServerCmdConfig, error) {
	rawConfig, err := newCmdRawConfig()
	return newServerCmdConfigFromRawConfig(rawConfig), err
}

func newServerCmdConfigFromRawConfig(rawConfig *CmdRawConfig) *ServerCmdConfig {
	return &(rawConfig.ServerCmdConfig)
}
