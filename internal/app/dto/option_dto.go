package dto

type CreateOptionRequest struct {
	QuestionID string `json:"question_id" binding:"required,uuid"`
	OptionText string `json:"option_text" binding:"required"`
}

type CreateOptionResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}
