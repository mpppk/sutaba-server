package itwitter

import (
	"reflect"
	"testing"

	"github.com/mpppk/sutaba-server/pkg/domain/model"
)

func TestTwitter_NewMessage(t *testing.T) {
	type args struct {
		tweet *Tweet
	}
	tests := []struct {
		name string
		args args
		want *model.Message
	}{
		{
			name: "OK",
			args: args{
				tweet: &Tweet{
					ID:              1,
					User:            model.User{ID: 2},
					Text:            "test text",
					QuoteTweet:      &Tweet{ID: 3},
					InReplyToUserID: 4,
					MediaURLs:       []string{"url1", "url2"},
				},
			},
			want: &model.Message{
				ID:                1,
				User:              model.User{ID: 2},
				Text:              "test text",
				ReferencedMessage: &model.Message{ID: 3},
				MediaNum:          2,
				ReplyUserID:       4,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewTwitter()
			if got := r.NewMessage(tt.args.tweet); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTwitter_NewMessageEvent(t *testing.T) {
	type args struct {
		forUserID model.UserID
		tweet     *Tweet
	}
	tests := []struct {
		name string
		args args
		want *model.MessageEvent
	}{
		{
			name: "not shared event",
			args: args{
				forUserID: 1,
				tweet: &Tweet{
					ID:   1,
					User: model.User{ID: 2},
					Text: "test text",
				},
			},
			want: &model.MessageEvent{
				TargetUserID: 1,
				IsShared:     false,
				Message: &model.Message{
					ID:   1,
					User: model.User{ID: 2},
					Text: "test text",
				},
			},
		},
		{
			name: "shared event",
			args: args{
				forUserID: 1,
				tweet: &Tweet{
					ID:              1,
					User:            model.User{ID: 2},
					Text:            "test text",
					RetweetedStatus: &Tweet{},
				},
			},
			want: &model.MessageEvent{
				TargetUserID: 1,
				IsShared:     true,
				Message: &model.Message{
					ID:   1,
					User: model.User{ID: 2},
					Text: "test text",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewTwitter()
			if got := r.NewMessageEvent(tt.args.forUserID, tt.args.tweet); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMessageEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}
