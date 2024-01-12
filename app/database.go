package app

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

func NewDB() *sql.DB {
	databaseUrl := os.Getenv("DATABASE_URL")
	db, err := sql.Open("mysql", databaseUrl)
	fmt.Println("DATABASE_URL", databaseUrl)
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}
