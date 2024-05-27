package postgresql

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func ConnectToDB(DNS string) (*sql.DB, error) {
	maxRetries := 10
	var db *sql.DB
	var err error
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", DNS)
		if err == nil {
			err = db.Ping()
			if err == nil {
				return db, nil
			}
		}
		log.Printf("Failed to connect to database, retrying... (%d/%d)", i+1, maxRetries)
		time.Sleep(2 * time.Second)
	}
	return nil, err
}
