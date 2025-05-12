package repo

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"github.com/google/uuid"
)

type Role struct {
	Id   int
	Uuid uuid.UUID
	Name string
}

type RoleRepo struct {
	DB *sql.DB
}

func (repo *RoleRepo) GetByName(roleName string) (Role, error) {
	sql, args, _ := sq.
		Select("*").
		From("users").
		Where(sq.Eq{"name": roleName}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	row := repo.DB.QueryRow(sql, args...)

	role := Role{}
	err := row.Scan(
		&role.Id,
		&role.Uuid,
		&role.Name,
	)
	return role, err
}

func (repo *RoleRepo) GetById(id uuid.UUID) (Role, error) {
	sql, args, _ := sq.
		Select("*").
		From("users").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	row := repo.DB.QueryRow(sql, args...)

	role := Role{}
	err := row.Scan(
		&role.Id,
		&role.Uuid,
		&role.Name,
	)
	return role, err
}
