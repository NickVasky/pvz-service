package service

import (
	"AvitoTechPVZ/codegen/dto"
	"AvitoTechPVZ/repo"
	"AvitoTechPVZ/security"
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"runtime"
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

type defaultAPIServicerImpl struct {
	Repo               *repo.Repo
	SecurityController *security.SecurityController
}

var allowedRoles = map[string]bool{
	"moderator": true,
	"client":    true,
}

type responseErr struct {
	Message string `json:"message"`
	code    int
	err     error
}

func NewDefaultAPIServicerImpl(r *repo.Repo, s *security.SecurityController) *defaultAPIServicerImpl {
	return &defaultAPIServicerImpl{
		Repo:               r,
		SecurityController: s,
	}
}

func (e responseErr) handle() (dto.ImplResponse, error) {
	pc, _, _, ok := runtime.Caller(1)
	funcName := "unknown"
	if ok {
		funcName = runtime.FuncForPC(pc).Name()
	}
	log.Printf("[%s] - Error: %v, Client Error message: %s", funcName, e.err, e.Message)

	return dto.Response(e.code, e), nil
}

func (s *defaultAPIServicerImpl) DummyLoginPost(ctx context.Context, r dto.DummyLoginPostRequest) (dto.ImplResponse, error) {
	// Implement your business logic here
	if _, ok := allowedRoles[r.Role]; !ok {
		return responseErr{
			Message: "Invalid role!",
			code:    http.StatusBadRequest,
		}.handle()
	}

	userClaims := security.NewDummyUser(r.Role)
	token := s.SecurityController.GenerateJwtToken(userClaims)

	return dto.ImplResponse{
		Body: token,
		Code: http.StatusOK,
	}, nil
}

func (s *defaultAPIServicerImpl) RegisterPost(ctx context.Context, r dto.RegisterPostRequest) (dto.ImplResponse, error) {
	return dto.ImplResponse{
		Body: "register",
		Code: http.StatusOK,
	}, nil
}

func (s *defaultAPIServicerImpl) LoginPost(ctx context.Context, r dto.LoginPostRequest) (dto.ImplResponse, error) {
	return dto.ImplResponse{
		Body: "login",
		Code: http.StatusOK,
	}, nil
}

func (s *defaultAPIServicerImpl) PvzGet(ctx context.Context, startDate time.Time, endDate time.Time, page int32, limit int32) (dto.ImplResponse, error) {
	if page <= 0 {
		return responseErr{
			Message: "Page number should be > 0!",
			code:    http.StatusUnprocessableEntity,
			err:     nil,
		}.handle()
	}
	if limit <= 0 {
		return responseErr{
			Message: "Limit number should be > 0!",
			code:    http.StatusUnprocessableEntity,
			err:     nil,
		}.handle()
	}

	products, err := s.Repo.Products.GetPageByDate(startDate, endDate, uint64(limit), uint64((page-1)*limit))
	if err != nil {
		return responseErr{
			Message: "Error!",
			code:    http.StatusInternalServerError,
			err:     err,
		}.handle()
	}

	receptionIDs := getUniqueReceptions(products)
	receptions, err := s.Repo.Receptions.GetByIds(receptionIDs)
	if err != nil {
		return responseErr{
			Message: "Error!",
			code:    http.StatusInternalServerError,
			err:     err,
		}.handle()
	}

	pvzIDs := getUniquePvzs(receptions)
	pvzs, err := s.Repo.Pvzs.GetByIds(pvzIDs)
	if err != nil {
		return responseErr{
			Message: "Error!",
			code:    http.StatusInternalServerError,
			err:     err,
		}.handle()
	}
	response := buildPvzGetResponse(pvzs, receptions, products)

	return dto.ImplResponse{
		Body: response,
		Code: http.StatusOK,
	}, nil
}

func (s *defaultAPIServicerImpl) PvzPost(ctx context.Context, pvz dto.Pvz) (dto.ImplResponse, error) {
	pvzId, err := uuid.Parse(pvz.Id)
	if err != nil {
		return responseErr{
			Message: "ID isn't a valid UUID!",
			code:    http.StatusBadRequest,
			err:     err,
		}.handle()
	}

	log.Println("Requested city: ", pvz.City)
	city, err := s.Repo.Cities.GetByName(pvz.City)
	if err != nil {
		return responseErr{
			Message: "Invalid city!",
			code:    http.StatusBadRequest,
			err:     err,
		}.handle()
	}

	_, err = s.Repo.Pvzs.GetById(pvzId)
	if err == nil {
		return responseErr{
			Message: "PVZ with such Id already exists!",
			code:    http.StatusBadRequest,
			err:     err,
		}.handle()
	}

	_, err = s.Repo.Pvzs.Add(pvzId, city.Id, pvz.RegistrationDate)

	if err != nil {
		return responseErr{
			Message: "Error!",
			code:    http.StatusInternalServerError,
			err:     err,
		}.handle()
	}

	newPvz, err := s.Repo.Pvzs.GetById(pvzId)

	if err != nil {
		return responseErr{
			Message: "Error!",
			code:    http.StatusInternalServerError,
			err:     err,
		}.handle()
	}

	return dto.ImplResponse{
		Body: newPvz,
		Code: http.StatusOK,
	}, nil
}

func (s *defaultAPIServicerImpl) PvzPvzIdCloseLastReceptionPost(ctx context.Context, pvzIdParam string) (dto.ImplResponse, error) {
	pvzId, err := uuid.Parse(pvzIdParam)
	if err != nil {
		return responseErr{
			Message: "ID isn't a valid UUID!",
			code:    http.StatusBadRequest,
			err:     err,
		}.handle()
	}

	reception, err := s.Repo.Receptions.GetOpened(pvzId)
	if err != nil {
		if err == sql.ErrNoRows {
			return responseErr{
				Message: "No opened receptions",
				code:    http.StatusBadRequest,
				err:     err,
			}.handle()
		}
		return responseErr{
			Message: "Error!",
			code:    http.StatusInternalServerError,
			err:     err,
		}.handle()
	}
	receptionId := uuid.MustParse(reception.Id)
	_, err = s.Repo.Products.GetLastByReception(receptionId)
	if err != nil {
		if err == sql.ErrNoRows {
			return responseErr{
				Message: "No products in reception yet",
				code:    http.StatusBadRequest,
				err:     err,
			}.handle()
		}
		return responseErr{
			Message: "Error!",
			code:    http.StatusInternalServerError,
			err:     err,
		}.handle()
	}
	err = s.Repo.Receptions.Close(receptionId)
	if err != nil {
		return responseErr{
			Message: "Error!",
			code:    http.StatusInternalServerError,
			err:     err,
		}.handle()
	}

	reception, err = s.Repo.Receptions.GetById(receptionId)
	if err != nil {
		return responseErr{
			Message: "Error!",
			code:    http.StatusInternalServerError,
			err:     err,
		}.handle()
	}

	return dto.ImplResponse{
		Body: reception,
		Code: http.StatusOK,
	}, nil
}

func (s *defaultAPIServicerImpl) PvzPvzIdDeleteLastProductPost(ctx context.Context, pvzIdParam string) (dto.ImplResponse, error) {
	pvzUUID, err := uuid.Parse(pvzIdParam)
	if err != nil {
		return responseErr{
			Message: "ID isn't a valid UUID!",
			code:    http.StatusBadRequest,
			err:     err,
		}.handle()
	}
	_, err = s.Repo.Tx.DeleteLastAddedProduct(pvzUUID)
	if err != nil {
		if errors.Is(err, repo.ErrNoProducts) {
			return responseErr{
				Message: err.Error(),
				code:    http.StatusBadRequest,
				err:     err,
			}.handle()
		}
		return responseErr{
			Message: "Error!",
			code:    http.StatusBadRequest,
			err:     err,
		}.handle()
	}

	return dto.ImplResponse{
		Body: nil,
		Code: http.StatusOK,
	}, nil
}

func (s *defaultAPIServicerImpl) ReceptionsPost(ctx context.Context, r dto.ReceptionsPostRequest) (dto.ImplResponse, error) {
	pvzUUID, err := uuid.Parse(r.PvzId)
	if err != nil {
		return responseErr{
			Message: "ID isn't a valid UUID!",
			code:    http.StatusBadRequest,
			err:     err,
		}.handle()
	}

	_, err = s.Repo.Pvzs.GetById(pvzUUID)
	if err != nil {
		return responseErr{
			Message: "Invalid city!",
			code:    http.StatusBadRequest,
			err:     err,
		}.handle()
	}

	rc, err := s.Repo.Receptions.GetOpened(pvzUUID)
	if err == nil && rc.PvzId == r.PvzId {
		return responseErr{
			Message: "Reception already opened for PVZ!",
			code:    http.StatusBadRequest,
			err:     err,
		}.handle()
	}

	newID, err := s.Repo.Receptions.Add(pvzUUID)
	if err != nil {
		return responseErr{
			Message: "Error!",
			code:    http.StatusInternalServerError,
			err:     err,
		}.handle()
	}

	newRc, err := s.Repo.Receptions.GetById(newID)
	if err != nil {
		return responseErr{
			Message: "Error!",
			code:    http.StatusInternalServerError,
			err:     err,
		}.handle()
	}
	return dto.ImplResponse{
		Body: newRc,
		Code: http.StatusOK,
	}, nil
}

func (s *defaultAPIServicerImpl) ProductsPost(ctx context.Context, r dto.ProductsPostRequest) (dto.ImplResponse, error) {
	pvzUUID, err := uuid.Parse(r.PvzId)
	if err != nil {
		return responseErr{
			Message: "ID isn't a valid UUID!",
			code:    http.StatusBadRequest,
			err:     err,
		}.handle()
	}

	productTypeId, err := s.Repo.ProductTypes.GetByName(r.Type)
	if err != nil {
		return responseErr{
			Message: "Invalid product type!",
			code:    http.StatusBadRequest,
			err:     err,
		}.handle()
	}

	newProductId, err := s.Repo.Tx.AddProductToReception(pvzUUID, productTypeId.Id)
	if err != nil {
		if errors.Is(err, repo.ErrNoOpenedReceptions) {
			return responseErr{
				Message: err.Error(), //think this errors through
				code:    http.StatusBadRequest,
				err:     err,
			}.handle()
		}
		return responseErr{
			Message: "Error!", //think this errors through
			code:    http.StatusInternalServerError,
			err:     err,
		}.handle()
	}

	product, err := s.Repo.Products.GetByID(newProductId)
	if err != nil {
		return responseErr{
			Message: "Error!",
			code:    http.StatusInternalServerError,
			err:     err,
		}.handle()
	}

	return dto.ImplResponse{
		Body: product,
		Code: http.StatusOK,
	}, nil
}
