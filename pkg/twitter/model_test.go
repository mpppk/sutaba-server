package twitter

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/ChimeraCoder/anaconda"
	"github.com/google/go-cmp/cmp"
)

var testDir = "../../testdata/aaa"

func TestAccountActivityEvents(t *testing.T) {
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
								RecipientId string `json:"recipient_id"`
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
				Apps: map[string]interface{}{
					"13090192": map[string]interface{}{
						"id":   "13090192",
						"name": "FuriousCamperTestApp1",
						"url":  "https://twitter.com/furiouscamper",
					},
					"3001969357": map[string]interface{}{
						"created_timestamp":       "1422556069340",
						"description":             "Alter Ego - Twitter PE opinions-are-my-own",
						"followers_count":         float64(22),
						"friends_count":           float64(45),
						"id":                      "3001969357",
						"location":                "Boulder, CO",
						"name":                    "Jordan Brinks",
						"profile_image_url":       "null",
						"profile_image_url_https": "https://pbs.twimg.com/profile_images/851526626785480705/cW4WTi7C_normal.jpg",
						"protected":               false,
						"screen_name":             "furiouscamper",
						"statuses_count":          float64(494),
						"url":                     "https://t.co/SnxaA15ZuY",
						"verified":                false,
					},
					"4337869213": map[string]interface{}{
						"created_timestamp":       "1448312972328",
						"followers_count":         float64(8),
						"friends_count":           float64(8),
						"id":                      "4337869213",
						"location":                "Burlington, MA",
						"name":                    "Harrison Test",
						"profile_image_url":       "null",
						"profile_image_url_https": "https://abs.twimg.com/sticky/default_profile_images/default_profile_normal.png",
						"protected":               false,
						"screen_name":             "Harris_0ff",
						"statuses_count":          float64(240),
						"verified":                false,
					},
					"users": map[string]interface{}{},
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
							RecipientId string `json:"recipient_id"`
						}{"3001969357"},
					},
				},
				Users: map[string]interface{}{
					"3001969357": map[string]interface{}{
						"created_timestamp":       "1422556069340",
						"description":             "Alter Ego - Twitter PE opinions-are-my-own",
						"followers_count":         float64(23),
						"friends_count":           float64(47),
						"id":                      "3001969357",
						"location":                "Boulder, CO",
						"name":                    "Jordan Brinks",
						"profile_image_url":       "null",
						"profile_image_url_https": "https://pbs.twimg.com/profile_images/851526626785480705/cW4WTi7C_normal.jpg",
						"protected":               false,
						"screen_name":             "furiouscamper",
						"statuses_count":          float64(509),
						"url":                     "https://t.co/SnxaA15ZuY",
						"verified":                false,
					},
					"3284025577": map[string]interface{}{
						"created_timestamp":       "1437281176085",
						"followers_count":         float64(1),
						"friends_count":           float64(4),
						"id":                      "3284025577",
						"name":                    "Bogus Bogart",
						"profile_image_url":       "null",
						"profile_image_url_https": "https://pbs.twimg.com/profile_images/763383202857779200/ndvZ96mE_normal.jpg",
						"protected":               true,
						"screen_name":             "bogusbogart",
						"statuses_count":          float64(35),
						"verified":                false,
					},
				},
			},
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
							RecipientId string `json:"recipient_id"`
						}{"3001969357"},
						LastReadEventId: "963085315333238788",
					},
				},
				Users: map[string]interface{}{
					"199566737": map[string]interface{}{
						"created_timestamp":       "1286429788000",
						"description":             "data by day @twitter, design by dusk",
						"followers_count":         float64(299),
						"friends_count":           float64(336),
						"id":                      "199566737",
						"location":                "Denver, CO",
						"name":                    "Le Braat",
						"profile_image_url":       "null",
						"profile_image_url_https": "https://pbs.twimg.com/profile_images/936652894371119105/YHEozVAg_normal.jpg",
						"protected":               false,
						"screen_name":             "LeBraat",
						"statuses_count":          float64(752),
						"verified":                false,
					},
					"3001969357": map[string]interface{}{
						"created_timestamp":       "1422556069340",
						"description":             "Alter Ego - Twitter PE opinions-are-my-own",
						"followers_count":         float64(23),
						"friends_count":           float64(48),
						"id":                      "3001969357",
						"location":                "Boulder, CO",
						"name":                    "Jordan Brinks",
						"profile_image_url":       "null",
						"profile_image_url_https": "https://pbs.twimg.com/profile_images/851526626785480705/cW4WTi7C_normal.jpg",
						"protected":               false,
						"screen_name":             "furiouscamper",
						"statuses_count":          float64(510),
						"url":                     "https://t.co/SnxaA15ZuY",
						"verified":                false,
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
							Id     string `json:"id"`
							UserId string `json:"user_id"`
						}{
							Id:     "1045405559317569537",
							UserId: "930524282358325248",
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
			t.Fatalf("failed to unmarshal tweet_create_events: %s", contents)
		}

		if diff := cmp.Diff(aae, *want); diff != "" {
			t.Errorf("AccountActivityEvent differs: (-got +want)\n%s", diff)
		}
	}
}

