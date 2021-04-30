package clinics

import (
	"encoding/json"
	"net/http"
)

// ResponseData Struct is used to store response data
type ResponseData struct {
	StatusCode int         `json:"status_code"`
	Status     bool        `json:"status"`
	Result     interface{} `json:"result"`
}

// ResponseError is used to store error
type ResponseError struct {
	StatusCode int    `json:"status_code"`
	Status     bool   `json:"status"`
	Message    string `json:"message"`
}

func SearchDentalClinicController(w http.ResponseWriter, r *http.Request) {
	data, statusCode, err := SearchDentalClinics(r)
	if err != nil {
		w.WriteHeader(statusCode)
		errData := ResponseError{
			StatusCode: statusCode,
			Status:     false,
			Message:    err.Error(),
		}
		json.NewEncoder(w).Encode(errData)
	} else {
		resData := ResponseData{
			StatusCode: statusCode,
			Status:     true,
			Result:     data,
		}
		json.NewEncoder(w).Encode(resData)
	}
}

func SearchVetClinicController(w http.ResponseWriter, r *http.Request) {
	data, statusCode, err := SearchVetClinics(r)
	if err != nil {
		w.WriteHeader(statusCode)
		errData := ResponseError{
			StatusCode: statusCode,
			Status:     false,
			Message:    err.Error(),
		}
		json.NewEncoder(w).Encode(errData)
	} else {
		resData := ResponseData{
			StatusCode: statusCode,
			Status:     true,
			Result:     data,
		}
		json.NewEncoder(w).Encode(resData)
	}
}
