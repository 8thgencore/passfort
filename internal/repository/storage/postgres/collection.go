package postgres

import (
	"context"
	"time"

	"github.com/8thgencore/passfort/internal/database"
	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/internal/repository/storage/postgres/converter"
	"github.com/8thgencore/passfort/internal/repository/storage/postgres/dao"
	sq "github.com/Masterminds/squirrel"
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
func (r *CollectionRepository) CreateCollection(ctx context.Context, userID uint64, collection *domain.Collection) (*domain.Collection, error) {
	var collectionDAO dao.CollectionDAO

	// Begin a transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Insert into collections table
	collectionQuery := r.db.QueryBuilder.Insert("collections").
		Columns("name", "description").
		Values(collection.Name, collection.Description).
		Suffix("RETURNING *")

	collectionSQL, collectionArgs, err := collectionQuery.ToSql()
	if err != nil {
		return nil, err
	}

	err = tx.QueryRow(ctx, collectionSQL, collectionArgs...).Scan(
		&collectionDAO.ID,
		&collectionDAO.Name,
		&collectionDAO.Description,
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

	return converter.ToCollection(&collectionDAO), nil
}

// GetCollectionByID gets a collection by ID from the database
func (r *CollectionRepository) GetCollectionByID(ctx context.Context, id uint64) (*domain.Collection, error) {
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
		&collectionDAO.CreatedAt,
		&collectionDAO.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return converter.ToCollection(&collectionDAO), nil
}

// ListCollectionsByUserID lists all collections from the database
func (r *CollectionRepository) ListCollectionsByUserID(ctx context.Context, userID uint64, skip, limit uint64) ([]domain.Collection, error) {
	var collectionDAO dao.CollectionDAO
	var collections []domain.Collection

	query := r.db.QueryBuilder.Select("c.id, c.name, c.description, c.created_at, c.updated_at").
		From("collections c").
		Join("users_collections uc ON c.id = uc.collection_id").
		Where(sq.Eq{"uc.user_id": userID}).
		OrderBy("c.id").
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
			&collectionDAO.CreatedAt,
			&collectionDAO.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Convert the CollectionDAO to a domain.Collection
		collection := converter.ToCollection(&collectionDAO)

		// Append the converted collection to the list
		collections = append(collections, *collection)
	}

	return collections, nil
}

// UpdateCollection updates a collection by ID in the database
func (r *CollectionRepository) UpdateCollection(ctx context.Context, collection *domain.Collection) (*domain.Collection, error) {
	var collectionDAO dao.CollectionDAO

	name := nullString(collection.Name)
	description := nullString(collection.Description)

	query := r.db.QueryBuilder.Update("collections").
		Set("name", sq.Expr("COALESCE(?, name)", name)).
		Set("description", sq.Expr("COALESCE(?, description)", description)).
		Set("updated_at", time.Now()).
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
		&collectionDAO.CreatedAt,
		&collectionDAO.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return converter.ToCollection(&collectionDAO), nil
}

// DeleteCollection deletes a collection by ID from the database
func (r *CollectionRepository) DeleteCollection(ctx context.Context, id uint64) error {
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
func (r *CollectionRepository) IsUserPartOfCollection(ctx context.Context, userID, collectionID uint64) (bool, error) {
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
