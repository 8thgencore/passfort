package postgres

import (
	"context"
	"time"

	"github.com/8thgencore/passfort/internal/database"
	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/internal/repository/storage/postgres/dao"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

/**
 * CollectionRepository implements postgres.CollectionRepository interface
 * and provides access to the PostgreSQL database
 */
type CollectionRepository struct {
	db *database.DB
}

// NewCollectionRepository creates a new user repository instance
func NewCollectionRepository(db *database.DB) *CollectionRepository {
	return &CollectionRepository{
		db,
	}
}

// CreateCollection creates a new collection in the database
func (r *CollectionRepository) CreateCollection(ctx context.Context, userID uuid.UUID, collection *dao.CollectionDAO) (*dao.CollectionDAO, error) {
	var collectionDAO dao.CollectionDAO

	// Begin a transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Insert into collections table
	collectionQuery := r.db.QueryBuilder.Insert("collections").
		Columns("name", "description", "created_by", "updated_by").
		Values(collection.Name, collection.Description, collection.CreatedBy, collection.UpdatedBy).
		Suffix("RETURNING *")

	collectionSQL, collectionArgs, err := collectionQuery.ToSql()
	if err != nil {
		return nil, err
	}

	err = tx.QueryRow(ctx, collectionSQL, collectionArgs...).Scan(
		&collectionDAO.ID,
		&collectionDAO.Name,
		&collectionDAO.Description,
		&collectionDAO.CreatedBy,
		&collectionDAO.UpdatedBy,
		&collectionDAO.CreatedAt,
		&collectionDAO.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Insert into users_collections table
	usersCollectionsQuery := r.db.QueryBuilder.Insert("users_collections").
		Columns("user_id", "collection_id").
		Values(userID, collectionDAO.ID).
		Suffix("RETURNING *")

	usersCollectionsSQL, usersCollectionsArgs, err := usersCollectionsQuery.ToSql()
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(ctx, usersCollectionsSQL, usersCollectionsArgs...)
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &collectionDAO, nil
}

// GetCollectionByID gets a collection by ID from the database
func (r *CollectionRepository) GetCollectionByID(ctx context.Context, id uuid.UUID) (*dao.CollectionDAO, error) {
	var collectionDAO dao.CollectionDAO

	query := r.db.QueryBuilder.Select("*").
		From("collections").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&collectionDAO.ID,
		&collectionDAO.Name,
		&collectionDAO.Description,
		&collectionDAO.CreatedBy,
		&collectionDAO.UpdatedBy,
		&collectionDAO.CreatedAt,
		&collectionDAO.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return &collectionDAO, nil
}

// ListCollectionsByUserID lists all collections from the database
func (r *CollectionRepository) ListCollectionsByUserID(ctx context.Context, userID uuid.UUID, skip, limit uint64) ([]dao.CollectionDAO, error) {
	var collectionDAO dao.CollectionDAO
	var collectionsDAO []dao.CollectionDAO

	query := r.db.QueryBuilder.Select("c.id, c.name, c.description, c.created_by, c.updated_by, c.created_at, c.updated_at").
		From("collections c").
		Join("users_collections uc ON c.id = uc.collection_id").
		Where(sq.Eq{"uc.user_id": userID}).
		OrderBy("c.created_at DESC").
		Limit(limit).
		Offset((skip - 1) * limit)

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
		err := rows.Scan(
			&collectionDAO.ID,
			&collectionDAO.Name,
			&collectionDAO.Description,
			&collectionDAO.CreatedBy,
			&collectionDAO.UpdatedBy,
			&collectionDAO.CreatedAt,
			&collectionDAO.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Append the converted collection to the list
		collectionsDAO = append(collectionsDAO, collectionDAO)
	}

	return collectionsDAO, nil
}

// UpdateCollection updates a collection by ID in the database
func (r *CollectionRepository) UpdateCollection(ctx context.Context, collection *dao.CollectionDAO) (*dao.CollectionDAO, error) {
	var collectionDAO dao.CollectionDAO

	name := NullString(collection.Name)
	description := NullString(collection.Description)

	query := r.db.QueryBuilder.Update("collections").
		Set("name", sq.Expr("COALESCE(?, name)", name)).
		Set("description", sq.Expr("COALESCE(?, description)", description)).
		Set("updated_at", time.Now()).
		Set("updated_by", collection.UpdatedBy).
		Where(sq.Eq{"id": collection.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&collectionDAO.ID,
		&collectionDAO.Name,
		&collectionDAO.Description,
		&collectionDAO.CreatedBy,
		&collectionDAO.UpdatedBy,
		&collectionDAO.CreatedAt,
		&collectionDAO.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &collectionDAO, nil
}

// DeleteCollection deletes a collection by ID from the database
func (r *CollectionRepository) DeleteCollection(ctx context.Context, id uuid.UUID) error {
	// Begin a transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Delete links in users_collections associated with the collection
	usersCollectionsQuery := r.db.QueryBuilder.Delete("users_collections").
		Where(sq.Eq{"collection_id": id})

	usersCollectionsSQL, usersCollectionsArgs, err := usersCollectionsQuery.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, usersCollectionsSQL, usersCollectionsArgs...)
	if err != nil {
		return err
	}

	// Delete the collection
	collectionQuery := r.db.QueryBuilder.Delete("collections").
		Where(sq.Eq{"id": id})

	collectionSQL, collectionArgs, err := collectionQuery.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, collectionSQL, collectionArgs...)
	if err != nil {
		return err
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

// IsUserPartOfCollection checks if the user is part of the specified collection
func (r *CollectionRepository) IsUserPartOfCollection(ctx context.Context, userID, collectionID uuid.UUID) (bool, error) {
	query := r.db.QueryBuilder.Select("1").
		From("users_collections").
		Where(sq.Eq{"user_id": userID, "collection_id": collectionID}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return false, err
	}

	var exists interface{}
	err = r.db.QueryRow(ctx, sql, args...).Scan(&exists)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil // User is not part of the collection
		}
		return false, err
	}

	return true, nil // User is part of the collection
}
