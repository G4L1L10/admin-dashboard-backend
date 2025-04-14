package service

import (
	"github.com/G4L1L10/admin-dashboard-backend/internal/app/model"
	"github.com/G4L1L10/admin-dashboard-backend/internal/app/repository"
)

type StatsService struct {
	statsRepository *repository.StatsRepository
}

func NewStatsService(statsRepository *repository.StatsRepository) *StatsService {
	return &StatsService{statsRepository: statsRepository}
}

func (s *StatsService) GetStats() (*model.Stats, error) {
	return s.statsRepository.GetStats()
}

