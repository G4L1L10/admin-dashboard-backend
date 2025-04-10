package service

import (
	"github.com/G4L1L10/admin-dashboard-backend/internal/app/model"
	"github.com/G4L1L10/admin-dashboard-backend/internal/app/repository"
)

type LessonService struct {
	lessonRepo *repository.LessonRepository
}

func NewLessonService(lessonRepo *repository.LessonRepository) *LessonService {
	return &LessonService{lessonRepo: lessonRepo}
}

// CREATE
func (s *LessonService) CreateLesson(lesson *model.Lesson) error {
	return s.lessonRepo.Create(lesson)
}

// READ
func (s *LessonService) GetLessonByID(id string) (*model.Lesson, error) {
	return s.lessonRepo.GetByID(id)
}

// UPDATE
func (s *LessonService) UpdateLesson(lesson *model.Lesson) error {
	return s.lessonRepo.Update(lesson)
}

// DELETE
func (s *LessonService) DeleteLesson(id string) error {
	return s.lessonRepo.Delete(id)
}

