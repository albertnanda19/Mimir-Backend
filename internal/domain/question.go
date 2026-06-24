package domain

type Question struct {
	ID      string   `json:"id"`
	Label   string   `json:"label"`
	Type    string   `json:"type"` // text, multiple_choice, rating, file
	Options []string `json:"options,omitempty"`
}
