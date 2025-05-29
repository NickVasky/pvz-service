package repo

import (
	"AvitoTechPVZ/codegen/dto"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type PvzRepo struct {
	DB iDB
}

var basePvzsSelector = psq.
	Select(
		"p.id",
		"c.name",
		"p.registration_date").
	From("pvzs p").
	Join("cities c ON p.city_id = c.id")

func (repo *PvzRepo) Add(pvzId uuid.UUID, cityID int, registrationDate time.Time) (uuid.UUID, error) {
	insertQuery := psq.
		Insert("pvzs").
		Columns(
			"id",
			"city_id",
			"registration_date").
		Values(
			pvzId.String(),
			cityID,
			registrationDate).
		Suffix("RETURNING id")

	var id string
	var u uuid.UUID
	err := insertQuery.RunWith(repo.DB).QueryRow().Scan(&id)

	if err != nil {
		log.Println(err)
		return u, err
	}

	u, err = uuid.Parse(id)
	if err == nil {
		log.Println(err)
		return u, err
	}

	return u, nil
}

func (repo *PvzRepo) GetById(id uuid.UUID) (dto.Pvz, error) {
	sql, args, err := basePvzsSelector.
		Where(sq.Eq{"p.id": id.String()}).
		ToSql()

	var p dto.Pvz

	if err != nil {
		log.Println(err)
		return p, err
	}

	row := repo.DB.QueryRow(sql, args...)

	err = row.Scan(
		&p.Id,
		&p.City,
		&p.RegistrationDate)

	log.Println("Pvz found: ", p)
	return p, err
}

func (repo *PvzRepo) GetByIds(IDs []uuid.UUID) ([]dto.Pvz, error) {
	sql, args, err := basePvzsSelector.
		Where(sq.Eq{"p.id": IDs}).
		ToSql()

	var pvzs []dto.Pvz

	if err != nil {
		log.Println(err)
		return pvzs, err
	}

	rows, err := repo.DB.Query(sql, args...)
	if err != nil {
		return pvzs, err
	}
	defer rows.Close()

	for rows.Next() {
		var p dto.Pvz
		if err := rows.Scan(&p.Id, &p.City, &p.RegistrationDate); err != nil {
			return pvzs, err
		}
		pvzs = append(pvzs, p)
	}

	return pvzs, nil
}
