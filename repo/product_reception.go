package repo

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var ErrNoOpenedReceptions = errors.New("no opened receptions")
var ErrNoProducts = errors.New("no products in reception")

func (t *txController) AddProductToReception(pvzId uuid.UUID, productTypeID int) (uuid.UUID, error) {
	fail := func(err error) (uuid.UUID, error) {
		return uuid.UUID{}, err
	}

	tx, err := t.DB.Begin()
	if err != nil {
		return fail(err)
	}
	defer tx.Rollback()

	pRepo := ProductRepo{DB: tx}
	rRepo := ReceptionsRepo{DB: tx}

	opened, err := rRepo.GetOpened(pvzId)
	if err != nil {
		if err == sql.ErrNoRows {
			return fail(ErrNoOpenedReceptions)
		}
		return fail(err)
	}

	receptionId, err := uuid.Parse(opened.Id)
	if err != nil {
		return fail(fmt.Errorf("reception ID isn't a valid UUID"))
	}

	newProductID, err := pRepo.Add(productTypeID, receptionId)
	if err != nil {
		return fail(err)
	}
	if err = tx.Commit(); err != nil {
		return fail(err)
	}
	return newProductID, nil
}

func (t *txController) DeleteLastAddedProduct(pvzId uuid.UUID) (uuid.UUID, error) {
	fail := func(err error) (uuid.UUID, error) {
		return uuid.UUID{}, err
	}

	tx, err := t.DB.Begin()
	if err != nil {
		return fail(err)
	}
	defer tx.Rollback()

	pRepo := ProductRepo{DB: tx}
	rRepo := ReceptionsRepo{DB: tx}

	opened, err := rRepo.GetOpened(pvzId)
	if err != nil {
		if err == sql.ErrNoRows {
			return fail(ErrNoOpenedReceptions)
		}
		return fail(err)
	}

	receptionId, err := uuid.Parse(opened.Id)
	if err != nil {
		return fail(fmt.Errorf("reception ID isn't a valid UUID"))
	}

	lastProduct, err := pRepo.GetLastByReception(receptionId)
	if err != nil {
		if err == sql.ErrNoRows {
			return fail(ErrNoProducts)
		}
		return fail(err)
	}

	lastProductID, err := uuid.Parse(lastProduct.Id)
	if err != nil {
		return fail(fmt.Errorf("reception ID isn't a valid UUID"))
	}

	err = pRepo.Remove(lastProductID)
	if err != nil {
		return fail(err)
	}

	if err = tx.Commit(); err != nil {
		return fail(err)
	}
	return lastProductID, nil
}
