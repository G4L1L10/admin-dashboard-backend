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
	// Step 1: Get current max position for this lesson
	var maxPos int
	err := r.db.QueryRow(`
		SELECT COALESCE(MAX(position), 0) 
		FROM questions 
		WHERE lesson_id = $1
	`, question.LessonID).Scan(&maxPos)
	if err != nil {
		return err
	}
	question.Position = maxPos + 1

	// Step 2: Insert question with position
	query := `INSERT INTO questions (id, lesson_id, question_text, question_type, image_url, audio_url, answer, explanation, position)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err = r.db.Exec(query, question.ID, question.LessonID, question.QuestionText, question.QuestionType, question.ImageURL, question.AudioURL, question.Answer, question.Explanation, question.Position)
	return err
}

// READ - basic question only (no tags/options)
func (r *QuestionRepository) GetByID(id string) (*model.Question, error) {
	query := `SELECT id, lesson_id, question_text, question_type, image_url, audio_url, answer, explanation, created_at, updated_at FROM questions WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var q model.Question
	err := row.Scan(&q.ID, &q.LessonID, &q.QuestionText, &q.QuestionType, &q.ImageURL, &q.AudioURL, &q.Answer, &q.Explanation, &q.CreatedAt, &q.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &q, nil
}

// READ - full question with tags and options
func (r *QuestionRepository) GetByIDWithTags(id string) (*model.QuestionWithOptions, error) {
	query := `
	SELECT 
		q.id, q.lesson_id, q.question_text, q.question_type, q.image_url, q.audio_url, q.answer, q.explanation,
		q.position,
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
		q.id = $1
	`

	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var question *model.QuestionWithOptions
	optionSet := make(map[string]bool)
	tagSet := make(map[string]bool)

	for rows.Next() {
		var (
			qID, lessonID, questionText, questionType string
			imageURL, audioURL, answer, explanation   *string
			position                                  int
			optionText, tagName                       *string
		)

		err := rows.Scan(&qID, &lessonID, &questionText, &questionType,
			&imageURL, &audioURL, &answer, &explanation, &position,
			&optionText, &tagName)
		if err != nil {
			return nil, err
		}

		if question == nil {
			question = &model.QuestionWithOptions{
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
				Position:     position,
			}
		}

		if optionText != nil && !optionSet[*optionText] {
			question.Options = append(question.Options, *optionText)
			optionSet[*optionText] = true
		}
		if tagName != nil && !tagSet[*tagName] {
			question.Tags = append(question.Tags, *tagName)
			tagSet[*tagName] = true
		}
	}

	if question == nil {
		return nil, sql.ErrNoRows
	}

	return question, nil
}

func (r *QuestionRepository) GetByLessonID(lessonID string) ([]*model.QuestionWithOptions, error) {
	query := `
	SELECT 
		q.id, q.lesson_id, q.question_text, q.question_type, q.image_url, q.audio_url, q.answer, q.explanation,
		q.position, -- ✅ Includes position
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
		q.position ASC, o.created_at ASC, t.name ASC
	`

	rows, err := r.db.Query(query, lessonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	questionMap := make(map[string]*model.QuestionWithOptions)
	questionList := []*model.QuestionWithOptions{} // ✅ Track order

	for rows.Next() {
		var (
			qID, lessonID, questionText, questionType string
			imageURL, audioURL, answer, explanation   *string
			position                                  int
			optionText, tagName                       *string
		)

		err := rows.Scan(&qID, &lessonID, &questionText, &questionType,
			&imageURL, &audioURL, &answer, &explanation,
			&position, &optionText, &tagName)
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
				Position:     position,
			}
			questionMap[qID] = q
			questionList = append(questionList, q) // ✅ preserve order
		}

		if optionText != nil && !slices.Contains(q.Options, *optionText) {
			q.Options = append(q.Options, *optionText)
		}
		if tagName != nil && !slices.Contains(q.Tags, *tagName) {
			q.Tags = append(q.Tags, *tagName)
		}
	}

	return questionList, nil
}

func (r *QuestionRepository) GetQuestionsByTag(tagName string) ([]*model.QuestionWithOptions, error) {
	query := `
	SELECT 
		q.id, q.lesson_id, q.question_text, q.question_type, q.image_url, q.audio_url, q.answer, q.explanation,
		q.position,
		o.option_text,
		t.name
	FROM 
		questions q
	INNER JOIN 
		question_tags qt ON q.id = qt.question_id
	INNER JOIN 
		tags t ON qt.tag_id = t.id
	LEFT JOIN 
		options o ON q.id = o.question_id
	WHERE 
		t.name = $1
	ORDER BY 
		q.position ASC, o.created_at ASC
	`

	rows, err := r.db.Query(query, tagName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	questionMap := make(map[string]*model.QuestionWithOptions)

	for rows.Next() {
		var (
			qID, lessonID, questionText, questionType string
			imageURL, audioURL, answer, explanation   *string
			position                                  int
			optionText, tagNameFromDB                 *string
		)

		err := rows.Scan(&qID, &lessonID, &questionText, &questionType,
			&imageURL, &audioURL, &answer, &explanation, &position,
			&optionText, &tagNameFromDB)
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
				Position:     position,
			}
			questionMap[qID] = q
		}

		if optionText != nil && !slices.Contains(q.Options, *optionText) {
			q.Options = append(q.Options, *optionText)
		}
		if tagNameFromDB != nil && !slices.Contains(q.Tags, *tagNameFromDB) {
			q.Tags = append(q.Tags, *tagNameFromDB)
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
