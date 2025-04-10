package dto

type CreateTagRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateTagResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}
