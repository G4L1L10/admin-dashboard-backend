package service

import (
	"github.com/G4L1L10/admin-dashboard-backend/internal/app/model"
	"github.com/G4L1L10/admin-dashboard-backend/internal/app/repository"
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

// UPDATE
func (s *TagService) UpdateTag(tag *model.Tag) error {
	return s.tagRepo.Update(tag)
}

// DELETE
func (s *TagService) DeleteTag(id string) error {
	return s.tagRepo.Delete(id)
}

