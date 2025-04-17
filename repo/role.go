package repo

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"github.com/google/uuid"
)

type Role struct {
	Id   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

type RoleRepo struct {
	DB *sql.DB
}

type rolesTableSchema struct {
	name, idCol, nameCol string
}

var rolesTable = rolesTableSchema{
	name:    "users",
	idCol:   "id",
	nameCol: "name",
}

func (repo RoleRepo) GetByName(roleName string) (Role, error) {
	sql, args, _ := sq.
		Select("*").
		From(rolesTable.name).
		Where(sq.Eq{rolesTable.nameCol: roleName}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	row := repo.DB.QueryRow(sql, args...)

	role := Role{}
	err := row.Scan(
		&role.Id,
		&role.Name,
	)
	return role, err
}

func (repo RoleRepo) GetById(id uuid.UUID) (Role, error) {
	sql, args, _ := sq.
		Select("*").
		From(rolesTable.name).
		Where(sq.Eq{rolesTable.idCol: id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	row := repo.DB.QueryRow(sql, args...)

	role := Role{}
	err := row.Scan(
		&role.Id,
		&role.Name,
	)
	return role, err
}
