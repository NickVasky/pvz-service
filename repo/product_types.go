package repo

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"github.com/google/uuid"
)

type ProductType struct {
	Id   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

type ProductTypeRepo struct {
	DB *sql.DB
}

type productTypesTableSchema struct {
	name, idCol, nameCol string
}

var productTypesTable = productTypesTableSchema{
	name:    "product_types",
	idCol:   "id",
	nameCol: "name",
}

func (repo *ProductTypeRepo) GetByName(productTypeName string) (ProductType, error) {
	sql, args, _ := sq.
		Select("*").
		From(productTypesTable.name).
		Where(sq.Eq{productTypesTable.nameCol: productTypeName}).
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
		From(productTypesTable.name).
		Where(sq.Eq{productTypesTable.idCol: id}).
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
