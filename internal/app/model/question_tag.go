package model

import "time"

type QuestionTag struct {
	QuestionID string    `json:"question_id"`
	TagID      string    `json:"tag_id"`
	CreatedAt  time.Time `json:"created_at"`
}
