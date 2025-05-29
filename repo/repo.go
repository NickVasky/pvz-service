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
	Products     ProductRepo
	Tx           txController
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{
		Users:        UserRepo{DB: db},
		Roles:        RoleRepo{DB: db},
		Cities:       CityRepo{DB: db},
		Statuses:     StatusRepo{DB: db},
		ProductTypes: ProductTypeRepo{DB: db},
		Pvzs:         PvzRepo{DB: db},
		Receptions:   ReceptionsRepo{DB: db},
		Products:     ProductRepo{DB: db},
		Tx:           txController{DB: db},
	}
}

type iDB interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

type txController struct {
	DB *sql.DB
}

func NewDbConn(cfg cfg.DbConfig) (*sql.DB, error) {
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DbName, cfg.DbSslMode)

	return sql.Open("postgres", connectionString)
}

func CloseDbConn(db *sql.DB) error {
	return db.Close()
}
