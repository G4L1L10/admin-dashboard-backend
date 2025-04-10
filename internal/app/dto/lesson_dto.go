package dto

type CreateLessonRequest struct {
	CourseID     string `json:"course_id" binding:"required,uuid"`
	Unit         int    `json:"unit" binding:"required"`
	Title        string `json:"title" binding:"required"`
	Description  string `json:"description"`
	Difficulty   string `json:"difficulty"`
	XPReward     int    `json:"xp_reward"`
	CrownsReward int    `json:"crowns_reward"`
}

type CreateLessonResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}
