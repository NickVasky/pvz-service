package service

import (
	"AvitoTechPVZ/codegen/dto"
	"AvitoTechPVZ/repo"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

/*
type DefaultAPIServicer interface {
	DummyLoginPost(context.Context, DummyLoginPostRequest) (ImplResponse, error)
	RegisterPost(context.Context, RegisterPostRequest) (ImplResponse, error)
	LoginPost(context.Context, LoginPostRequest) (ImplResponse, error)
	PvzGet(context.Context, time.Time, time.Time, int32, int32) (ImplResponse, error)
	PvzPost(context.Context, Pvz) (ImplResponse, error)
	PvzPvzIdCloseLastReceptionPost(context.Context, string) (ImplResponse, error)
	PvzPvzIdDeleteLastProductPost(context.Context, string) (ImplResponse, error)
	ReceptionsPost(context.Context, ReceptionsPostRequest) (ImplResponse, error)
	ProductsPost(context.Context, ProductsPostRequest) (ImplResponse, error)
}
*/

type DefaultAPIServicerImpl struct {
	Repo repo.Repo
}

var allowedRoles = map[string]bool{
	"moderator": true,
	"client":    true,
}

type errMessage struct {
	Message string `json:"message"`
}

func (s *DefaultAPIServicerImpl) DummyLoginPost(ctx context.Context, r dto.DummyLoginPostRequest) (dto.ImplResponse, error) {
	// Implement your business logic here
	if _, ok := allowedRoles[r.Role]; !ok {
		msg := errMessage{Message: "Invalid role!"}
		return dto.Response(http.StatusBadRequest, msg), nil
	}

	userClaims := NewDummyUser(r.Role)
	token := userClaims.GenerateJwtToken()

	return dto.ImplResponse{
		Body: token,
		Code: http.StatusOK,
	}, nil
}

func (s *DefaultAPIServicerImpl) RegisterPost(ctx context.Context, r dto.RegisterPostRequest) (dto.ImplResponse, error) {
	return dto.ImplResponse{
		Body: "register",
		Code: http.StatusOK,
	}, nil
}

func (s *DefaultAPIServicerImpl) LoginPost(ctx context.Context, r dto.LoginPostRequest) (dto.ImplResponse, error) {
	return dto.ImplResponse{
		Body: "login",
		Code: http.StatusOK,
	}, nil
}

func (s *DefaultAPIServicerImpl) PvzGet(ctx context.Context, startDate time.Time, endDate time.Time, page int32, limit int32) (dto.ImplResponse, error) {
	return dto.ImplResponse{
		Body: "pvz (get)",
		Code: http.StatusOK,
	}, nil
}

func (s *DefaultAPIServicerImpl) PvzPost(ctx context.Context, pvz dto.Pvz) (dto.ImplResponse, error) {
	pvzId, err := uuid.Parse(pvz.Id)
	if err != nil {
		msg := errMessage{Message: "ID isn't a valid UUID!"}
		return dto.Response(http.StatusBadRequest, msg), nil
	}

	log.Println("Requested city: ", pvz.City)
	city, err := s.Repo.Cities.GetByName(pvz.City)
	if err != nil {
		log.Println(err)
		msg := errMessage{Message: "Invalid city!"}
		return dto.Response(http.StatusBadRequest, msg), nil
	}

	_, err = s.Repo.Pvzs.GetById(pvzId)
	if err == nil {
		log.Println(err)
		msg := errMessage{Message: "PVZ with such Id already exists!"}
		return dto.Response(http.StatusBadRequest, msg), nil
	}

	_, err = s.Repo.Pvzs.Add(pvzId, city.Id, pvz.RegistrationDate)

	if err != nil {
		msg := errMessage{Message: "Error!"}
		return dto.Response(http.StatusInternalServerError, msg), nil
	}

	newPvz, err := s.Repo.Pvzs.GetById(pvzId)

	if err != nil {
		msg := errMessage{Message: "Error!"}
		return dto.Response(http.StatusInternalServerError, msg), nil
	}

	return dto.ImplResponse{
		Body: newPvz,
		Code: http.StatusOK,
	}, nil
}

func (s *DefaultAPIServicerImpl) PvzPvzIdCloseLastReceptionPost(ctx context.Context, pvzIdParam string) (dto.ImplResponse, error) {
	return dto.ImplResponse{
		Body: "CloseLastReceptionPost",
		Code: http.StatusOK,
	}, nil
}

func (s *DefaultAPIServicerImpl) PvzPvzIdDeleteLastProductPost(ctx context.Context, pvzIdParam string) (dto.ImplResponse, error) {
	return dto.ImplResponse{
		Body: "DeleteLastProductPost",
		Code: http.StatusOK,
	}, nil
}

func (s *DefaultAPIServicerImpl) ReceptionsPost(ctx context.Context, r dto.ReceptionsPostRequest) (dto.ImplResponse, error) {
	pvzUUID, err := uuid.Parse(r.PvzId)
	if err != nil {
		msg := errMessage{Message: "ID isn't a valid UUID!"}
		return dto.Response(http.StatusBadRequest, msg), nil
	}

	_, err = s.Repo.Pvzs.GetById(pvzUUID)
	if err != nil {
		log.Println(err)
		msg := errMessage{Message: "Invalid city!"}
		return dto.Response(http.StatusBadRequest, msg), nil
	}

	rc, err := s.Repo.Receptions.GetOpened(pvzUUID)
	if err == nil && rc.PvzId == r.PvzId {
		log.Println(err)
		msg := errMessage{Message: "Reception already opened for PVZ!"}
		return dto.Response(http.StatusBadRequest, msg), nil
	}

	newID, err := s.Repo.Receptions.Add(pvzUUID)
	if err != nil {
		msg := errMessage{Message: "Error!"}
		return dto.Response(http.StatusInternalServerError, msg), nil
	}

	newRc, err := s.Repo.Receptions.GetById(newID)
	if err != nil {
		msg := errMessage{Message: "Error!"}
		return dto.Response(http.StatusInternalServerError, msg), nil
	}
	return dto.ImplResponse{
		Body: newRc,
		Code: http.StatusOK,
	}, nil
}

func (s *DefaultAPIServicerImpl) ProductsPost(ctx context.Context, r dto.ProductsPostRequest) (dto.ImplResponse, error) {
	return dto.ImplResponse{
		Body: "products",
		Code: http.StatusOK,
	}, nil
}
