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
		`SELECT u.id, u.name, u.email, u.password, u.is_verified, u.created_at, u.role_id, r.name FROM users u 
		JOIN roles r ON u.role_id = r.id
		WHERE u.email = $1
		`)
	if err != nil {
		return err
	}
	// If user doesn't exist return error.
	res := query.QueryRowContext(ctx, user.Email)
	// Update user pointer with CreatedAt field.
	if err := res.Scan(&user.Id, &user.Name, &user.Email, &user.Password.Hashed, &user.IsVerified, &user.CreatedAt, &user.RoleId, &user.Role.Name); err != nil {
		return err
	}
	return nil
}

func (s *UserStore) GetUserById(ctx context.Context, user *types.User) error {
	query, err := s.db.PrepareContext(ctx,
		`SELECT u.id, u.name, u.email, u.is_verified, u.created_at, u.role_id, r.name FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE u.id = $1
		`)
	if err != nil {
		return err
	}
	// If user doesn't exist return error.
	res := query.QueryRowContext(ctx, user.Id)
	// Update user pointer with CreatedAt field.
	if err := res.Scan(&user.Id, &user.Name, &user.Email, &user.IsVerified, &user.CreatedAt, &user.RoleId, &user.Role.Name); err != nil {
		return err
	}
	return nil
}

func (s *UserStore) UpdateUserRole(ctx context.Context, user *types.User, roleName string) error {
	var roleId int64
	query, err := s.db.PrepareContext(ctx,
		`SELECT id FROM roles WHERE name = $1`)
	if err != nil {
		return err
	}
	res := query.QueryRowContext(ctx, roleName)
	if err := res.Scan(&roleId); err != nil {
		return err
	}
	query, err = s.db.PrepareContext(ctx,
		`UPDATE users SET role_id = $1 WHERE id = $2`)
	if err != nil {
		return err
	}
	_, err = query.ExecContext(ctx, roleId, user.Id)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) UpdateUserPassword(ctx context.Context, user *types.User) error {
	query, err := s.db.PrepareContext(ctx, `UPDATE users SET password =$1 WHERE id =$2`)
	if err != nil {
		return err
	}

	_, err = query.ExecContext(ctx, user.Password.Hashed, user.Id)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) Verify(ctx context.Context, user *types.User) error {
	query, err := s.db.PrepareContext(ctx, `UPDATE users SET is_verified = TRUE WHERE id =$1`)
	if err != nil {
		return err
	}
	_, err = query.ExecContext(ctx, user.Id)
	if err != nil {
		return err
	}
	return nil
}
