package repository

import (
	"database/sql"

	"github.com/G4L1L10/admin-dashboard-backend/internal/app/model"
)

type QuestionRepository struct {
	db *sql.DB
}

func NewQuestionRepository(db *sql.DB) *QuestionRepository {
	return &QuestionRepository{db: db}
}

// CREATE
func (r *QuestionRepository) Create(question *model.Question) error {
	query := `INSERT INTO questions (id, lesson_id, question_text, question_type, image_url, audio_url, answer, explanation)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(query, question.ID, question.LessonID, question.QuestionText, question.QuestionType, question.ImageURL, question.AudioURL, question.Answer, question.Explanation)
	return err
}

// READ
func (r *QuestionRepository) GetByID(id string) (*model.Question, error) {
	query := `SELECT id, lesson_id, question_text, question_type, image_url, audio_url, answer, explanation, created_at, updated_at FROM questions WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var question model.Question
	err := row.Scan(&question.ID, &question.LessonID, &question.QuestionText, &question.QuestionType, &question.ImageURL, &question.AudioURL, &question.Answer, &question.Explanation, &question.CreatedAt, &question.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &question, nil
}

func (r *QuestionRepository) GetByLessonID(lessonID string) ([]*model.Question, error) {
	query := `SELECT id, lesson_id, question_text, question_type, image_url, audio_url, answer, explanation, created_at, updated_at 
	          FROM questions WHERE lesson_id = $1 ORDER BY created_at ASC`
	rows, err := r.db.Query(query, lessonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []*model.Question
	for rows.Next() {
		var q model.Question
		err := rows.Scan(&q.ID, &q.LessonID, &q.QuestionText, &q.QuestionType, &q.ImageURL, &q.AudioURL, &q.Answer, &q.Explanation, &q.CreatedAt, &q.UpdatedAt)
		if err != nil {
			return nil, err
		}
		questions = append(questions, &q)
	}
	return questions, nil
}

// UPDATE
func (r *QuestionRepository) Update(question *model.Question) error {
	query := `UPDATE questions SET question_text = $1, question_type = $2, image_url = $3, audio_url = $4, answer = $5, explanation = $6, updated_at = NOW() WHERE id = $7`
	_, err := r.db.Exec(query, question.QuestionText, question.QuestionType, question.ImageURL, question.AudioURL, question.Answer, question.Explanation, question.ID)
	return err
}

// DELETE
func (r *QuestionRepository) Delete(id string) error {
	query := `DELETE FROM questions WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

