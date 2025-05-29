package repo

import (
	"AvitoTechPVZ/codegen/dto"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/google/uuid"
)

type ProductRepo struct {
	DB iDB
}

var baseProductsSelector = psq.
	Select(
		"p.id",
		"p.date_time",
		"t.name",
		"p.reception_id").
	From("products p").
	Join("product_types t ON t.id = p.type_id")

func (repo *ProductRepo) Add(productTypeId int, receptionId uuid.UUID) (uuid.UUID, error) {
	sql, args, err := psq.
		Insert("products").
		Columns(
			"id",
			"type_id",
			"reception_id",
			"date_time").
		Values(
			uuid.New(),
			productTypeId,
			receptionId,
			time.Now()).
		Suffix("RETURNING id").
		ToSql()

	var u uuid.UUID
	if err != nil {
		return u, err
	}

	var id string
	err = repo.DB.QueryRow(sql, args...).Scan(&id)
	if err != nil {
		return u, err
	}

	u, err = uuid.Parse(id)
	if err != nil {
		return u, err
	}
	return u, nil
}

func (repo *ProductRepo) Remove(id uuid.UUID) error {
	sql, args, err := psq.
		Delete("products").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return err
	}

	_, err = repo.DB.Exec(sql, args...)

	if err != nil {
		return err
	}
	return nil
}

func (repo *ProductRepo) GetByID(id uuid.UUID) (dto.Product, error) {
	sql, args, err := baseProductsSelector.
		Where(sq.Eq{"p.id": id.String()}).
		ToSql()

	var p dto.Product

	if err != nil {
		log.Println(err)
		return p, err
	}

	row := repo.DB.QueryRow(sql, args...)

	err = row.Scan(
		&p.Id,
		&p.DateTime,
		&p.Type,
		&p.ReceptionId)

	if err != nil {
		return p, err
	}

	log.Println("Product found: ", p)
	return p, err
}

func (repo *ProductRepo) GetLastByReception(receptionId uuid.UUID) (dto.Product, error) {
	sql, args, err := baseProductsSelector.
		Where(sq.Eq{"p.reception_id": receptionId.String()}).
		OrderBy("p.date_time DESC").
		Limit(1).
		ToSql()

	var p dto.Product

	if err != nil {
		log.Println(err)
		return p, err
	}

	row := repo.DB.QueryRow(sql, args...)

	err = row.Scan(
		&p.Id,
		&p.DateTime,
		&p.Type,
		&p.ReceptionId)

	if err != nil {
		log.Println(err)
		return p, err
	}

	log.Println("Product found: ", p)
	return p, err
}

func (repo *ProductRepo) GetPageByDate(startDate, endDate time.Time, limit, offset uint64) ([]dto.Product, error) {
	sql, args, err := baseProductsSelector.
		Where(
			sq.And{
				sq.GtOrEq{"p.date_time": startDate},
				sq.LtOrEq{"p.date_time": endDate}}).
		OrderBy("p.date_time DESC").
		Limit(limit).
		Offset(offset).
		ToSql()

	var products []dto.Product

	if err != nil {
		log.Println(err)
		return products, err
	}

	rows, err := repo.DB.Query(sql, args...)
	if err != nil {
		return products, err
	}
	defer rows.Close()

	for rows.Next() {
		var p dto.Product
		if err := rows.Scan(&p.Id, &p.DateTime, &p.Type, &p.ReceptionId); err != nil {
			return products, err
		}
		products = append(products, p)
	}

	return products, nil
}
