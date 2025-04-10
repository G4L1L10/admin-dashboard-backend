package model

import "time"

type Question struct {
	ID           string    `json:"id"`
	LessonID     string    `json:"lesson_id"`
	QuestionText string    `json:"question_text"`
	QuestionType string    `json:"question_type"`
	ImageURL     *string   `json:"image_url,omitempty"` // optional
	AudioURL     *string   `json:"audio_url,omitempty"` // optional
	Answer       string    `json:"answer"`
	Explanation  string    `json:"explanation"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
