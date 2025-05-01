package repository

import (
	"database/sql"

	"github.com/G4L1L10/admin-dashboard-backend/internal/app/model"
)

type OptionRepository struct {
	db *sql.DB
}

func NewOptionRepository(db *sql.DB) *OptionRepository {
	return &OptionRepository{db: db}
}

// CREATE
func (r *OptionRepository) Create(option *model.Option) error {
	query := `INSERT INTO options (id, question_id, option_text) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, option.ID, option.QuestionID, option.OptionText)
	return err
}

// READ
func (r *OptionRepository) GetByID(id string) (*model.Option, error) {
	query := `SELECT id, question_id, option_text, created_at, updated_at FROM options WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var option model.Option
	err := row.Scan(&option.ID, &option.QuestionID, &option.OptionText, &option.CreatedAt, &option.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &option, nil
}

// UPDATE
func (r *OptionRepository) Update(option *model.Option) error {
	query := `UPDATE options SET option_text = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(query, option.OptionText, option.ID)
	return err
}

// Delete all options for a specific question
func (r *OptionRepository) DeleteByQuestionID(questionID string) error {
	query := `DELETE FROM options WHERE question_id = $1`
	_, err := r.db.Exec(query, questionID)
	return err
}

// DELETE
func (r *OptionRepository) Delete(id string) error {
	query := `DELETE FROM options WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
