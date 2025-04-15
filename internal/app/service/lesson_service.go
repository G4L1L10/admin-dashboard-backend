package service

import (
	"github.com/G4L1L10/admin-dashboard-backend/internal/app/model"
	"github.com/G4L1L10/admin-dashboard-backend/internal/app/repository"
)

type LessonService struct {
	lessonRepo   *repository.LessonRepository
	questionRepo *repository.QuestionRepository
}

func NewLessonService(
	lessonRepo *repository.LessonRepository,
	questionRepo *repository.QuestionRepository,
) *LessonService {
	return &LessonService{
		lessonRepo:   lessonRepo,
		questionRepo: questionRepo,
	}
}

// CREATE
func (s *LessonService) CreateLesson(lesson *model.Lesson) error {
	return s.lessonRepo.Create(lesson)
}

// READ
func (s *LessonService) GetLessonByID(id string) (*model.Lesson, error) {
	return s.lessonRepo.GetByID(id)
}

func (s *LessonService) GetLessonsByCourseID(courseID string) ([]*model.Lesson, error) {
	return s.lessonRepo.GetByCourseID(courseID)
}

func (s *LessonService) GetFullLesson(lessonID string) (*model.FullLesson, error) {
	// Get lesson metadata
	lesson, err := s.lessonRepo.GetByID(lessonID)
	if err != nil {
		return nil, err
	}

	// Get all questions (with options and tags)
	questions, err := s.questionRepo.GetByLessonID(lessonID)
	if err != nil {
		return nil, err
	}

	// Bundle it together
	fullLesson := &model.FullLesson{
		Lesson:    lesson,
		Questions: questions,
	}

	return fullLesson, nil
}

// UPDATE
func (s *LessonService) UpdateLesson(lesson *model.Lesson) error {
	return s.lessonRepo.Update(lesson)
}

// DELETE
func (s *LessonService) DeleteLesson(id string) error {
	return s.lessonRepo.Delete(id)
}
