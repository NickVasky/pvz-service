package repo

import (
	sq "github.com/Masterminds/squirrel"

	"github.com/google/uuid"
)

type Status struct {
	Id   int
	Uuid uuid.UUID
	Name string
}

type StatusRepo struct {
	DB iDB
}

func (repo *StatusRepo) GetByName(statusName string) (Status, error) {
	sql, args, _ := sq.
		Select("*").
		From("statuses").
		Where(sq.Eq{"name": statusName}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	row := repo.DB.QueryRow(sql, args...)

	status := Status{}
	err := row.Scan(
		&status.Id,
		&status.Uuid,
		&status.Name,
	)
	return status, err
}

func (repo *StatusRepo) GetById(id uuid.UUID) (Status, error) {
	sql, args, _ := sq.
		Select("*").
		From("statuses").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	row := repo.DB.QueryRow(sql, args...)

	status := Status{}
	err := row.Scan(
		&status.Id,
		&status.Uuid,
		&status.Name,
	)
	return status, err
}
