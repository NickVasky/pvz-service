package repo

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"github.com/google/uuid"
)

type City struct {
	Id   uuid.UUID
	Name string
}

type CityRepo struct {
	DB *sql.DB
}

func (repo *CityRepo) GetByName(cityName string) (City, error) {
	sql, args, _ := sq.
		Select("*").
		From("cities").
		Where(sq.Eq{"name": cityName}).
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
		From("cities").
		Where(sq.Eq{"id": id}).
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
