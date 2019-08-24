package twitter

import "github.com/ChimeraCoder/anaconda"

type CRCRequest struct {
	CRCToken string `json:"crc_token" query:"crc_token"`
}

type CRCResponse struct {
	ResponseToken string `json:"response_token"`
}

type FavoriteEvent struct {
	Id              string         `json:"id"`
	CreatedAt       string         `json:"created_at"`
	TimestampMs     string         `json:"timestamp_ms"`
	FavoritedStatus anaconda.Tweet `json:"favorited_status"`
	User            anaconda.User  `json:"user"`
}

type UserRelationEvent struct {
	Type             string         `json:"type"`
	CreatedTimestamp string         `json:"created_timestamp"`
	Target           *anaconda.User `json:"target"`
	Source           *anaconda.User `json:"source"`
}

type Revoke struct {
	DateTime string
	Target   struct {
		AppId string `json:"app_id"`
	}
	Source struct {
		UserId string `json:"user_id"`
	}
}

type UserEvent struct {
	Revoke *Revoke `json:"revoke"`
}

type DirectMessageEvent struct {
	Type             string `json:"type"`
	ID               string `json:"id"`
	CreatedTimestamp string `json:"created_timestamp"`
	MessageCreate    struct {
		Target struct {
			RecipientID string `json:"recipient_id"`
		} `json:"target"`
		SenderID    string `json:"sender_id"`
		SourceAppID string `json:"source_app_id"`
		MessageData struct {
			Text     string `json:"text"`
			Entities struct {
				Hashtags     []interface{} `json:"hashtags"`
				Symbols      []interface{} `json:"symbols"`
				UserMentions []interface{} `json:"user_mentions"`
				Urls         []interface{} `json:"urls"`
			} `json:"entities"`
		} `json:"message_data"`
	} `json:"message_create"`
}

type DirectMessageIndicateTypingEvent struct {
	CreatedTimestamp string `json:"created_timestamp"`
	SenderID         string `json:"sender_id"`
	Target           struct {
		RecipientID string `json:"recipient_id"`
	} `json:"target"`
}

type DirectMessageMarkReadEvent struct {
	CreatedTimestamp string `json:"created_timestamp"`
	SenderID         string `json:"sender_id"`
	Target           struct {
		RecipientID string `json:"recipient_id"`
	} `json:"target"`
	LastReadEventID string `json:"last_read_event_id"`
}

type TweetDeleteEvent struct {
	Status struct {
		ID     string `json:"id"`
		UserID string `json:"user_id"`
	} `json:"status"`
	TimestampMs string `json:"timestamp_ms"`
}

type AccountActivityEvent struct {
	ForUserId                         string                              `json:"for_user_id"`
	TweetCreateEvents                 []anaconda.Tweet                    `json:"tweet_create_events"`
	FavoriteEvents                    []*FavoriteEvent                    `json:"favorite_events"`
	FollowEvents                      []*UserRelationEvent                `json:"follow_events"`
	BlockEvents                       []*UserRelationEvent                `json:"block_events"`
	MuteEvents                        []*UserRelationEvent                `json:"mute_events"`
	UserEvent                         UserEvent                           `json:"user_event"`
	DirectMessageEvents               []*DirectMessageEvent               `json:"direct_message_events"`
	DirectMessageIndicateTypingEvents []*DirectMessageIndicateTypingEvent `json:"direct_message_indicate_typing_events"`
	DirectMessageMarkReadEvents       []*DirectMessageMarkReadEvent       `json:"direct_message_mark_read_events"`
	TweetDeleteEvents                 []*TweetDeleteEvent                 `json:"tweet_delete_events"`
}
