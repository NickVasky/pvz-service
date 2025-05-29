package service

import (
	"AvitoTechPVZ/codegen/dto"
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter creates a new router for any number of api routers
func NewRouter(routers ...dto.Router) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, api := range routers {
		for name, route := range api.Routes() {
			var handler http.Handler = route.HandlerFunc
			if endpointAccess[name].isProtected {
				handler = AuthMiddleware(handler, endpointAccess[name].roles)
			}
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
