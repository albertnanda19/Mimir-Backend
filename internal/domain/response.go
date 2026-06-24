package domain

import "time"

type Response struct {
	FormID      string            `json:"form_id"`
	Answers     map[string]string `json:"answers"`
	SubmittedAt time.Time         `json:"submitted_at"`
}
