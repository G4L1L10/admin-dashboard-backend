package service

import (
	"github.com/G4L1L10/admin-dashboard-backend/internal/app/model"
	"github.com/G4L1L10/admin-dashboard-backend/internal/app/repository"
	"github.com/google/uuid"
)

type UserProgressService struct {
	Repository repository.UserProgressRepository
}

// ✅ Constructor
func NewUserProgressService(repo repository.UserProgressRepository) *UserProgressService {
	return &UserProgressService{
		Repository: repo,
	}
}

// ✅ Replace slice with a summary object
func (s *UserProgressService) GetUserProgress(userID uuid.UUID) (*model.UserProgressSummary, error) {
	// TEMP: hardcoded dummy data
	return &model.UserProgressSummary{
		XP:     150,
		Streak: 5,
		Hearts: 3,
		Crowns: 2,
	}, nil
}

// ✅ Mark lesson completed
func (s *UserProgressService) MarkLessonCompleted(input model.MarkProgressInput) error {
	return s.Repository.UpsertUserProgress(input)
}

