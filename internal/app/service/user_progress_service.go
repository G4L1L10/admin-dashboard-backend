// internal/app/service/user_progress_service.go
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

// ✅ Get all progress for a user
func (s *UserProgressService) GetUserProgress(userID uuid.UUID) ([]model.UserProgress, error) {
	return s.Repository.GetProgressByUser(userID)
}

// ✅ Mark lesson completed
func (s *UserProgressService) MarkLessonCompleted(input model.MarkProgressInput) error {
	return s.Repository.UpsertUserProgress(input)
}

