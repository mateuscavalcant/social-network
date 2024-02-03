package database

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func DB() *sql.DB {

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	_db := os.Getenv("DB")

	db, _ := sql.Open("mysql", user+":"+password+"@tcp(localhost:3306)/"+_db)
	err := db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}
