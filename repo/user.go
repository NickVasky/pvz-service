package repo

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type User struct {
	Id       uuid.UUID
	RoleId   uuid.UUID
	Email    string
	Password string
}

type UserRepo struct {
	DB *sql.DB
}

func (repo *UserRepo) All() ([]User, error) {
	sql, _, _ := psq.Select("*").From("users").ToSql()

	rows, err := repo.DB.Query(sql)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.Id,
			&user.RoleId,
			&user.Email,
			&user.Password)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *UserRepo) Add(u User) error {
	sql, args, _ := psq.
		Insert("users").
		Columns(
			"id",
			"role_id",
			"email",
			"password").
		Values(
			u.Id,
			u.RoleId,
			u.Email,
			u.Password).
		ToSql()

	_, err := repo.DB.Exec(sql, args...)

	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepo) GetById(id uuid.UUID) (User, error) {
	sql, args, _ := psq.
		Select("*").
		From("users").
		Where(sq.Eq{"id": id}).
		ToSql()

	row := repo.DB.QueryRow(sql, args...)

	user := User{}

	err := row.Scan(
		&user.Id,
		&user.RoleId,
		&user.Email,
		&user.Password)

	return user, err
}

func (repo *UserRepo) GetByEmail(email string) (User, error) {
	sql, args, _ := psq.
		Select("*").
		From("users").
		Where(sq.Eq{"email": email}).
		ToSql()

	row := repo.DB.QueryRow(sql, args...)

	user := User{}

	err := row.Scan(
		&user.Id,
		&user.RoleId,
		&user.Email,
		&user.Password)

	return user, err
}
