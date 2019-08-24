package twitter

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/ChimeraCoder/anaconda"
)

func TestAccountActivityEvents(t *testing.T) {
	testDir := "../../testdata/aaa"
	tests := []struct {
		name     string
		fileName string
		want     AccountActivityEvent
	}{
		{
			name:     "tweet_create_events",
			fileName: "tweet_create_events.json",
			want: AccountActivityEvent{
				ForUserId: "2244994945",
				TweetCreateEvents: []*anaconda.Tweet{
					{
						Id: 1,
					},
				},
			},
		},
		{
			name:     "favorite_events",
			fileName: "favorite_events.json",
			want: AccountActivityEvent{
				ForUserId: "2244994945",
				FavoriteEvents: []*FavoriteEvent{
					{
						Id:              "a7ba59eab0bfcba386f7acedac279542",
						CreatedAt:       "Mon Mar 26 16:33:26 +0000 2018",
						TimestampMs:     1522082006140,
						FavoritedStatus: anaconda.Tweet{Id: 1},
						User:            anaconda.User{Id: 2},
					},
				},
			},
		},
		{
			name:     "follow_events",
			fileName: "follow_events.json",
			want: AccountActivityEvent{
				ForUserId: "2244994945",
				FollowEvents: []*UserRelationEvent{
					{
						Type:             EventType(FollowEventType),
						CreatedTimestamp: "1517588749178",
						Target:           &anaconda.User{Id: 1},
						Source:           &anaconda.User{Id: 2},
					},
				},
			},
		},
		{
			name:     "unfollow_events",
			fileName: "unfollow_events.json",
			want: AccountActivityEvent{
				ForUserId: "2244994945",
				FollowEvents: []*UserRelationEvent{
					{
						Type:             EventType(UnfollowEventType),
						CreatedTimestamp: "1517588749178",
						Target:           &anaconda.User{Id: 1},
						Source:           &anaconda.User{Id: 2},
					},
				},
			},
		},
		{
			name:     "block_events",
			fileName: "block_events.json",
			want: AccountActivityEvent{
				ForUserId: "2244994945",
				BlockEvents: []*UserRelationEvent{
					{
						Type:             EventType(BlockEventType),
						CreatedTimestamp: "1518127020304",
						Target:           &anaconda.User{Id: 1},
						Source:           &anaconda.User{Id: 2},
					},
				},
			},
		},
		{
			name:     "unblock_events",
			fileName: "unblock_events.json",
			want: AccountActivityEvent{
				ForUserId: "2244994945",
				BlockEvents: []*UserRelationEvent{
					{
						Type:             EventType(UnblockEventType),
						CreatedTimestamp: "1518127020304",
						Target:           &anaconda.User{Id: 1},
						Source:           &anaconda.User{Id: 2},
					},
				},
			},
		},
		{
			name:     "mute_events",
			fileName: "mute_events.json",
			want: AccountActivityEvent{
				ForUserId: "2244994945",
				MuteEvents: []*UserRelationEvent{
					{
						Type:             EventType(MuteEventType),
						CreatedTimestamp: "1518127020304",
						Target:           &anaconda.User{Id: 1},
						Source:           &anaconda.User{Id: 2},
					},
				},
			},
		},
		{
			name:     "unmute_events",
			fileName: "unmute_events.json",
			want: AccountActivityEvent{
				ForUserId: "2244994945",
				MuteEvents: []*UserRelationEvent{
					{
						Type:             EventType(UnmuteEventType),
						CreatedTimestamp: "1518127020304",
						Target:           &anaconda.User{Id: 1},
						Source:           &anaconda.User{Id: 2},
					},
				},
			},
		},
		{
			name:     "user_event",
			fileName: "user_event.json",
			want: AccountActivityEvent{
				UserEvent: &UserEvent{
					Revoke: &Revoke{
						DateTime: "2018-05-24T09:48:12+00:00",
						Target: struct {
							AppId string `json:"app_id"`
						}{"13090192"},
						Source: struct {
							UserId string `json:"user_id"`
						}{"63046977"},
					},
				},
			},
		},
		{
			name:     "direct_message_events",
			fileName: "direct_message_events.json",
			want: AccountActivityEvent{
				ForUserId: "4337869213",
				DirectMessageEvents: []*DirectMessageEvent{
					{
						Type:             EventType(MessageCreateEventType),
						Id:               "954491830116155396",
						CreatedTimestamp: "1516403560557",
						MessageCreate: &MessageCreate{
							Target: struct {
								RecipientID string `json:"recipient_id"`
							}{"4337869213"},
							SenderId:    "3001969357",
							SourceAppId: "13090192",
							MessageData: MessageData{
								Text: "Hello World!",
								Entities: &Entities{
									Hashtags:     []interface{}{},
									Symbols:      []interface{}{},
									UserMentions: []interface{}{},
									Urls:         []interface{}{},
								},
							},
						},
					},
				},
			},
		},
		{
			name:     "direct_message_indicate_typing_events",
			fileName: "direct_message_indicate_typing_events.json",
			want: AccountActivityEvent{
				ForUserId: "4337869213",
				DirectMessageIndicateTypingEvents: []*DirectMessageIndicateTypingEvent{
					{
						CreatedTimestamp: "1518127183443",
						SenderId:         "3284025577",
						Target: struct {
							RecipientID string `json:"recipient_id"`
						}{"3001969357"},
					},
				},
			},
			// TODO
		},
		{
			name:     "direct_message_mark_read_events",
			fileName: "direct_message_mark_read_events.json",
			want: AccountActivityEvent{
				ForUserId: "4337869213",
				DirectMessageMarkReadEvents: []*DirectMessageMarkReadEvent{
					{
						CreatedTimestamp: "1518452444662",
						SenderId:         "199566737",
						Target: struct {
							RecipientID string `json:"recipient_id"`
						}{"3001969357"},
						LastReadEventId: "963085315333238788",
					},
				},
			},
		},
		{
			name:     "tweet_delete_events",
			fileName: "tweet_delete_events.json",
			want: AccountActivityEvent{
				ForUserId: "930524282358325248",
				TweetDeleteEvents: []*TweetDeleteEvent{
					{
						Status: struct {
							ID     string `json:"id"`
							UserID string `json:"user_id"`
						}{
							ID:     "1045405559317569537",
							UserID: "930524282358325248",
						},
						TimestampMs: "1432228155593",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		filePath := filepath.Join(testDir, tt.fileName)
		t.Run(tt.name, generateTweetEventTestFunc(filePath, &tt.want))
	}
}

func generateTweetEventTestFunc(filePath string, want *AccountActivityEvent) func(t *testing.T) {
	return func(t *testing.T) {
		contents, err := ioutil.ReadFile(filePath)
		if err != nil {
			t.Fatalf("failed to read test file from %s", filePath)
		}

		var aae AccountActivityEvent
		if err := json.Unmarshal(contents, &aae); err != nil {
			t.Fatalf("failed to unmarshal tweet_create_events: %s", string(contents))
		}

		if diff := cmp.Diff(aae, *want); diff != "" {
			t.Errorf("AccountActivityEvent differs: (-got +want)\n%s", diff)
		}
	}
}
