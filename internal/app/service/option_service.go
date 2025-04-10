package service

import (
	"github.com/G4L1L10/admin-dashboard-backend/internal/app/model"
	"github.com/G4L1L10/admin-dashboard-backend/internal/app/repository"
)

type OptionService struct {
	optionRepo *repository.OptionRepository
}

func NewOptionService(optionRepo *repository.OptionRepository) *OptionService {
	return &OptionService{optionRepo: optionRepo}
}

// CREATE
func (s *OptionService) CreateOption(option *model.Option) error {
	return s.optionRepo.Create(option)
}

// READ
func (s *OptionService) GetOptionByID(id string) (*model.Option, error) {
	return s.optionRepo.GetByID(id)
}

// UPDATE
func (s *OptionService) UpdateOption(option *model.Option) error {
	return s.optionRepo.Update(option)
}

// DELETE
func (s *OptionService) DeleteOption(id string) error {
	return s.optionRepo.Delete(id)
}

