package domain

import (
	"testing"

	"github.com/mpppk/sutaba-server/pkg/domain/model"
)

func TestIsTargetMessage(t *testing.T) {
	type args struct {
		botUser *model.User
		message *model.Message
	}
	tests := []struct {
		name           string
		args           args
		msgIsTarget    bool
		refMsgIsTarget bool
	}{
		{
			name: "自分自身でない、Mediaを含むメッセージは対象",
			args: args{
				botUser: &model.User{
					ID: 1,
				},
				message: &model.Message{
					User: model.User{
						ID: 2,
					},
					ReplyUserID: 1,
					MediaNum:    1,
				},
			},
			msgIsTarget:    true,
			refMsgIsTarget: false,
		},
		{
			name: "自分自身でない、Mediaを含むメッセージをreferしたメッセージは対象",
			args: args{
				botUser: &model.User{
					ID: 1,
				},
				message: &model.Message{
					User: model.User{
						ID: 2,
					},
					ReplyUserID: 1,
					MediaNum:    0,
					ReferencedMessage: &model.Message{
						User: model.User{
							ID: 3,
						},
						MediaNum: 1,
					},
				},
			},
			msgIsTarget:    false,
			refMsgIsTarget: true,
		},
		{
			name: "ReplyUserIDが異なっていても、Textにusernameが含まれていれば対象",
			args: args{
				botUser: &model.User{
					ID:   1,
					Name: "username",
				},
				message: &model.Message{
					User: model.User{
						ID: 2,
					},
					MediaNum:    1,
					ReplyUserID: 2,
					Text:        "@username text",
				},
			},
			msgIsTarget:    true,
			refMsgIsTarget: false,
		},
		{
			name: "メディアが付与されたメッセージでなければ対象外",
			args: args{
				botUser: &model.User{},
				message: &model.Message{
					MediaNum: 0,
				},
			},
			msgIsTarget: false,
		},
		{
			name: "自分自身のメッセージは対象外",
			args: args{
				botUser: &model.User{
					ID: 1,
				},
				message: &model.Message{
					User: model.User{
						ID: 1,
					},
					ReplyUserID: 1,
					MediaNum:    1,
				},
			},
			msgIsTarget:    false,
			refMsgIsTarget: false,
		},
		{
			name: "自分自身のメッセージをreferしているメッセージは対象外",
			args: args{
				botUser: &model.User{
					ID: 1,
				},
				message: &model.Message{
					User: model.User{
						ID: 2,
					},
					MediaNum:    1,
					ReplyUserID: 1,
					ReferencedMessage: &model.Message{
						User: model.User{
							ID: 1,
						},
						MediaNum: 1,
					},
				},
			},
			msgIsTarget:    false,
			refMsgIsTarget: false,
		},
		{
			name: "referしているメッセージが対象でも、元が自分自身のメッセージであれば対象外",
			args: args{
				botUser: &model.User{
					ID: 1,
				},
				message: &model.Message{
					User: model.User{
						ID: 1,
					},
					MediaNum:    1,
					ReplyUserID: 1,
					ReferencedMessage: &model.Message{
						User: model.User{
							ID: 2,
						},
						MediaNum: 1,
					},
				},
			},
			msgIsTarget:    false,
			refMsgIsTarget: false,
		},
		{
			name: "対象ユーザへのリプライでなければ対象外",
			args: args{
				botUser: &model.User{
					ID: 1,
				},
				message: &model.Message{
					User: model.User{
						ID: 2,
					},
					MediaNum:    1,
					ReplyUserID: 2,
				},
			},
			msgIsTarget:    false,
			refMsgIsTarget: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msgIsTarget, refMsgIsTarget, reason := IsTargetMessage(tt.args.botUser, tt.args.message)
			if msgIsTarget != tt.msgIsTarget {
				t.Errorf("IsTargetMessage() = %v %v, msgIsTarget %v, reason %s", msgIsTarget, refMsgIsTarget, tt.msgIsTarget, reason)
			}
			if refMsgIsTarget != tt.refMsgIsTarget {
				t.Errorf("IsTargetMessage() = %v %v, refMsgIsTarget %v, reason %s", msgIsTarget, refMsgIsTarget, tt.refMsgIsTarget, reason)
			}
		})
	}
}

func TestIsTargetMessageEvent(t *testing.T) {
	type args struct {
		botUser      *model.User
		messageEvent *model.MessageEvent
	}
	tests := []struct {
		name           string
		args           args
		msgIsTarget    bool
		refMsgIsTarget bool
	}{
		{
			name: "イベントの対象ユーザがbotでは無かった場合は対象外",
			args: args{
				botUser: &model.User{
					ID: 1,
				},
				messageEvent: &model.MessageEvent{
					TargetUserID: 2,
					Message: &model.Message{
						User: model.User{
							ID: 2,
						},
						MediaNum: 1,
					},
				},
			},
			msgIsTarget:    false,
			refMsgIsTarget: false,
		},
		{
			name: "イベントがShare(twitterの場合はretweet)だった場合は対象外",
			args: args{
				botUser: &model.User{
					ID: 1,
				},
				messageEvent: &model.MessageEvent{
					TargetUserID: 3,
					Message: &model.Message{
						User: model.User{
							ID: 2,
						},
						MediaNum: 1,
					},
				},
			},
			msgIsTarget:    false,
			refMsgIsTarget: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msgIsTarget, refMsgIsTarget, reason := IsTargetMessageEvent(tt.args.botUser, tt.args.messageEvent)
			if msgIsTarget != tt.msgIsTarget {
				t.Errorf("IsTargetMessage() = %v %v, msgIsTarget %v, reason %s", msgIsTarget, refMsgIsTarget, tt.msgIsTarget, reason)
			}
			if refMsgIsTarget != tt.refMsgIsTarget {
				t.Errorf("IsTargetMessage() = %v %v, refMsgIsTarget %v, reason %s", msgIsTarget, refMsgIsTarget, tt.refMsgIsTarget, reason)
			}
		})
	}
}
