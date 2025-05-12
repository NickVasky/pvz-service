package repo

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"github.com/google/uuid"
)

type ProductType struct {
	Id   uuid.UUID
	Name string
}

type ProductTypeRepo struct {
	DB *sql.DB
}

func (repo *ProductTypeRepo) GetByName(productTypeName string) (ProductType, error) {
	sql, args, _ := sq.
		Select("*").
		From("product_types").
		Where(sq.Eq{"name": productTypeName}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	row := repo.DB.QueryRow(sql, args...)

	productType := ProductType{}
	err := row.Scan(
		&productType.Id,
		&productType.Name,
	)
	return productType, err
}

func (repo *ProductTypeRepo) GetById(id uuid.UUID) (ProductType, error) {
	sql, args, _ := sq.
		Select("*").
		From("product_types").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	row := repo.DB.QueryRow(sql, args...)

	productType := ProductType{}
	err := row.Scan(
		&productType.Id,
		&productType.Name,
	)
	return productType, err
}
