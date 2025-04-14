package repository

import (
	"database/sql"

	"github.com/G4L1L10/admin-dashboard-backend/internal/app/model"
)

type StatsRepository struct {
	db *sql.DB
}

func NewStatsRepository(db *sql.DB) *StatsRepository {
	return &StatsRepository{db: db}
}

func (r *StatsRepository) GetStats() (*model.Stats, error) {
	var stats model.Stats

	err := r.db.QueryRow(`SELECT COUNT(*) FROM courses`).Scan(&stats.Courses)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(`SELECT COUNT(*) FROM lessons`).Scan(&stats.Lessons)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(`SELECT COUNT(*) FROM questions`).Scan(&stats.Questions)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(`SELECT COUNT(*) FROM tags`).Scan(&stats.Tags)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}
