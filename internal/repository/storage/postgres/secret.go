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
	var createdSecret dao.SecretDAO

	query := r.db.QueryBuilder.Insert("secrets").
		Columns("collection_id", "secret_type", "name", "description", "created_by", "updated_by", "linked_secret_id").
		Values(collectionID, secret.SecretType, secret.Name, secret.Description, secret.CreatedBy, secret.UpdatedBy, secret.LinkedSecretId).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&createdSecret.ID,
		&createdSecret.CollectionID,
		&createdSecret.SecretType,
		&createdSecret.Name,
		&createdSecret.Description,
		&createdSecret.CreatedBy,
		&createdSecret.UpdatedBy,
		&createdSecret.CreatedAt,
		&createdSecret.UpdatedAt,
		&createdSecret.LinkedSecretId,
	)

	if err != nil {
		if errCode := r.db.ErrorCode(err); errCode == "23503" {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return &createdSecret, nil
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
		&secret.Name,
		&secret.Description,
		&secret.CreatedBy,
		&secret.UpdatedBy,
		&secret.CreatedAt,
		&secret.UpdatedAt,
		&secret.LinkedSecretId,
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
			&secret.Name,
			&secret.Description,
			&secret.CreatedBy,
			&secret.UpdatedBy,
			&secret.CreatedAt,
			&secret.UpdatedAt,
			&secret.LinkedSecretId,
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
		Set("name", secret.Name).
		Set("description", secret.Description).
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
		&updatedSecret.Name,
		&updatedSecret.Description,
		&updatedSecret.CreatedBy,
		&updatedSecret.UpdatedBy,
		&updatedSecret.CreatedAt,
		&updatedSecret.UpdatedAt,
		&updatedSecret.LinkedSecretId,
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

// Create Password Secret creates a new password secret in the data warehouse.
func (r *SecretRepository) CreatePasswordSecret(ctx context.Context, secret *dao.PasswordSecretDAO) (*dao.PasswordSecretDAO, error) {
	query := r.db.QueryBuilder.Insert("password_secrets").
		Columns("url", "login", "password").
		Values(secret.URL, secret.Login, secret.Password).
		Suffix("RETURNING *")

	var createdSecret dao.PasswordSecretDAO
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&createdSecret.ID,
		&createdSecret.URL,
		&createdSecret.Login,
		&createdSecret.Password,
	)
	if err != nil {
		return nil, err
	}

	return &createdSecret, nil
}

// Get Password Secret By ID returns the password secret by the specified ID.
func (r *SecretRepository) GetPasswordSecretByID(ctx context.Context, id uuid.UUID) (*dao.PasswordSecretDAO, error) {
	var secret dao.PasswordSecretDAO

	query := r.db.QueryBuilder.Select("*").From("password_secrets").Where(sq.Eq{"id": id}).Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&secret.ID,
		&secret.URL,
		&secret.Login,
		&secret.Password,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return &secret, nil
}

// Update Password Secret updates the password secret.
func (r *SecretRepository) UpdatePasswordSecret(ctx context.Context, secret *dao.PasswordSecretDAO) (*dao.PasswordSecretDAO, error) {
	var updatedSecret dao.PasswordSecretDAO

	query := r.db.QueryBuilder.Update("password_secrets").
		Set("url", secret.URL).
		Set("login", secret.Login).
		Set("password", secret.Password).
		Where(sq.Eq{"id": secret.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&updatedSecret.ID,
		&updatedSecret.URL,
		&updatedSecret.Login,
		&updatedSecret.Password,
	)
	if err != nil {
		return nil, err
	}

	return &updatedSecret, nil
}

// CreateTextSecret creates a new text secret in the data warehouse.
func (r *SecretRepository) CreateTextSecret(ctx context.Context, secret *dao.TextSecretDAO) (*dao.TextSecretDAO, error) {
	query := r.db.QueryBuilder.Insert("text_secrets").
		Columns("text").
		Values(secret.Text).
		Suffix("RETURNING *")

	var createdSecret dao.TextSecretDAO
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&createdSecret.ID,
		&createdSecret.Text,
	)
	if err != nil {
		return nil, err
	}

	return &createdSecret, nil
}

// GetTextSecretByID returns a text secret by the specified identifier.
func (r *SecretRepository) GetTextSecretByID(ctx context.Context, id uuid.UUID) (*dao.TextSecretDAO, error) {
	var secret dao.TextSecretDAO

	query := r.db.QueryBuilder.Select("*").From("text_secrets").Where(sq.Eq{"id": id}).Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&secret.ID,
		&secret.Text,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return &secret, nil
}

// UpdateText Secret updates the text secret.
func (r *SecretRepository) UpdateTextSecret(ctx context.Context, secret *dao.TextSecretDAO) (*dao.TextSecretDAO, error) {
	var updatedSecret dao.TextSecretDAO

	query := r.db.QueryBuilder.Update("text_secrets").
		Set("text", secret.Text).
		Where(sq.Eq{"id": secret.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&updatedSecret.ID,
		&updatedSecret.Text,
	)
	if err != nil {
		return nil, err
	}

	return &updatedSecret, nil
}
