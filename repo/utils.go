package repo

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// I know it's awful. TODO: read it from env at least
var (
	DbHost     = "localhost"
	DbPort     = 25913
	DbUsername = "avitopvz"
	DbPassword = "test"
	DbName     = "db_pvz"
)

type Repo struct {
	Users        UserRepo
	Roles        RoleRepo
	Cities       CityRepo
	Statuses     StatusRepo
	ProductTypes ProductTypeRepo
}

func NewRepo(db *sql.DB) Repo {
	return Repo{
		Users:        UserRepo{DB: db},
		Roles:        RoleRepo{DB: db},
		Cities:       CityRepo{DB: db},
		Statuses:     StatusRepo{DB: db},
		ProductTypes: ProductTypeRepo{DB: db},
	}
}

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
