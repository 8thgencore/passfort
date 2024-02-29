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
 * UserRepository implements postgres.UserRepository interface
 * and provides an access to the postgres database
 */
type UserRepository struct {
	db *database.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{
		db,
	}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(ctx context.Context, user *dao.UserDAO) (*dao.UserDAO, error) {
	var userDao dao.UserDAO

	query := r.db.QueryBuilder.Insert("users").
		Columns("name", "email", "password").
		Values(user.Name, user.Email, user.Password).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&userDao.ID,
		&userDao.Name,
		&userDao.Email,
		&userDao.Password,
		&userDao.IsVerified,
		&userDao.Role,
		&userDao.CreatedAt,
		&userDao.UpdatedAt,
	)
	if err != nil {
		if errCode := r.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return &userDao, nil
}

// GetUserByID gets a user by ID from the database
func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*dao.UserDAO, error) {
	var userDao dao.UserDAO

	query := r.db.QueryBuilder.Select("*").
		From("users").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&userDao.ID,
		&userDao.Name,
		&userDao.Email,
		&userDao.Password,
		&userDao.IsVerified,
		&userDao.Role,
		&userDao.CreatedAt,
		&userDao.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return &userDao, nil
}

// GetUserByEmailAndPassword gets a user by email from the database
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*dao.UserDAO, error) {
	var userDao dao.UserDAO

	query := r.db.QueryBuilder.Select("*").
		From("users").
		Where(sq.Eq{"email": email}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&userDao.ID,
		&userDao.Name,
		&userDao.Email,
		&userDao.Password,
		&userDao.IsVerified,
		&userDao.Role,
		&userDao.CreatedAt,
		&userDao.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return &userDao, nil
}

// ListUsers lists all users from the database
func (r *UserRepository) ListUsers(ctx context.Context, skip, limit uint64) ([]dao.UserDAO, error) {
	var userDao dao.UserDAO
	var usersDAO []dao.UserDAO

	query := r.db.QueryBuilder.Select("*").
		From("users").
		OrderBy("id").
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
			&userDao.ID,
			&userDao.Name,
			&userDao.Email,
			&userDao.Password,
			&userDao.IsVerified,
			&userDao.Role,
			&userDao.CreatedAt,
			&userDao.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Append the user to the list
		usersDAO = append(usersDAO, userDao)
	}

	return usersDAO, nil
}

// UpdateUser updates a user by ID in the database
func (r *UserRepository) UpdateUser(ctx context.Context, user *dao.UserDAO) (*dao.UserDAO, error) {
	var userDao dao.UserDAO

	name := nullString(user.Name)
	email := nullString(user.Email)
	password := nullString(user.Password)
	isVerified := nullBool(user.IsVerified)
	role := nullString(string(user.Role))

	query := r.db.QueryBuilder.Update("users").
		Set("name", sq.Expr("COALESCE(?, name)", name)).
		Set("email", sq.Expr("COALESCE(?, email)", email)).
		Set("password", sq.Expr("COALESCE(?, password)", password)).
		Set("is_verified", sq.Expr("COALESCE(?, is_verified)", isVerified)).
		Set("role", sq.Expr("COALESCE(?, role)", role)).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": user.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&userDao.ID,
		&userDao.Name,
		&userDao.Email,
		&userDao.Password,
		&userDao.IsVerified,
		&userDao.Role,
		&userDao.CreatedAt,
		&userDao.UpdatedAt,
	)
	if err != nil {
		if errCode := r.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return &userDao, nil
}

// DeleteUser deletes a user by ID from the database
func (r *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	query := r.db.QueryBuilder.Delete("users").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
