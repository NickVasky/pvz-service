package repo

import (
	cfg "AvitoTechPVZ/config"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
)

var psq = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type Repo struct {
	Users        UserRepo
	Roles        RoleRepo
	Cities       CityRepo
	Statuses     StatusRepo
	ProductTypes ProductTypeRepo
	Pvzs         PvzRepo
	Receptions   ReceptionsRepo
}

func NewRepo(db *sql.DB) Repo {
	return Repo{
		Users:        UserRepo{DB: db},
		Roles:        RoleRepo{DB: db},
		Cities:       CityRepo{DB: db},
		Statuses:     StatusRepo{DB: db},
		ProductTypes: ProductTypeRepo{DB: db},
		Pvzs:         PvzRepo{DB: db},
		Receptions:   ReceptionsRepo{DB: db},
	}
}

func OpenDbConnection(cfg cfg.DbConfig) *sql.DB {
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DbName, cfg.DbSslMode)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	return db
}