func TestAccountActivityEvent_GetApps(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		want     Apps
		wantErr  bool
	}{
		{
			name:     "direct_message_events",
			fileName: "direct_message_events.json",
			want: Apps{
				SourceApp: &SourceApp{
					Id:   "13090192",
					Name: "FuriousCamperTestApp1",
					URL:  "https://twitter.com/furiouscamper",
				},
				Sender: &AppUser{
					Id:                   "3001969357",
					CreatedTimestamp:     "1422556069340",
					Name:                 "Jordan Brinks",
					ScreenName:           "furiouscamper",
					Location:             "Boulder, CO",
					Description:          "Alter Ego - Twitter PE opinions-are-my-own",
					URL:                  "https://t.co/SnxaA15ZuY",
					Protected:            false,
					Verified:             false,
					FollowersCount:       22,
					FriendsCount:         45,
					StatusesCount:        494,
					ProfileImageURL:      "null",
					ProfileImageURLHTTPS: "https://pbs.twimg.com/profile_images/851526626785480705/cW4WTi7C_normal.jpg",
				},
				Recipient: &AppUser{
					Id:                   "4337869213",
					CreatedTimestamp:     "1448312972328",
					Name:                 "Harrison Test",
					ScreenName:           "Harris_0ff",
					Location:             "Burlington, MA",
					Protected:            false,
					Verified:             false,
					FollowersCount:       8,
					FriendsCount:         8,
					ProfileImageURL:      "null",
					StatusesCount:        240,
					ProfileImageURLHTTPS: "https://abs.twimg.com/sticky/default_profile_images/default_profile_normal.png",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := filepath.Join(testDir, tt.fileName)
			contents, err := ioutil.ReadFile(filePath)
			if err != nil {
				t.Fatalf("failed to read test file from %s", filePath)
			}

			var aae AccountActivityEvent
			if err := json.Unmarshal(contents, &aae); err != nil {
				t.Fatalf("failed to unmarshal tweet_create_events: %s", contents)
			}

			got, err := aae.GetApps()
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountActivityEvent.GetApps() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, &tt.want); diff != "" {
				t.Errorf("AccountActivityEvent.Apps differs: (-got +want)\n%s", diff)
			}
		})
	}
}

func TestAccountActivityEvent_GetUsers(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		want     []*AppUser
		wantErr  bool
	}{
		{
			name:     "direct_message_indicate_typing_events",
			fileName: "direct_message_indicate_typing_events.json",
			want: []*AppUser{
				{
					Id:                   "3001969357",
					CreatedTimestamp:     "1422556069340",
					Name:                 "Jordan Brinks",
					ScreenName:           "furiouscamper",
					Location:             "Boulder, CO",
					Description:          "Alter Ego - Twitter PE opinions-are-my-own",
					URL:                  "https://t.co/SnxaA15ZuY",
					Protected:            false,
					Verified:             false,
					FollowersCount:       23,
					FriendsCount:         47,
					StatusesCount:        509,
					ProfileImageURL:      "null",
					ProfileImageURLHTTPS: "https://pbs.twimg.com/profile_images/851526626785480705/cW4WTi7C_normal.jpg",
				},
				{
					Id:                   "3284025577",
					CreatedTimestamp:     "1437281176085",
					Name:                 "Bogus Bogart",
					ScreenName:           "bogusbogart",
					Protected:            true,
					Verified:             false,
					FollowersCount:       1,
					FriendsCount:         4,
					StatusesCount:        35,
					ProfileImageURL:      "null",
					ProfileImageURLHTTPS: "https://pbs.twimg.com/profile_images/763383202857779200/ndvZ96mE_normal.jpg",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := filepath.Join(testDir, tt.fileName)
			contents, err := ioutil.ReadFile(filePath)
			if err != nil {
				t.Fatalf("failed to read test file from %s", filePath)
			}

			var aae AccountActivityEvent
			if err := json.Unmarshal(contents, &aae); err != nil {
				t.Fatalf("failed to unmarshal tweet_create_events: %s", contents)
			}

			gotUsers, err := aae.GetUsers()
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountActivityEvent.GetUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			compare := func(gotUser *AppUser, wantUsers []*AppUser) {
				for _, wantUser := range wantUsers {
					if gotUser.Id == wantUser.Id {
						if diff := cmp.Diff(gotUser, wantUser); diff != "" {
							t.Errorf("AccountActivityEvent.Users differs: (-got +want)\n%s", diff)
						}
						return
					}
				}
				t.Errorf("gotUser not found: expected: %#v", gotUser)
			}

			for _, gotUser := range gotUsers {
				compare(gotUser, tt.want)
			}

		})
	}
}
