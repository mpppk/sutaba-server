package option

// StartCmdConfig is config for start command
type StartCmdConfig struct {
	TwitterConsumerKey            string `mapstructure:"TWITTER_CONSUMER_KEY"`
	TwitterConsumerSecret         string `mapstructure:"TWITTER_CONSUMER_SECRET"`
	OwnerTwitterAccessToken       string `mapstructure:"OWNER_TWITTER_ACCESS_TOKEN"`
	OwnerTwitterAccessTokenSecret string `mapstructure:"OWNER_TWITTER_ACCESS_TOKEN_SECRET"`
	BotTwitterAccessToken         string `mapstructure:"BOT_TWITTER_ACCESS_TOKEN"`
	BotTwitterAccessTokenSecret   string `mapstructure:"BOT_TWITTER_ACCESS_TOKEN_SECRET"`
	ErrorTweetMessage             string
	SorryTweetMessage             string
	TweetKeyword                  string
	OwnerTwitterUserID            int64
	BotTwitterUserID              int64
	ClassifierServerHost          string
	Port                          string `mapstructure:"PORT"`
}

// NewStartCmdConfigFromViper generate config for start command from viper
func NewStartCmdConfigFromViper() (*StartCmdConfig, error) {
	rawConfig, err := newCmdRawConfig()
	return newStartCmdConfigFromRawConfig(rawConfig), err
}

func newStartCmdConfigFromRawConfig(rawConfig *CmdRawConfig) *StartCmdConfig {
	return &(rawConfig.StartCmdConfig)
}
