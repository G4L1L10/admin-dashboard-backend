package service

import (
	"context"
	"fmt"
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
	ctx := context.Background()
	const bucketName = "cms-media-1" // replace with your actual bucket if different

	// Step 0: Get current question record for comparison
	oldQuestion, err := s.questionRepo.GetByID(question.ID)
	if err != nil {
		return err
	}

	// Step 1: Compare image and delete old if changed
	if oldQuestion.ImageURL != nil && question.ImageURL != nil && *oldQuestion.ImageURL != *question.ImageURL {
		go func(oldPath string) {
			if err := utils.DeleteFromGCS(ctx, bucketName, oldPath); err != nil {
				fmt.Printf("Failed to delete old image: %v\n", err)
			}
		}(*oldQuestion.ImageURL)
	}

	// Step 2: Compare audio and delete old if changed
	if oldQuestion.AudioURL != nil && question.AudioURL != nil && *oldQuestion.AudioURL != *question.AudioURL {
		go func(oldPath string) {
			if err := utils.DeleteFromGCS(ctx, bucketName, oldPath); err != nil {
				fmt.Printf("Failed to delete old audio: %v\n", err)
			}
		}(*oldQuestion.AudioURL)
	}

	// Step 3: Update question fields in DB
	if err := s.questionRepo.Update(question); err != nil {
		return err
	}

	// Step 4: Remove existing tag links
	if err := s.questionTagRepo.DeleteByQuestionID(question.ID); err != nil {
		return err
	}

	// Step 5: Recreate tag links
	for _, tagName := range question.Tags {
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

	// Step 6: Cleanup unused tags
	if err := s.tagRepo.DeleteUnusedTags(); err != nil {
		return err
	}

	return nil
}

// DELETE
func (s *QuestionService) DeleteQuestion(id string) error {
	ctx := context.Background()
	const bucketName = "cms-media-1" // change if needed

	// Step 1: Fetch the question to get media references
	q, err := s.questionRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Step 2: Delete the DB record
	if err := s.questionRepo.Delete(id); err != nil {
		return err
	}

	// Step 3: Remove associated media if present
	if q.ImageURL != nil && *q.ImageURL != "" {
		go func(path string) {
			if err := utils.DeleteFromGCS(ctx, bucketName, path); err != nil {
				fmt.Printf("Failed to delete image on question deletion: %v\n", err)
			}
		}(*q.ImageURL)
	}

	if q.AudioURL != nil && *q.AudioURL != "" {
		go func(path string) {
			if err := utils.DeleteFromGCS(ctx, bucketName, path); err != nil {
				fmt.Printf("Failed to delete audio on question deletion: %v\n", err)
			}
		}(*q.AudioURL)
	}

	// Step 4: Cleanup unused tags
	return s.tagRepo.DeleteUnusedTags()
}
