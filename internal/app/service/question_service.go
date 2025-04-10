package service

import (
	"github.com/G4L1L10/admin-dashboard-backend/internal/app/model"
	"github.com/G4L1L10/admin-dashboard-backend/internal/app/repository"
	"github.com/G4L1L10/admin-dashboard-backend/pkg/utils"
)

type QuestionService struct {
	questionRepo    *repository.QuestionRepository
	optionRepo      *repository.OptionRepository
	tagRepo         *repository.TagRepository
	questionTagRepo *repository.QuestionTagRepository
}

func NewQuestionService(
	questionRepo *repository.QuestionRepository,
	optionRepo *repository.OptionRepository,
	tagRepo *repository.TagRepository,
	questionTagRepo *repository.QuestionTagRepository,
) *QuestionService {
	return &QuestionService{
		questionRepo:    questionRepo,
		optionRepo:      optionRepo,
		tagRepo:         tagRepo,
		questionTagRepo: questionTagRepo,
	}
}

// CREATE
func (s *QuestionService) CreateQuestion(question *model.Question, options []*model.Option, tags []string) error {
	// Insert Question
	if err := s.questionRepo.Create(question); err != nil {
		return err
	}

	// Insert Options
	for _, option := range options {
		if err := s.optionRepo.Create(option); err != nil {
			return err
		}
	}

	// Insert Tags + Link
	for _, tagName := range tags {
		tag, err := s.tagRepo.FindByName(tagName)
		if err != nil {
			newTag := &model.Tag{
				ID:   utils.GenerateUUID(), // assume you have utils.GenerateUUID()
				Name: tagName,
			}
			err = s.tagRepo.Create(newTag)
			if err != nil {
				return err
			}
			tag = newTag
		}

		questionTag := &model.QuestionTag{
			QuestionID: question.ID,
			TagID:      tag.ID,
		}
		if err := s.questionTagRepo.Create(questionTag); err != nil {
			return err
		}
	}

	return nil
}

// READ
func (s *QuestionService) GetQuestionByID(id string) (*model.Question, error) {
	return s.questionRepo.GetByID(id)
}

func (s *QuestionService) GetQuestionsByLessonID(lessonID string) ([]*model.Question, error) {
	return s.questionRepo.GetByLessonID(lessonID)
}

// UPDATE
func (s *QuestionService) UpdateQuestion(question *model.Question) error {
	return s.questionRepo.Update(question)
}

// DELETE
func (s *QuestionService) DeleteQuestion(id string) error {
	return s.questionRepo.Delete(id)
}

