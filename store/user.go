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
		`INSERT INTO users (name, email, password, role_id) 
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_on`)
	if err != nil {
		return err
	}
	// If user doesn't exist Insert user in database.
	res := query.QueryRowContext(ctx, user.Name, user.Email, user.Password.Hashed, user.RoleId)
	// Update user pointer with CreatedAt field.
	if err := res.Scan(&user.Id, &user.CreatedOn); err != nil {
		// TODO: Identify which error is thrown and convert it to custome error.
		return err
	}
	return nil
}

func (s *UserStore) GetUserByEmail(ctx context.Context, user *types.User) error {
	// prepare query to get user using email from database.
	query, err := s.db.PrepareContext(ctx,
		`SELECT u.id, u.name, u.email, u.password, u.is_verified, u.created_on, u.role_id, r.name FROM users u 
		JOIN roles r ON u.role_id = r.id
		WHERE u.email = $1
		`)
	if err != nil {
		return err
	}
	// If user doesn't exist return error.
	res := query.QueryRowContext(ctx, user.Email)
	// Update user pointer with CreatedAt field.
	if err := res.Scan(&user.Id, &user.Name, &user.Email, &user.Password.Hashed, &user.IsVerified, &user.CreatedOn, &user.RoleId, &user.Role.Name); err != nil {
		return err
	}
	return nil
}

func (s *UserStore) GetUserById(ctx context.Context, user *types.User) error {
	query, err := s.db.PrepareContext(ctx,
		`SELECT u.id, u.name, u.email, u.is_verified, u.created_on, u.role_id, r.name FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE u.id = $1
		`)
	if err != nil {
		return err
	}
	// If user doesn't exist return error.
	res := query.QueryRowContext(ctx, user.Id)
	// Update user pointer with CreatedAt field.
	if err := res.Scan(&user.Id, &user.Name, &user.Email, &user.IsVerified, &user.CreatedOn, &user.RoleId, &user.Role.Name); err != nil {
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

func (s *UserStore) AddEditorById(ctx context.Context, creatorUserId, userId int64) error {
	query, err := s.db.PrepareContext(ctx,
		`INSERT INTO editors (user_id, editor_id) VALUES ($1, $2)`)
	if err != nil {
		return err
	}
	_, err = query.ExecContext(ctx, creatorUserId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) GetEditorById(ctx context.Context, currentUserId int64, editor *types.Editor) error {
	query, err := s.db.PrepareContext(ctx,
		`SELECT user_id FROM editors WHERE user_id = $1 AND editor_id = $2`)
	if err != nil {
		return err
	}
	res := query.QueryRowContext(ctx, currentUserId, editor.Id)
	// Maybe need to add a editor struct in types package.
	if err := res.Scan(&editor.UserId); err != nil {
		return err
	}
	return nil
}
