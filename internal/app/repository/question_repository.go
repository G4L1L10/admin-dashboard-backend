package repository

import (
	"database/sql"
	"slices"

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

func (r *QuestionRepository) GetByLessonID(lessonID string) ([]*model.QuestionWithOptions, error) {
	query := `
	SELECT 
		q.id, q.lesson_id, q.question_text, q.question_type, q.image_url, q.audio_url, q.answer, q.explanation,
		o.option_text,
		t.name
	FROM 
		questions q
	LEFT JOIN 
		options o ON q.id = o.question_id
	LEFT JOIN 
		question_tags qt ON q.id = qt.question_id
	LEFT JOIN 
		tags t ON qt.tag_id = t.id
	WHERE 
		q.lesson_id = $1
	ORDER BY 
		q.created_at ASC, o.created_at ASC, t.name ASC
	`

	rows, err := r.db.Query(query, lessonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	questionMap := make(map[string]*model.QuestionWithOptions)

	for rows.Next() {
		var (
			qID, lessonID, questionText, questionType string
			imageURL, audioURL, answer, explanation   *string
			optionText, tagName                       *string
		)

		err := rows.Scan(&qID, &lessonID, &questionText, &questionType, &imageURL, &audioURL, &answer, &explanation, &optionText, &tagName)
		if err != nil {
			return nil, err
		}

		q, exists := questionMap[qID]
		if !exists {
			q = &model.QuestionWithOptions{
				ID:           qID,
				LessonID:     lessonID,
				QuestionText: questionText,
				QuestionType: questionType,
				ImageURL:     imageURL,
				AudioURL:     audioURL,
				Answer:       answer,
				Explanation:  explanation,
				Options:      []string{},
				Tags:         []string{},
			}
			questionMap[qID] = q
		}

		// Collect Options
		if optionText != nil && !slices.Contains(q.Options, *optionText) {
			q.Options = append(q.Options, *optionText)
		}

		// Collect Tags
		if tagName != nil && !slices.Contains(q.Tags, *tagName) {
			q.Tags = append(q.Tags, *tagName)
		}
	}

	var questions []*model.QuestionWithOptions
	for _, q := range questionMap {
		questions = append(questions, q)
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

