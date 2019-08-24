package twitter

import (
	"fmt"

	"github.com/ChimeraCoder/anaconda"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/xerrors"
)

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

type SourceApp struct {
	Id   string `mapstructure:"id"`
	Name string `mapstructure:"name"`
	URL  string `mapstructure:"url"`
}

type AppUser struct {
	Id                   string `mapstructure:"id"`
	CreatedTimestamp     string `mapstructure:"created_timestamp"`
	Name                 string `mapstructure:"name"`
	ScreenName           string `mapstructure:"screen_name"`
	Location             string `mapstructure:"location"`
	Description          string `mapstructure:"description"`
	URL                  string `mapstructure:"url"`
	Protected            bool   `mapstructure:"protected"`
	Verified             bool   `mapstructure:"verified"`
	FollowersCount       int    `mapstructure:"followers_count"`
	FriendsCount         int    `mapstructure:"friends_count"`
	StatusesCount        int    `mapstructure:"statuses_count"`
	ProfileImageURL      string `mapstructure:"profile_image_url"`
	ProfileImageURLHTTPS string `mapstructure:"profile_image_url_https"`
}

type Apps struct {
	SourceApp *SourceApp
	Sender    *AppUser
	Recipient *AppUser
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
	Apps                              map[string]interface{}
	Users                             map[string]interface{}
}

func (a *AccountActivityEvent) GetApps() (*Apps, error) {
	if a.DirectMessageEvents == nil || len(a.DirectMessageEvents) == 0 {
		return nil, xerrors.Errorf("failed to get apps: DirectMessageEvents is nil")
	}
	event := a.DirectMessageEvents[0] // FIXME
	sourceAppId := event.MessageCreate.SourceAppId
	sourceAppMap, ok := a.Apps[sourceAppId]
	if !ok {
		return nil, xerrors.Errorf("invalid apps: source app(id: %s) is empty", sourceAppId)
	}
	var sourceApp SourceApp
	if err := mapstructure.Decode(sourceAppMap, &sourceApp); err != nil {
		return nil, xerrors.Errorf("failed to decode source app: %#v", sourceAppMap)
	}

	senderId := event.MessageCreate.SenderId
	senderMap, ok := a.Apps[senderId]
	if !ok {
		return nil, xerrors.Errorf("invalid apps: sender (id: %s) is empty", senderId)
	}
	var sender AppUser
	if err := mapstructure.Decode(senderMap, &sender); err != nil {
		return nil, xerrors.Errorf("failed to decode sender: %#v", senderMap)
	}

	recipientId := event.MessageCreate.Target.RecipientId
	recipientMap, ok := a.Apps[recipientId]
	if !ok {
		return nil, xerrors.Errorf("invalid apps: sender (id: %s) is empty", senderId)
	}
	var recipient AppUser
	if err := mapstructure.Decode(recipientMap, &recipient); err != nil {
		return nil, xerrors.Errorf("failed to decode recipient: %#v", recipientMap)
	}

	return &Apps{
		SourceApp: &sourceApp,
		Sender:    &sender,
		Recipient: &recipient,
	}, nil
}

func (a *AccountActivityEvent) GetUsers() (users []*AppUser, err error) {
	for _, v := range a.Users {
		var user AppUser
		fmt.Printf("%#v\n", v)
		if err := mapstructure.Decode(v, &user); err != nil {
			return nil, xerrors.Errorf("failed to decode sender: %#v", v)
		}
		users = append(users, &user)
	}
	return
}
