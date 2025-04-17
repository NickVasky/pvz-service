package repo

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var (
	DbHost     = "localhost"
	DbPort     = 25913
	DbUsername = "avitopvz"
	DbPassword = "test"
	DbName     = "db_pvz"
)

// TODO: To think about some common controller for all repos

func OpenDbConnection() *sql.DB {
	connectionString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		DbHost, DbPort, DbUsername, DbPassword, DbName)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	return db
}
