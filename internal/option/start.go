package option

// StartCmdConfig is config for start command
type StartCmdConfig struct {
	TwitterConsumerKey       string `mapstructure:"TWITTER_CONSUMER_KEY"`
	TwitterConsumerSecret    string `mapstructure:"TWITTER_CONSUMER_SECRET"`
	TwitterAccessToken       string `mapstructure:"TWITTER_ACCESS_TOKEN"`
	TwitterAccessTokenSecret string `mapstructure:"TWITTER_ACCESS_TOKEN_SECRET"`
	ErrorMessage             string
	TweetKeyword             string
	InReplyToUserID          int64
	ClassifierServerHost     string
	Port                     string `mapstructure:"PORT"`
}

// NewStartCmdConfigFromViper generate config for start command from viper
func NewStartCmdConfigFromViper() (*StartCmdConfig, error) {
	rawConfig, err := newCmdRawConfig()
	return newStartCmdConfigFromRawConfig(rawConfig), err
}

func newStartCmdConfigFromRawConfig(rawConfig *CmdRawConfig) *StartCmdConfig {
	return &(rawConfig.StartCmdConfig)
}
