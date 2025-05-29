package service

import (
	"AvitoTechPVZ/codegen/dto"

	"github.com/google/uuid"
)

type UUIDable interface {
	ID() string
}

func getUniqueReceptions(items []dto.Product) []uuid.UUID {
	seen := make(map[string]struct{})
	result := []uuid.UUID{}

	for _, product := range items {
		if _, ok := seen[product.ReceptionId]; !ok {
			seen[product.ReceptionId] = struct{}{}
			uuid := uuid.MustParse(product.ReceptionId)
			result = append(result, uuid)
		}
	}

	return result
}

func getUniquePvzs(items []dto.Reception) []uuid.UUID {
	seen := make(map[string]struct{})
	result := []uuid.UUID{}

	for _, reception := range items {
		if _, ok := seen[reception.PvzId]; !ok {
			seen[reception.PvzId] = struct{}{}
			uuid := uuid.MustParse(reception.PvzId)
			result = append(result, uuid)
		}
	}

	return result
}

func buildPvzGetResponse(pvzs []dto.Pvz, receptions []dto.Reception, products []dto.Product) []dto.PvzGet200ResponseInner {
	var response = []dto.PvzGet200ResponseInner{}

	for _, pvz := range pvzs {
		response = append(response, dto.PvzGet200ResponseInner{
			Pvz:        pvz,
			Receptions: buildReception(pvz.Id, receptions, products),
		})
	}

	return response
}

func buildReception(pvzId string, receptions []dto.Reception, products []dto.Product) []dto.PvzGet200ResponseInnerReceptionsInner {
	var response = []dto.PvzGet200ResponseInnerReceptionsInner{}

	for _, reception := range receptions {
		if reception.PvzId == pvzId {
			response = append(response, dto.PvzGet200ResponseInnerReceptionsInner{
				Reception: reception,
				Products:  buildReceptionProducts(reception.Id, products),
			})
		}
	}

	return response
}

func buildReceptionProducts(receptionId string, products []dto.Product) []dto.Product {
	var response = []dto.Product{}

	for _, product := range products {
		if product.ReceptionId == receptionId {
			response = append(response, product)
		}
	}

	return response
}

// 3 builders seems to provide a better trade-offs in comparison to this:
/*
SELECT
	pr.id as product_id, pr.date_time as product_dateTime, pt.name as product_type, r.id as reception_id,
	r.date_time as reception_dateTime, s.name as reception_status, pv.id as pvz_id,
	pv.registration_date as pvz_registrationDate, c.name as pvz_city
FROM products AS pr
JOIN receptions AS r ON pr.reception_id = r.id
JOIN pvzs AS pv ON r.pvz_id = pv.id
JOIN cities AS c ON pv.city_id = c.id
JOIN product_types AS pt ON pr.type_id = pt.id
JOIN statuses AS s ON r.status_id = s.id
WHERE pr.date_time > '2025-05-14 13:05:00' AND pr.date_time < now()
ORDER BY pvz_id, reception_id, product_dateTime
*/
// Less time to build response json, but a lot of duplicate data to travel over network
// Also more processing on the DB side
