package service

import (
	"github.com/G4L1L10/admin-dashboard-backend/internal/app/model"
	"github.com/G4L1L10/admin-dashboard-backend/internal/app/repository"
	"github.com/G4L1L10/admin-dashboard-backend/pkg/utils"
)

type TagService struct {
	tagRepo *repository.TagRepository
}

func NewTagService(tagRepo *repository.TagRepository) *TagService {
	return &TagService{tagRepo: tagRepo}
}

// CREATE
func (s *TagService) CreateTag(tag *model.Tag) error {
	return s.tagRepo.Create(tag)
}

// READ
func (s *TagService) GetTagByID(id string) (*model.Tag, error) {
	return s.tagRepo.GetByID(id)
}

// SEARCH FOR TAGS
func (s *TagService) SearchTags(keyword string) ([]*model.Tag, error) {
	return s.tagRepo.SearchByName(keyword)
}

func (s *TagService) FindOrCreate(name string) (string, error) {
	existing, err := s.tagRepo.GetByName(name)
	if err != nil {
		return "", err
	}
	if existing != nil {
		return existing.ID, nil
	}

	// Not found, so create
	tag := &model.Tag{
		ID:   utils.GenerateUUID(),
		Name: name,
	}
	if err := s.tagRepo.Create(tag); err != nil {
		return "", err
	}
	return tag.ID, nil
}

// UPDATE
func (s *TagService) UpdateTag(tag *model.Tag) error {
	return s.tagRepo.Update(tag)
}

// DELETE
func (s *TagService) DeleteTag(id string) error {
	return s.tagRepo.Delete(id)
}
