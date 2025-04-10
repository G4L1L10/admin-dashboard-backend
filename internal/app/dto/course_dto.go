package dto

type CreateCourseRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type CreateCourseResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}
