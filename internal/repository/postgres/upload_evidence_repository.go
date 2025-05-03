package postgres

import (
	"backend/internal/entity"
	"context"
	"database/sql"
)

type UploadEvidenceRepository struct {
	db *sql.DB
}

func NewUploadEvidenceRepository(db *sql.DB) *UploadEvidenceRepository {
	return &UploadEvidenceRepository{
		db: db,
	}
}

func (r *UploadEvidenceRepository) Upload(ctx context.Context, uploadEvidence *entity.UploadEvidence) error {
	query := `
	INSERT INTO images (user_id, image, created_at)
	VALUES ($1, $2, $3)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		uploadEvidence.UserID,
		uploadEvidence.Image,
		uploadEvidence.CreatedAt,
	)

	return err
}