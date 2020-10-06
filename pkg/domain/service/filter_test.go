package domain

import (
	"testing"

	"github.com/mpppk/sutaba-server/pkg/domain/model"
)

func TestIsTargetMessage(t *testing.T) {
	type args struct {
		user    *model.User
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
				user: &model.User{
					ID: 1,
				},
				message: &model.Message{
					User: model.User{
						ID: 2,
					},
					MediaNum: 1,
				},
			},
			msgIsTarget:    true,
			refMsgIsTarget: false,
		},
		{
			name: "自分自身でない、Mediaを含むメッセージをreferしたメッセージは対象",
			args: args{
				user: &model.User{
					ID: 1,
				},
				message: &model.Message{
					User: model.User{
						ID: 2,
					},
					MediaNum: 0,
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
			name: "メディアが付与されたメッセージでなければ対象外",
			args: args{
				user: &model.User{},
				message: &model.Message{
					MediaNum: 0,
				},
			},
			msgIsTarget: false,
		},
		{
			name: "自分自身のメッセージは対象外",
			args: args{
				user: &model.User{
					ID: 1,
				},
				message: &model.Message{
					User: model.User{
						ID: 1,
					},
					MediaNum: 1,
				},
			},
			msgIsTarget:    false,
			refMsgIsTarget: false,
		},
		{
			name: "自分自身のメッセージをreferしているメッセージは対象外",
			args: args{
				user: &model.User{
					ID: 1,
				},
				message: &model.Message{
					User: model.User{
						ID: 2,
					},
					MediaNum: 1,
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
				user: &model.User{
					ID: 1,
				},
				message: &model.Message{
					User: model.User{
						ID: 1,
					},
					MediaNum: 1,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msgIsTarget, refMsgIsTarget, reason := IsTargetMessage(tt.args.user, tt.args.message)
			if msgIsTarget != tt.msgIsTarget {
				t.Errorf("IsTargetMessage() = %v %v, msgIsTarget %v, reason %s", msgIsTarget, refMsgIsTarget, tt.msgIsTarget, reason)
			}
			if refMsgIsTarget != tt.refMsgIsTarget {
				t.Errorf("IsTargetMessage() = %v %v, refMsgIsTarget %v, reason %s", msgIsTarget, refMsgIsTarget, tt.refMsgIsTarget, reason)
			}
		})
	}
}
