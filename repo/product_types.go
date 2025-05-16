package repo

import (
	"log"

	sq "github.com/Masterminds/squirrel"

	"github.com/google/uuid"
)

type ProductType struct {
	Id   int
	Uuid uuid.UUID
	Name string
}

type ProductTypeRepo struct {
	DB iDB
}

var productTypesCache = map[string]ProductType{}

func (repo *ProductTypeRepo) GetByName(productTypeName string) (ProductType, error) {
	var pt ProductType
	pt, ok := productTypesCache[productTypeName]
	if ok {
		log.Printf("Found product type \"%v\" in cache", pt)
		return pt, nil
	}

	sql, args, err := sq.
		Select("*").
		From("product_types").
		Where(sq.Eq{"name": productTypeName}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return pt, err
	}
	row := repo.DB.QueryRow(sql, args...)

	err = row.Scan(
		&pt.Id,
		&pt.Uuid,
		&pt.Name,
	)
	if err != nil {
		return pt, err
	}
	log.Printf("Found product type \"%v\" in DB, caching...", pt)
	productTypesCache[pt.Name] = pt
	return pt, err
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
		&productType.Uuid,
		&productType.Name,
	)
	return productType, err
}
