package repository

import (
	"database/sql"

	"github.com/G4L1L10/admin-dashboard-backend/internal/app/model"
)

type TagRepository struct {
	db *sql.DB
}

func NewTagRepository(db *sql.DB) *TagRepository {
	return &TagRepository{db: db}
}

// CREATE
func (r *TagRepository) Create(tag *model.Tag) error {
	query := `INSERT INTO tags (id, name) VALUES ($1, $2) ON CONFLICT (name) DO NOTHING`
	_, err := r.db.Exec(query, tag.ID, tag.Name)
	return err
}

// READ
func (r *TagRepository) GetByID(id string) (*model.Tag, error) {
	query := `SELECT id, name, created_at, updated_at FROM tags WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var tag model.Tag
	err := row.Scan(&tag.ID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *TagRepository) FindByName(name string) (*model.Tag, error) {
	query := `SELECT id, name, created_at, updated_at FROM tags WHERE name = $1`
	row := r.db.QueryRow(query, name)

	var tag model.Tag
	err := row.Scan(&tag.ID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *TagRepository) SearchByName(keyword string) ([]*model.Tag, error) {
	query := `
	SELECT id, name, created_at, updated_at
	FROM tags
	WHERE name ILIKE '%' || $1 || '%'
	ORDER BY name ASC
	`
	rows, err := r.db.Query(query, keyword)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*model.Tag
	for rows.Next() {
		var tag model.Tag
		err := rows.Scan(&tag.ID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tags = append(tags, &tag)
	}

	return tags, nil
}

func (r *TagRepository) GetByName(name string) (*model.Tag, error) {
	query := `SELECT id, name, created_at, updated_at FROM tags WHERE name = $1`
	row := r.db.QueryRow(query, name)

	var tag model.Tag
	err := row.Scan(&tag.ID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil // tag doesn't exist yet
	}
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// UPDATE
func (r *TagRepository) Update(tag *model.Tag) error {
	query := `UPDATE tags SET name = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(query, tag.Name, tag.ID)
	return err
}

// DELETE
func (r *TagRepository) Delete(id string) error {
	query := `DELETE FROM tags WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// CLEANUP
func (r *TagRepository) DeleteUnusedTags() error {
	query := `
	DELETE FROM tags
	WHERE id NOT IN (
		SELECT DISTINCT tag_id FROM question_tags
	)
	`
	_, err := r.db.Exec(query)
	return err
}
