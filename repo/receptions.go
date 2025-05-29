package repo

import (
	"AvitoTechPVZ/codegen/dto"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type ReceptionsRepo struct {
	DB iDB
}

type ReceptionStatus uint8

const (
	InProgress ReceptionStatus = 1
	Closed     ReceptionStatus = 2
)

var baseReceptionsSelector = psq.Select(
	"r.id",
	"r.date_time",
	"r.pvz_id",
	"s.name").
	From("receptions r").
	Join("statuses s ON r.status_id = s.id")

func (s ReceptionStatus) String() string {
	switch s {
	case InProgress:
		return "In Progress"
	case Closed:
		return "Closed"
	default:
		return "Unknown status"
	}
}

func (repo *ReceptionsRepo) Add(pvzId uuid.UUID) (uuid.UUID, error) {
	insertQuery := psq.
		Insert("receptions").
		Columns(
			"id",
			"pvz_id",
			"status_id",
			"date_time").
		Values(
			uuid.New(),
			pvzId.String(),
			InProgress,
			time.Now()).
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

func (repo *ReceptionsRepo) Close(receptionID uuid.UUID) error {
	sql, args, err := psq.
		Update("receptions").
		Set("status_id", Closed).
		Where(sq.Eq{"id": receptionID.String()}).
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

func (repo *ReceptionsRepo) GetById(ID uuid.UUID) (dto.Reception, error) {
	sql, args, err := baseReceptionsSelector.
		Where(sq.Eq{"r.id": ID}).
		ToSql()

	var r dto.Reception

	if err != nil {
		log.Println(err)
		return r, err
	}

	row := repo.DB.QueryRow(sql, args...)

	err = row.Scan(&r.Id, &r.DateTime, &r.PvzId, &r.Status)

	if err != nil {
		return r, err
	}
	log.Println("Pvz found: ", r)
	return r, nil
}

func (repo *ReceptionsRepo) GetByIds(IDs []uuid.UUID) ([]dto.Reception, error) {
	sql, args, err := baseReceptionsSelector.
		Where(sq.Eq{"r.id": IDs}).
		ToSql()

	var receptions []dto.Reception

	if err != nil {
		log.Println(err)
		return receptions, err
	}

	rows, err := repo.DB.Query(sql, args...)
	if err != nil {
		return receptions, err
	}
	defer rows.Close()

	for rows.Next() {
		var r dto.Reception
		if err := rows.Scan(&r.Id, &r.DateTime, &r.PvzId, &r.Status); err != nil {
			return receptions, err
		}
		receptions = append(receptions, r)
	}

	return receptions, nil
}

func (repo *ReceptionsRepo) GetOpened(pvzID uuid.UUID) (dto.Reception, error) {
	return repo.getByStatus(pvzID, InProgress)
}

func (repo *ReceptionsRepo) GetClosed(pvzID uuid.UUID) (dto.Reception, error) {
	return repo.getByStatus(pvzID, Closed)
}

func (repo *ReceptionsRepo) getByStatus(pvzID uuid.UUID, status ReceptionStatus) (dto.Reception, error) {
	sql, args, err := baseReceptionsSelector.
		Where(sq.And{
			sq.Eq{"r.pvz_id": pvzID.String()},
			sq.Eq{"s.id": status},
		}).
		ToSql()

	var r dto.Reception

	if err != nil {
		log.Println(err)
		return r, err
	}

	row := repo.DB.QueryRow(sql, args...)

	err = row.Scan(
		&r.Id,
		&r.DateTime,
		&r.PvzId,
		&r.Status)

	if err != nil {
		log.Println(err)
		return r, err
	}

	log.Println("Reception found: ", r)
	return r, nil
}
