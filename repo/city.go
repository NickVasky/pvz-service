package repo

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"github.com/google/uuid"
)

type City struct {
	Id   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

type CityRepo struct {
	DB *sql.DB
}

type citiesTableSchema struct {
	name, idCol, nameCol string
}

var citiesTable = citiesTableSchema{
	name:    "cities",
	idCol:   "id",
	nameCol: "name",
}

func (repo *CityRepo) GetByName(cityName string) (City, error) {
	sql, args, _ := sq.
		Select("*").
		From(citiesTable.name).
		Where(sq.Eq{citiesTable.nameCol: cityName}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	row := repo.DB.QueryRow(sql, args...)

	city := City{}
	err := row.Scan(
		&city.Id,
		&city.Name,
	)
	return city, err
}

func (repo *CityRepo) GetById(id uuid.UUID) (City, error) {
	sql, args, _ := sq.
		Select("*").
		From(citiesTable.name).
		Where(sq.Eq{citiesTable.idCol: id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	row := repo.DB.QueryRow(sql, args...)

	city := City{}
	err := row.Scan(
		&city.Id,
		&city.Name,
	)
	return city, err
}
