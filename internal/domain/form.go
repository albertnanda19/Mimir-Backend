package domain

import "time"

type Form struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	CreatedAt time.Time  `json:"created_at"`
	Questions []Question `json:"questions,omitempty"`
}
