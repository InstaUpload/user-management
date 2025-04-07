package store

import (
	"context"
	"database/sql"

	"github.com/InstaUpload/user-management/types"
)

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *types.User) error {
	// prepare query to insert user in database.
	query, err := s.db.PrepareContext(ctx,
		`INSERT INTO users (name, email, password) 
	VALUES ($1, $2, $3)
	RETURNING id, created_at`)
	if err != nil {
		return err
	}
	// If user doesn't exist Insert user in database.
	res := query.QueryRowContext(ctx, user.Name, user.Email, user.Password.Hashed)
	// Update user pointer with CreatedAt field.
	if err := res.Scan(&user.Id, &user.CreatedAt); err != nil {
		// TODO: Identify which error is thrown and convert it to custome error.
		return err
	}
	return nil
}

func (s *UserStore) GetUserByEmail(ctx context.Context, user *types.User) error {
	// prepare query to get user using email from database.
	query, err := s.db.PrepareContext(ctx,
		`SELECT id, name, email, password FROM users WHERE email = $1`)
	if err != nil {
		return err
	}
	// If user doesn't exist return error.
	res := query.QueryRowContext(ctx, user.Email)
	// Update user pointer with CreatedAt field.
	if err := res.Scan(&user.Id, &user.Name, &user.Email, &user.Password.Hashed); err != nil {
		return err
	}
	return nil
}

func (s *UserStore) GetUserById(ctx context.Context, user *types.User) error {
	query, err := s.db.PrepareContext(ctx,
		`SELECT id, name, email, is_verified, created_at FROM users WHERE id = $1`)
	if err != nil {
		return err
	}
	// If user doesn't exist return error.
	res := query.QueryRowContext(ctx, user.Id)
	// Update user pointer with CreatedAt field.
	if err := res.Scan(&user.Id, &user.Name, &user.Email, &user.IsVerified, &user.CreatedAt); err != nil {
		return err
	}
	return nil
}
