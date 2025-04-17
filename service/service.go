package service

import (
	"AvitoTechPVZ/codegen/dto"
	"AvitoTechPVZ/repo"
	"context"
	"fmt"
	"net/http"
	"time"
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
	Users repo.UserRepo
}

func (s *DefaultAPIServicerImpl) DummyLoginPost(ctx context.Context, r dto.DummyLoginPostRequest) (dto.ImplResponse, error) {
	// Implement your business logic here

	// TODO: change handler to do actual thing
	users, err := s.Users.All()
	if err != nil {
		panic(err)
	}

	return dto.ImplResponse{
		Body: fmt.Sprint(users),
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
	return dto.ImplResponse{
		Body: "pvz (post)",
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
	return dto.ImplResponse{
		Body: "receptions",
		Code: http.StatusOK,
	}, nil
}

func (s *DefaultAPIServicerImpl) ProductsPost(ctx context.Context, r dto.ProductsPostRequest) (dto.ImplResponse, error) {
	return dto.ImplResponse{
		Body: "products",
		Code: http.StatusOK,
	}, nil
}
