package repo

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"github.com/google/uuid"
)

type Status struct {
	Id   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

type StatusRepo struct {
	DB *sql.DB
}

type statusesTableSchema struct {
	name, idCol, nameCol string
}

var statusesTable = statusesTableSchema{
	name:    "statuses",
	idCol:   "id",
	nameCol: "name",
}

func (repo *StatusRepo) GetByName(statusName string) (Status, error) {
	sql, args, _ := sq.
		Select("*").
		From(statusesTable.name).
		Where(sq.Eq{statusesTable.nameCol: statusName}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	row := repo.DB.QueryRow(sql, args...)

	status := Status{}
	err := row.Scan(
		&status.Id,
		&status.Name,
	)
	return status, err
}

func (repo *StatusRepo) GetById(id uuid.UUID) (Status, error) {
	sql, args, _ := sq.
		Select("*").
		From(statusesTable.name).
		Where(sq.Eq{statusesTable.idCol: id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	row := repo.DB.QueryRow(sql, args...)

	status := Status{}
	err := row.Scan(
		&status.Id,
		&status.Name,
	)
	return status, err
}
