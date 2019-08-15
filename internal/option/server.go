package option

// ServerCmdConfig is config for server command
type ServerCmdConfig struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

// NewServerCmdConfigFromViper generate config for server command from viper
func NewServerCmdConfigFromViper() (*ServerCmdConfig, error) {
	rawConfig, err := newCmdRawConfig()
	return newServerCmdConfigFromRawConfig(rawConfig), err
}

func newServerCmdConfigFromRawConfig(rawConfig *CmdRawConfig) *ServerCmdConfig {
	return &(rawConfig.ServerCmdConfig)
}
