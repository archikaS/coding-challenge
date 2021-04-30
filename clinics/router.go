package clinics

import (
	"coding-challenge/middleware"

	"github.com/gorilla/mux"
)

func SetClinicRoutes(router *mux.Router) *mux.Router {

	router.HandleFunc("/clinics/get_dental_clinics",
		middleware.SetMiddlewareJSON(SearchDentalClinicController)).Methods("GET")
	router.HandleFunc("/clinics/get_vet_clinics",
		middleware.SetMiddlewareJSON(SearchVetClinicController)).Methods("GET")

	return router
}
