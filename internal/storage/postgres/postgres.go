package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Storage struct {
	Db *sql.DB
}

func New(dbUrl string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{Db: db}, nil
}
