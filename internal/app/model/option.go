package model

import "time"

type Option struct {
	ID         string    `json:"id"`
	QuestionID string    `json:"question_id"`
	OptionText string    `json:"option_text"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
