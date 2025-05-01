package model

type QuestionWithOptions struct {
	ID           string   `json:"id"`
	LessonID     string   `json:"lesson_id"`
	QuestionText string   `json:"question_text"`
	QuestionType string   `json:"question_type"`
	ImageURL     *string  `json:"image_url,omitempty"`
	AudioURL     *string  `json:"audio_url,omitempty"`
	Answer       *string  `json:"answer"`
	Explanation  *string  `json:"explanation"`
	Options      []string `json:"options"`
	Tags         []string `json:"tags"`
	Position     int      `json:"position"` // ✅ Add this line
}

