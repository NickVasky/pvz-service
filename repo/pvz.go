package repo

import (
	"AvitoTechPVZ/codegen/dto"
	"database/sql"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

// type Pvz struct {
// 	Id               uuid.UUID
// 	City             string
// 	RegistrationDate time.Time
// }

type PvzRepo struct {
	DB *sql.DB
}

type pvzsTableSchema struct {
	id, cityId, registrationDate string
}

var tblName = "pvzs"
var tbl = pvzsTableSchema{
	id:               "id",
	cityId:           "city_id",
	registrationDate: "registration_date",
}

func (repo *PvzRepo) Add(pvzId, cityID uuid.UUID, registrationDate time.Time) (uuid.UUID, error) {
	insertQuery := psq.
		Insert(tblName).
		Columns(
			tbl.id,
			tbl.cityId,
			tbl.registrationDate).
		Values(
			pvzId.String(),
			cityID.String(),
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
	sql, args, err := psq.
		Select(
			"p.id",
			"c.name",
			"p.registration_date").
		From(tblName + " p").
		Join("cities c ON p.city_id = c.id").
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
