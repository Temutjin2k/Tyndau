package postgres

import (
	"context"
	"errors"

	"github.com/Temutjin2k/Tyndau/user_service/internal/adapter/postgres/dao"

	"github.com/Temutjin2k/Tyndau/user_service/internal/model"

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

func (r *UserRepo) Update(ctx context.Context, update *model.User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, avatar_link = $3, version = version + 1
		WHERE id = $4 AND version = $5 AND is_deleted = false
		RETURNING version
	`

	args := []any{
		update.Name,
		update.Email,
		update.AvatarLink,
		update.ID,
		update.Version,
	}

	err := r.db.QueryRow(ctx, query, args...).Scan(&update.Version)
	if err != nil {
		return err
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
			return model.User{}, ErrNotFound
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

func (r *UserRepo) Delete(ctx context.Context, id int64) error {
	query := `
		UPDATE users 
		SET is_deleted = TRUE 
		WHERE id = $1`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
