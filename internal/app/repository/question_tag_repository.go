package repository

import (
	"database/sql"

	"github.com/G4L1L10/admin-dashboard-backend/internal/app/model"
)

type QuestionTagRepository struct {
	db *sql.DB
}

func NewQuestionTagRepository(db *sql.DB) *QuestionTagRepository {
	return &QuestionTagRepository{db: db}
}

func (r *QuestionTagRepository) AttachTag(questionID, tagID string) error {
	query := `
		INSERT INTO question_tags (question_id, tag_id)
		SELECT $1, $2
		WHERE NOT EXISTS (
			SELECT 1 FROM question_tags WHERE question_id = $1 AND tag_id = $2
		)
	`
	_, err := r.db.Exec(query, questionID, tagID)
	return err
}

// CREATE
func (r *QuestionTagRepository) Create(qt *model.QuestionTag) error {
	query := `INSERT INTO question_tags (question_id, tag_id) VALUES ($1, $2)`
	_, err := r.db.Exec(query, qt.QuestionID, qt.TagID)
	return err
}

// DELETE
func (r *QuestionTagRepository) Delete(questionID, tagID string) error {
	query := `DELETE FROM question_tags WHERE question_id = $1 AND tag_id = $2`
	_, err := r.db.Exec(query, questionID, tagID)
	return err
}

func (r *QuestionTagRepository) DeleteByQuestionID(questionID string) error {
	query := `DELETE FROM question_tags WHERE question_id = $1`
	_, err := r.db.Exec(query, questionID)
	return err
}

func (r *QuestionTagRepository) DeleteAllByQuestionID(questionID string) error {
	query := `DELETE FROM question_tags WHERE question_id = $1`
	_, err := r.db.Exec(query, questionID)
	return err
}

// READ
func (r *QuestionTagRepository) GetTagsByQuestionID(questionID string) ([]string, error) {
	query := `SELECT tag_id FROM question_tags WHERE question_id = $1`
	rows, err := r.db.Query(query, questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tagID string
		if err := rows.Scan(&tagID); err != nil {
			return nil, err
		}
		tags = append(tags, tagID)
	}
	return tags, nil
}
