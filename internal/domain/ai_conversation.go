package domain

import "time"

type AIConversation struct {
	ID        string    `json:"id"`
	FormID    string    `json:"form_id"`
	Messages  []Message `json:"messages"`
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	Role    string `json:"role"` // user, assistant
	Content string `json:"content"`
}
