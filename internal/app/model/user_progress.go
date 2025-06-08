package model

import (
	"time"

	"github.com/google/uuid"
)

type MarkProgressInput struct {
	UserID   uuid.UUID `json:"user_id" binding:"required"`
	LessonID uuid.UUID `json:"lesson_id" binding:"required"`
}

type UserProgress struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	LessonID    uuid.UUID  `json:"lesson_id"`
	Completed   bool       `json:"completed"`
	CompletedAt *time.Time `json:"completed_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ✅ This is the struct your iOS app expects
type UserProgressSummary struct {
	XP     int `json:"xp"`
	Streak int `json:"streak"`
	Hearts int `json:"hearts"`
	Crowns int `json:"crowns"`
}

