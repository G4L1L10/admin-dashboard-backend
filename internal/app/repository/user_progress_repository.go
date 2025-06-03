// internal/app/repository/user_progress_repository.go
package repository

import (
	"database/sql"
	"log"
	"time"

	"github.com/G4L1L10/admin-dashboard-backend/internal/app/model"
	"github.com/G4L1L10/admin-dashboard-backend/pkg/utils"
	"github.com/google/uuid"
)

type UserProgressRepository interface {
	GetProgressByUser(userID uuid.UUID) ([]model.UserProgress, error)
	UpsertUserProgress(input model.MarkProgressInput) error
}

type userProgressRepo struct {
	db *sql.DB
}

func NewUserProgressRepository(db *sql.DB) UserProgressRepository {
	return &userProgressRepo{db: db}
}

// ✅ Fetch all lesson progress for a user
func (r *userProgressRepo) GetProgressByUser(userID uuid.UUID) ([]model.UserProgress, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, lesson_id, completed, completed_at, updated_at
		FROM user_progress
		WHERE user_id = $1
	`, userID)
	if err != nil {
		log.Printf("GetProgressByUser query error: %v", err)
		return nil, err
	}
	defer utils.SafeCloseRows(rows)

	var progress []model.UserProgress
	for rows.Next() {
		var p model.UserProgress
		err := rows.Scan(&p.ID, &p.UserID, &p.LessonID, &p.Completed, &p.CompletedAt, &p.UpdatedAt)
		if err != nil {
			log.Printf("GetProgressByUser scan error: %v", err)
			continue
		}
		progress = append(progress, p)
	}

	return progress, nil
}

// ✅ Insert or update a user's lesson progress
func (r *userProgressRepo) UpsertUserProgress(input model.MarkProgressInput) error {
	now := time.Now()
	_, err := r.db.Exec(`
		INSERT INTO user_progress (id, user_id, lesson_id, completed, completed_at, updated_at)
		VALUES ($1, $2, $3, TRUE, $4, $5)
		ON CONFLICT (user_id, lesson_id)
		DO UPDATE SET completed = TRUE, completed_at = $4, updated_at = $5
	`,
		uuid.New(),     // $1 - id
		input.UserID,   // $2
		input.LessonID, // $3
		now,            // $4 - completed_at
		now,            // $5 - updated_at
	)
	if err != nil {
		log.Printf("UpsertUserProgress error: %v", err)
	}

	return err
}

