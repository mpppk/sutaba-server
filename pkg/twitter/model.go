package twitter

import "github.com/ChimeraCoder/anaconda"

type TweetCreateEvents struct {
	ForUserId         string           `json:"for_user_id"`
	TweetCreateEvents []anaconda.Tweet `json:"tweet_create_events"`
}

type CRCRequest struct {
	CRCToken string `json:"crc_token" query:"crc_token"`
}

type CRCResponse struct {
	ResponseToken string `json:"response_token"`
}
