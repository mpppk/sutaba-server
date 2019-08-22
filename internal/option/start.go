package option

// StartCmdConfig is config for start command
type StartCmdConfig struct {
	TwitterConsumerKey            string `mapstructure:"TWITTER_CONSUMER_KEY"`
	TwitterConsumerSecret         string `mapstructure:"TWITTER_CONSUMER_SECRET"`
	OwnerTwitterAccessToken       string `mapstructure:"OWNER_TWITTER_ACCESS_TOKEN"`
	OwnerTwitterAccessTokenSecret string `mapstructure:"OWNER_TWITTER_ACCESS_TOKEN_SECRET"`
	BotTwitterAccessToken         string `mapstructure:"BOT_TWITTER_ACCESS_TOKEN"`
	BotTwitterAccessTokenSecret   string `mapstructure:"BOT_TWITTER_ACCESS_TOKEN_SECRET"`
	ErrorTweetMessage             string `mapstructure:"ERROR_TWEET_MESSAGE"`
	SorryTweetMessage             string `mapstructure:"SORRY_TWEET_MESSAGE"`
	TweetKeyword                  string `mapstructure:"TWEET_KEYWORD"`
	OwnerTwitterUserID            int64  `mapstructure:"OWNER_TWITTER_USER_ID"`
	BotTwitterUserID              int64  `mapstructure:"BOT_TWITTER_USER_ID"`
	ClassifierServerHost          string `mapstructure:"CLASSIFIER_SERVER_HOST"`
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
