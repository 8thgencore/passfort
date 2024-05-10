package postgres

import (
	"context"

	"github.com/8thgencore/passfort/internal/database"
	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/internal/repository/storage/postgres/dao"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

/**
 * SecretRepository implements postgres.SecretRepository interface
 * and provides access to the PostgreSQL database
 */
type SecretRepository struct {
	db *database.DB
}

// NewSecretRepository creates a new user repository instance.
func NewSecretRepository(db *database.DB) *SecretRepository {
	return &SecretRepository{
		db,
	}
}

// CreateSecret creates a new secret in the data warehouse.
func (r *SecretRepository) CreateSecret(ctx context.Context, collectionID uuid.UUID, secret *dao.SecretDAO) (*dao.SecretDAO, error) {
	var secretDao dao.SecretDAO

	query := r.db.QueryBuilder.Insert("secrets").
		Columns("collection_id", "secret_type", "created_by", "updated_by").
		Values(collectionID, secret.SecretType, secret.CreatedBy, secret.UpdatedBy).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&secretDao.ID,
		&secretDao.CollectionID,
		&secretDao.SecretType,
		&secretDao.CreatedBy,
		&secretDao.UpdatedBy,
		&secretDao.CreatedAt,
		&secretDao.UpdatedAt,
	)
	if err != nil {
		if errCode := r.db.ErrorCode(err); errCode == "23503" {
			return nil, domain.ErrCollectionNotExists
		}
		return nil, err
	}

	return &secretDao, nil
}

// GetSecretByID returns the secret by the specified identifier.
func (r *SecretRepository) GetSecretByID(ctx context.Context, id uuid.UUID) (*dao.SecretDAO, error) {
	var secret dao.SecretDAO

	query := r.db.QueryBuilder.Select("*").From("secrets").Where(sq.Eq{"id": id}).Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&secret.ID,
		&secret.CollectionID,
		&secret.SecretType,
		&secret.CreatedBy,
		&secret.UpdatedBy,
		&secret.CreatedAt,
		&secret.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return &secret, nil
}

// ListSecretsByCollectionID selects a list of secrets for a specific collection ID
func (r *SecretRepository) ListSecretsByCollectionID(ctx context.Context, collectionID uuid.UUID, skip, limit uint64) ([]dao.SecretDAO, error) {
	var secrets []dao.SecretDAO

	query := r.db.QueryBuilder.Select("*").From("secrets").
		Where(sq.Eq{"collection_id": collectionID}).
		Offset(skip).Limit(limit)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var secret dao.SecretDAO
		err := rows.Scan(
			&secret.ID,
			&secret.CollectionID,
			&secret.SecretType,
			&secret.CreatedBy,
			&secret.UpdatedBy,
			&secret.CreatedAt,
			&secret.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		secrets = append(secrets, secret)
	}

	return secrets, nil
}

// UpdateSecret updates a secret
func (r *SecretRepository) UpdateSecret(ctx context.Context, secret *dao.SecretDAO) (*dao.SecretDAO, error) {
	var updatedSecret dao.SecretDAO

	query := r.db.QueryBuilder.Update("secrets").
		Set("collection_id", secret.CollectionID).
		Set("secret_type", secret.SecretType).
		Set("updated_at", secret.UpdatedAt).
		Set("updated_by", secret.UpdatedBy).
		Where(sq.Eq{"id": secret.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&updatedSecret.ID,
		&updatedSecret.CollectionID,
		&updatedSecret.SecretType,
		&updatedSecret.CreatedBy,
		&updatedSecret.UpdatedBy,
		&updatedSecret.CreatedAt,
		&updatedSecret.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &updatedSecret, nil
}

// DeleteSecret deletes a secret
func (r *SecretRepository) DeleteSecret(ctx context.Context, id uuid.UUID) error {
	query := r.db.QueryBuilder.Delete("secrets").Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sql, args...)
	return err
}
