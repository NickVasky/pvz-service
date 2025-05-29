package service

import (
	"AvitoTechPVZ/codegen/dto"
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter creates a new router from service
// func NewRouter(s security.SecurityController, routers ...dto.Router) *mux.Router {
// 	router := mux.NewRouter().StrictSlash(true)
// 	for _, api := range routers {
// 		for name, route := range api.Routes() {
// 			var handler http.Handler = route.HandlerFunc

// 			handler = s.AuthMiddleware(handler, name)
// 			handler = dto.Logger(handler, name)

// 			router.
// 				Methods(route.Method).
// 				Path(route.Pattern).
// 				Name(name).
// 				Handler(handler)
// 		}
// 	}

// 	return router
// }

func (s *defaultAPIServicerImpl) NewRouter(routers ...dto.Router) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, api := range routers {
		for name, route := range api.Routes() {
			var handler http.Handler = route.HandlerFunc

			handler = s.SecurityController.AuthMiddleware(handler, name)
			handler = dto.Logger(handler, name)

			router.
				Methods(route.Method).
				Path(route.Pattern).
				Name(name).
				Handler(handler)
		}
	}

	return router
}
