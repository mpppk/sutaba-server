package model

import "testing"

func TestMessageID_ToString(t *testing.T) {
	tests := []struct {
		name string
		m    MessageID
		want string
	}{
		{
			m:    1,
			want: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.ToString(); got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessageText_HasKeyword(t *testing.T) {
	type args struct {
		keyword string
	}
	tests := []struct {
		name string
		m    MessageText
		args args
		want bool
	}{
		{
			m:    "text",
			args: args{keyword: "keyword"},
			want: false,
		},
		{
			m:    "text keyword",
			args: args{keyword: "keyword"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.HasKeyword(tt.args.keyword); got != tt.want {
				t.Errorf("HasKeyword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessage_GetIDStr(t *testing.T) {
	type fields struct {
		ID                MessageID
		User              User
		Text              MessageText
		ReferencedMessage *Message
		MediaNum          int
		ReplyUserID       UserID
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{
				ID: 1,
			},
			want: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Message{
				ID:                tt.fields.ID,
				User:              tt.fields.User,
				Text:              tt.fields.Text,
				ReferencedMessage: tt.fields.ReferencedMessage,
				MediaNum:          tt.fields.MediaNum,
				ReplyUserID:       tt.fields.ReplyUserID,
			}
			if got := m.GetIDStr(); got != tt.want {
				t.Errorf("GetIDStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessage_HasKeyWord(t *testing.T) {
	type fields struct {
		ID                MessageID
		User              User
		Text              MessageText
		ReferencedMessage *Message
		MediaNum          int
		ReplyUserID       UserID
	}
	type args struct {
		keyword string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			fields: fields{Text: "text"},
			args:   args{keyword: "keyword"},
			want:   false,
		},
		{
			fields: fields{Text: "text keyword"},
			args:   args{keyword: "keyword"},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Message{
				ID:                tt.fields.ID,
				User:              tt.fields.User,
				Text:              tt.fields.Text,
				ReferencedMessage: tt.fields.ReferencedMessage,
				MediaNum:          tt.fields.MediaNum,
				ReplyUserID:       tt.fields.ReplyUserID,
			}
			if got := m.HasKeyWord(tt.args.keyword); got != tt.want {
				t.Errorf("HasKeyWord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessage_HasMessageReference(t *testing.T) {
	type fields struct {
		ID                MessageID
		User              User
		Text              MessageText
		ReferencedMessage *Message
		MediaNum          int
		ReplyUserID       UserID
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "",
			fields: fields{ReferencedMessage: &Message{}},
			want:   true,
		},
		{
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Message{
				ID:                tt.fields.ID,
				User:              tt.fields.User,
				Text:              tt.fields.Text,
				ReferencedMessage: tt.fields.ReferencedMessage,
				MediaNum:          tt.fields.MediaNum,
				ReplyUserID:       tt.fields.ReplyUserID,
			}
			if got := m.HasMessageReference(); got != tt.want {
				t.Errorf("HasMessageReference() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessage_IsDebugMode(t *testing.T) {
	type fields struct {
		ID                MessageID
		User              User
		Text              MessageText
		ReferencedMessage *Message
		MediaNum          int
		ReplyUserID       UserID
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			fields: fields{Text: "text"},
			want:   false,
		},
		{
			fields: fields{Text: MessageText("text " + debugKeyword)},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Message{
				ID:                tt.fields.ID,
				User:              tt.fields.User,
				Text:              tt.fields.Text,
				ReferencedMessage: tt.fields.ReferencedMessage,
				MediaNum:          tt.fields.MediaNum,
				ReplyUserID:       tt.fields.ReplyUserID,
			}
			if got := m.IsDebugMode(); got != tt.want {
				t.Errorf("IsDebugMode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessage_IsRepliedTo(t *testing.T) {
	type fields struct {
		ID                MessageID
		User              User
		Text              MessageText
		ReferencedMessage *Message
		MediaNum          int
		ReplyUserID       UserID
	}
	type args struct {
		user *User
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "ユーザIDが一致する場合はtrue",
			fields: fields{ReplyUserID: 1, Text: "text"},
			args:   args{user: &User{ID: 1}},
			want:   true,
		},
		{
			name:   "@ユーザ名をテキストに含む場合はtrue",
			fields: fields{ReplyUserID: 0, Text: "@username text"},
			args:   args{user: &User{ID: 1, Name: "username"}},
			want:   true,
		},
		{
			name:   "ユーザIDが一致せず、@ユーザ名をテキストに含まない場合はtrue",
			fields: fields{ReplyUserID: 0, Text: "@username text"},
			args:   args{user: &User{ID: 1, Name: "bob"}},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Message{
				ID:                tt.fields.ID,
				User:              tt.fields.User,
				Text:              tt.fields.Text,
				ReferencedMessage: tt.fields.ReferencedMessage,
				MediaNum:          tt.fields.MediaNum,
				ReplyUserID:       tt.fields.ReplyUserID,
			}
			if got := m.IsRepliedTo(tt.args.user); got != tt.want {
				t.Errorf("IsRepliedTo() = %v, want %v", got, tt.want)
			}
		})
	}
}
