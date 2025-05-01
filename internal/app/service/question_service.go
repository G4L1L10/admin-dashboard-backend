package service

import (
	"strings"

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
	if err := s.questionRepo.Create(question); err != nil {
		return err
	}

	for _, option := range options {
		if err := s.optionRepo.Create(option); err != nil {
			return err
		}
	}

	// ✅ Deduplicate tag names
	seen := make(map[string]bool)

	for _, tagName := range tags {
		tagName = strings.TrimSpace(tagName)
		if tagName == "" || seen[tagName] {
			continue
		}
		seen[tagName] = true

		tag, err := s.tagRepo.FindByName(tagName)
		if err != nil {
			newTag := &model.Tag{
				ID:   utils.GenerateUUID(),
				Name: tagName,
			}
			if err := s.tagRepo.Create(newTag); err != nil {
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

// Basic question only (no options/tags)
func (s *QuestionService) GetQuestionByID(id string) (*model.Question, error) {
	return s.questionRepo.GetByID(id)
}

// Full question with tags and options
func (s *QuestionService) GetFullQuestionByID(id string) (*model.QuestionWithOptions, error) {
	return s.questionRepo.GetByIDWithTags(id)
}

func (s *QuestionService) GetQuestionsByLessonID(lessonID string) ([]*model.QuestionWithOptions, error) {
	return s.questionRepo.GetByLessonID(lessonID)
}

func (s *QuestionService) GetQuestionsByTag(tagName string) ([]*model.QuestionWithOptions, error) {
	return s.questionRepo.GetQuestionsByTag(tagName)
}

func (s *QuestionService) RemoveAllTagsForQuestion(questionID string) error {
	return s.questionTagRepo.DeleteAllByQuestionID(questionID)
}

func (s *QuestionService) AttachTagToQuestion(questionID, tagID string) error {
	return s.questionTagRepo.AttachTag(questionID, tagID)
}

// UPDATE
func (s *QuestionService) UpdateQuestion(question *model.Question) error {
	// Step 1: update base question fields
	if err := s.questionRepo.Update(question); err != nil {
		return err
	}

	// Step 2: clear existing tag links
	if err := s.questionTagRepo.DeleteByQuestionID(question.ID); err != nil {
		return err
	}

	// Step 3: recreate tag links
	for _, tagName := range question.Tags {
		tag, err := s.tagRepo.FindByName(tagName)
		if err != nil {
			// Tag doesn't exist, create it
			newTag := &model.Tag{
				ID:   utils.GenerateUUID(),
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

	// Step 4: cleanup unused tags
	if err := s.tagRepo.DeleteUnusedTags(); err != nil {
		return err
	}

	return nil
}

// DELETE
func (s *QuestionService) DeleteQuestion(id string) error {
	if err := s.questionRepo.Delete(id); err != nil {
		return err
	}

	// Step 2: cleanup unused tags
	return s.tagRepo.DeleteUnusedTags()
}
