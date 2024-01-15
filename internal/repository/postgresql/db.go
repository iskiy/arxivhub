package postgresql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func ConnectToDB(DNS string) (*sql.DB, error) {
	db, err := sql.Open("postgres", DNS)
	if err != nil {
		return nil, fmt.Errorf("connection to repository error: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("check conn to repository error : %w", err)
	}

	return db, nil
}
