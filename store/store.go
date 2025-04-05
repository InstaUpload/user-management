package store

import (
	"context"
	"database/sql"

	"github.com/InstaUpload/user-management/types"
)

type Store struct {
	User interface {
		Create(context.Context, *types.User) error
		GetUserByEmail(context.Context, *types.User) error
	}
}

var MockStore Store

func NewStore(db *sql.DB) Store {
	return Store{
		User: &UserStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err = fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
