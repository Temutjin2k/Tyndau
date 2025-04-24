package postgres

import (
	"context"
	"errors"

	"user_service/internal/adapter/postgres/dao"
	"user_service/internal/model"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("user not found")

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

// Create(ctx context.Context, User model.User) (model.User, error)
// 	Update(ctx context.Context, update model.User) error
// 	GetProfile(ctx context.Context, email string) (model.User, error)
// 	GetByID(ctx context.Context, id int64) (model.User, error)

func (r *UserRepo) Create(ctx context.Context, user model.User) (model.User, error) {
	daoUser := dao.FromUser(user)

	query := `
		INSERT INTO users (name, email, avatar_link, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := r.db.QueryRow(ctx, query,
		daoUser.Name,
		daoUser.Email,
		daoUser.AvatarLink,
		daoUser.PasswordHash,
	).Scan(&daoUser.ID)

	if err != nil {
		return model.User{}, err
	}

	return dao.ToUser(daoUser), nil
}

func (r *UserRepo) Update(ctx context.Context, update model.User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, avatar_link = $3, version = version + 1
		WHERE id = $4 AND is_deleted = false
	`

	cmdTag, err := r.db.Exec(ctx, query,
		update.Name,
		update.Email,
		update.AvatarLink,
		update.ID,
	)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *UserRepo) GetProfile(ctx context.Context, email string) (model.User, error) {
	query := `
		SELECT id, created_at, name, email, avatar_link, password_hash, version
		FROM users
		WHERE email = $1 AND is_deleted = false
	`

	var daoUser dao.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&daoUser.ID,
		&daoUser.CreatedAt,
		&daoUser.Name,
		&daoUser.Email,
		&daoUser.AvatarLink,
		&daoUser.PasswordHash,
		&daoUser.Version,
	)
	if err != nil {
		switch {
		case err == pgx.ErrNoRows:
			return model.User{}, model.ErrDuplicateEmail
		default:
			return model.User{}, err
		}

	}

	return dao.ToUser(daoUser), nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (model.User, error) {
	query := `
		SELECT id, created_at, name, email, avatar_link, password_hash,version
		FROM users
		WHERE id = $1 AND is_deleted = false
	`

	var daoUser dao.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&daoUser.ID,
		&daoUser.CreatedAt,
		&daoUser.Name,
		&daoUser.Email,
		&daoUser.AvatarLink,
		&daoUser.PasswordHash,
		&daoUser.Version,
	)
	if err != nil {
		return model.User{}, err
	}

	return dao.ToUser(daoUser), nil
}
