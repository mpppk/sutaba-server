package twitter

import "github.com/ChimeraCoder/anaconda"

type EventType string

const (
	FollowEventType        EventType = "follow"
	UnfollowEventType      EventType = "unfollow"
	BlockEventType         EventType = "block"
	UnblockEventType       EventType = "unblock"
	MuteEventType          EventType = "mute"
	UnmuteEventType        EventType = "unmute"
	MessageCreateEventType EventType = "message_create"
)

type CRCRequest struct {
	CRCToken string `json:"crc_token" query:"crc_token"`
}

type CRCResponse struct {
	ResponseToken string `json:"response_token"`
}

type FavoriteEvent struct {
	Id              string         `json:"id"`
	CreatedAt       string         `json:"created_at"`
	TimestampMs     int            `json:"timestamp_ms"`
	FavoritedStatus anaconda.Tweet `json:"favorited_status"`
	User            anaconda.User  `json:"user"`
}

type UserRelationEvent struct {
	Type             EventType      `json:"type"`
	CreatedTimestamp string         `json:"created_timestamp"`
	Target           *anaconda.User `json:"target"`
	Source           *anaconda.User `json:"source"`
}

type Revoke struct {
	DateTime string `json:"date_time"`
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

type Entities struct {
	Hashtags     []interface{} `json:"hashtags"`
	Symbols      []interface{} `json:"symbols"`
	UserMentions []interface{} `json:"user_mentions"`
	Urls         []interface{} `json:"urls"`
}

type MessageData struct {
	Text     string    `json:"text"`
	Entities *Entities `json:"entities"`
}

type MessageCreate struct {
	Target struct {
		RecipientId string `json:"recipient_id"`
	} `json:"target"`
	SenderId    string      `json:"sender_id"`
	SourceAppId string      `json:"source_app_id"`
	MessageData MessageData `json:"message_data"`
}

type DirectMessageEvent struct {
	Type             EventType      `json:"type"`
	Id               string         `json:"id"`
	CreatedTimestamp string         `json:"created_timestamp"`
	MessageCreate    *MessageCreate `json:"message_create"`
}

type DirectMessageIndicateTypingEvent struct {
	CreatedTimestamp string `json:"created_timestamp"`
	SenderId         string `json:"sender_id"`
	Target           struct {
		RecipientId string `json:"recipient_id"`
	} `json:"target"`
}

type DirectMessageMarkReadEvent struct {
	CreatedTimestamp string `json:"created_timestamp"`
	SenderId         string `json:"sender_id"`
	Target           struct {
		RecipientId string `json:"recipient_id"`
	} `json:"target"`
	LastReadEventId string `json:"last_read_event_id"`
}

type TweetDeleteEvent struct {
	Status struct {
		Id     string `json:"id"`
		UserId string `json:"user_id"`
	} `json:"status"`
	TimestampMs string `json:"timestamp_ms"`
}

type AccountActivityEvent struct {
	ForUserId                         string                              `json:"for_user_id"`
	TweetCreateEvents                 []*anaconda.Tweet                   `json:"tweet_create_events"`
	FavoriteEvents                    []*FavoriteEvent                    `json:"favorite_events"`
	FollowEvents                      []*UserRelationEvent                `json:"follow_events"`
	BlockEvents                       []*UserRelationEvent                `json:"block_events"`
	MuteEvents                        []*UserRelationEvent                `json:"mute_events"`
	UserEvent                         *UserEvent                          `json:"user_event"`
	DirectMessageEvents               []*DirectMessageEvent               `json:"direct_message_events"`
	DirectMessageIndicateTypingEvents []*DirectMessageIndicateTypingEvent `json:"direct_message_indicate_typing_events"`
	DirectMessageMarkReadEvents       []*DirectMessageMarkReadEvent       `json:"direct_message_mark_read_events"`
	TweetDeleteEvents                 []*TweetDeleteEvent                 `json:"tweet_delete_events"`
}
