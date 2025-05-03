package postgres

import (
	"backend/internal/entity"
	"context"
	"database/sql"
	"errors"
	"time"
)

type EducationRepository struct {
	db *sql.DB
}

func NewEducationRepository(db *sql.DB) *EducationRepository {
	return &EducationRepository{
		db: db,
	}
}

func (r *EducationRepository) Create(ctx context.Context, education *entity.Education) error {
	query := `
	INSERT INTO educations (image, title, content, type, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	now := time.Now()
	education.CreatedAt = now
	education.UpdatedAt = now

	_, err := r.db.ExecContext(
		ctx,
		query,
		education.Image,
		education.Title,
		education.Content,
		education.Type,
		education.CreatedAt,
		education.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *EducationRepository) GetAllEducation(ctx context.Context) ([]*entity.Education, error) {
	query := `
	SELECT id, image, title, content, type, created_at, updated_at FROM educations
	`

	rows, err := r.db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var educations []*entity.Education
	for rows.Next() {
		e := &entity.Education{}
		err := rows.Scan(
			&e.ID,
			&e.Image,
			&e.Title,
			&e.Content,
			&e.Type,
			&e.CreatedAt,
			&e.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		educations = append(educations, e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return educations, nil
}

func (r *EducationRepository) FindById(ctx context.Context, id uint) (*entity.Education, error) {
	query := `
	SELECT id, image, title, content, type, created_at, updated_at
	FROM educations
	WHERE id = $1
	`

	e := &entity.Education{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&e.ID,
		&e.Image,
		&e.Title,
		&e.Content,
		&e.Type,
		&e.CreatedAt,
		&e.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return e, nil
}

func (r *EducationRepository) FindEducationPendonor(ctx context.Context) ([]*entity.Education, error) {
	query := `
	SELECT id, image, title, content, type, created_at, updated_at
	FROM educations
	WHERE type = 'pendonor'
	`

	rows, err := r.db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var educations []*entity.Education
	for rows.Next() {
		e := &entity.Education{}
		err := rows.Scan(
			&e.ID,
			&e.Image,
			&e.Title,
			&e.Content,
			&e.Type,
			&e.CreatedAt,
			&e.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		educations = append(educations, e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return educations, nil
}

func (r *EducationRepository) FindEducationPencariDonor(ctx context.Context) ([]*entity.Education, error) {
	query := `
	SELECT id, image, title, content, type, created_at, updated_at
	FROM educations
	WHERE type = 'donor'
	`

	rows, err := r.db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var educations []*entity.Education
	for rows.Next() {
		e := &entity.Education{}
		err := rows.Scan(
			&e.ID,
			&e.Image,
			&e.Title,
			&e.Content,
			&e.Type,
			&e.CreatedAt,
			&e.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		educations = append(educations, e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return educations, nil
}

func (r *EducationRepository) Update(ctx context.Context, education *entity.Education) error {
	query := `
	UPDATE educations
	SET image = $1, title = $2, content = $3, type = $4, updated_at = $5
	WHERE id = $6
	`

	education.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(
		ctx,
		query,
		education.Image,
		education.Title,
		education.Content,
		education.Type,
		education.UpdatedAt,
		education.ID,
	)

	return err
}

func (r *EducationRepository) Delete(ctx context.Context, id uint) error {
	query := `
	DELETE FROM educations WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)

	return err
}
