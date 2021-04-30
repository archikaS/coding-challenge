package routers

import (
	clinicsService "coding-challenge/clinics"
	"coding-challenge/middleware"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {
	router := mux.NewRouter()
	//Default Router
	router.HandleFunc("/", middleware.SetMiddlewareJSON(defaultRouterHandler)).Methods("GET")

	// Clinics Router
	clinicsService.SetClinicRoutes(router)
	return router
}

func defaultRouterHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-type", "application/json")
	fmt.Println("https://github.com/archikaS/coding-challenge.git clicked")
	json.NewEncoder(w).Encode("Welcome To Coding API")
}
