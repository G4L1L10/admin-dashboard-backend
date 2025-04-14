package repository

import (
	"database/sql"

	"github.com/G4L1L10/admin-dashboard-backend/internal/app/model"
)

type CourseRepository struct {
	db *sql.DB
}

func NewCourseRepository(db *sql.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

// CREATE
func (r *CourseRepository) Create(course *model.Course) error {
	query := `INSERT INTO courses (id, title, description) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, course.ID, course.Title, course.Description)
	return err
}

// READ - Single Course
func (r *CourseRepository) GetByID(id string) (*model.Course, error) {
	query := `SELECT id, title, description, created_at, updated_at FROM courses WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var course model.Course
	err := row.Scan(&course.ID, &course.Title, &course.Description, &course.CreatedAt, &course.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &course, nil
}

// READ - List All Courses ✅
func (r *CourseRepository) ListCourses() ([]*model.Course, error) {
	query := `SELECT id, title, description, created_at, updated_at FROM courses ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []*model.Course
	for rows.Next() {
		var course model.Course
		err := rows.Scan(&course.ID, &course.Title, &course.Description, &course.CreatedAt, &course.UpdatedAt)
		if err != nil {
			return nil, err
		}
		courses = append(courses, &course)
	}

	return courses, nil
}

// UPDATE
func (r *CourseRepository) Update(course *model.Course) error {
	query := `UPDATE courses SET title = $1, description = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.db.Exec(query, course.Title, course.Description, course.ID)
	return err
}

// DELETE
func (r *CourseRepository) Delete(id string) error {
	query := `DELETE FROM courses WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
