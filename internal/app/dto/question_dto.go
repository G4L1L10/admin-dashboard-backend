package dto

type CreateQuestionRequest struct {
	LessonID     string   `json:"lesson_id" binding:"required,uuid"`
	QuestionText string   `json:"question_text" binding:"required"`
	QuestionType string   `json:"question_type" binding:"required"`
	ImageURL     *string  `json:"image_url"`
	AudioURL     *string  `json:"audio_url"`
	Answer       string   `json:"answer" binding:"required"`
	Explanation  string   `json:"explanation"`
	Options      []string `json:"options" binding:"required,dive,required"` // non-empty options
	Tags         []string `json:"tags"`
}

type CreateQuestionResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}
