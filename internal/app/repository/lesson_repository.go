package repository

import (
	"database/sql"

	"github.com/G4L1L10/admin-dashboard-backend/internal/app/model"
)

type LessonRepository struct {
	db *sql.DB
}

func NewLessonRepository(db *sql.DB) *LessonRepository {
	return &LessonRepository{db: db}
}

// CREATE
func (r *LessonRepository) Create(lesson *model.Lesson) error {
	query := `INSERT INTO lessons (id, course_id, unit, title, description, difficulty, xp_reward, crowns_reward) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(query, lesson.ID, lesson.CourseID, lesson.Unit, lesson.Title, lesson.Description, lesson.Difficulty, lesson.XPReward, lesson.CrownsReward)
	return err
}

// READ
func (r *LessonRepository) GetByID(id string) (*model.Lesson, error) {
	query := `SELECT id, course_id, unit, title, description, difficulty, xp_reward, crowns_reward, created_at, updated_at FROM lessons WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var lesson model.Lesson
	err := row.Scan(
		&lesson.ID,
		&lesson.CourseID,
		&lesson.Unit,
		&lesson.Title,
		&lesson.Description,
		&lesson.Difficulty,
		&lesson.XPReward,
		&lesson.CrownsReward,
		&lesson.CreatedAt,
		&lesson.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &lesson, nil
}

// READ - List lessons by course ID (🆕)
func (r *LessonRepository) GetByCourseID(courseID string) ([]*model.Lesson, error) {
	query := `SELECT id, course_id, unit, title, description, difficulty, xp_reward, crowns_reward, created_at, updated_at 
	          FROM lessons WHERE course_id = $1 ORDER BY unit ASC`
	rows, err := r.db.Query(query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []*model.Lesson
	for rows.Next() {
		var lesson model.Lesson
		err := rows.Scan(
			&lesson.ID,
			&lesson.CourseID,
			&lesson.Unit,
			&lesson.Title,
			&lesson.Description,
			&lesson.Difficulty,
			&lesson.XPReward,
			&lesson.CrownsReward,
			&lesson.CreatedAt,
			&lesson.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		lessons = append(lessons, &lesson)
	}

	return lessons, nil
}

// UPDATE
func (r *LessonRepository) Update(lesson *model.Lesson) error {
	query := `UPDATE lessons 
	          SET title = $1, description = $2, unit = $3, difficulty = $4, xp_reward = $5, crowns_reward = $6, updated_at = NOW() 
	          WHERE id = $7`
	_, err := r.db.Exec(
		query,
		lesson.Title,
		lesson.Description,
		lesson.Unit, // ✅ Now correctly included and ordered
		lesson.Difficulty,
		lesson.XPReward,
		lesson.CrownsReward,
		lesson.ID,
	)
	return err
}

// DELETE
func (r *LessonRepository) Delete(id string) error {
	query := `DELETE FROM lessons WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

