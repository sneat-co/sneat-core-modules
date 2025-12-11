package models4auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEmail(t *testing.T) {
	type args struct {
		id   int64
		data *EmailData
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "empty",
			args: args{
				id: 0,
				data: &EmailData{
					From:     "me@example.com",
					To:       "stranger@example.com",
					Subject:  "Hello stranger",
					BodyText: "Hello stranger, how are you?",
					BodyHtml: "<p>Hello <b>stranger</b>, how are you?</p>",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email := NewEmail(tt.args.id, tt.args.data)
			assert.Equal(t, tt.args.id, email.ID)
			assert.Equal(t, tt.args.data.From, email.Data.From)
			assert.Equal(t, tt.args.data.To, email.Data.To)
			assert.Equal(t, tt.args.data.Subject, email.Data.Subject)
			assert.Equal(t, tt.args.data.BodyText, email.Data.BodyText)
			assert.Equal(t, tt.args.data.BodyHtml, email.Data.BodyHtml)
		})
	}
}
