package model

import "time"

type Lesson struct {
	ID           string    `json:"id"`
	CourseID     string    `json:"course_id"`
	Unit         int       `json:"unit"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Difficulty   string    `json:"difficulty"`
	XPReward     int       `json:"xp_reward"`
	CrownsReward int       `json:"crowns_reward"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
