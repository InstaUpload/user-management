package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/InstaUpload/user-management/types"
	_ "github.com/lib/pq"
)

func New(dbConfig *types.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbConfig.GetConnectionString())
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(dbConfig.MaxOpenConns)
	db.SetMaxIdleConns(dbConfig.MaxIdleConns)
	duration, err := time.ParseDuration(dbConfig.MaxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}
