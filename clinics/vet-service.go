package clinics

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type vetClinicInfo struct {
	Name        string  `json:"clinicName"`
	State       string  `json:"stateCode"`
	Availablity timings `json:"opening"`
}

/*================================================================================================
			[SearchVetClinics] - Search Vet Clinics
	1) Fetch all clinics if no search condition is provided
	2) Fetch clinincs which satisfy search conditions
================================================================================================*/
func SearchVetClinics(r *http.Request) ([]vetClinicInfo, int, error) {
	var filteredClinicData = []vetClinicInfo{}

	queryParams := r.URL.Query()
	searchConditionKeys, searchOperator, onlyTimeConditionExists, err := validateQueryParams(queryParams)
	if err != nil {
		return nil, 400, err
	}

	vetClinicData, err := getVetClinicList()
	if err != nil {
		return nil, 500, err
	}

	// return all clinics if there is no search condition
	if len(r.URL.Query()) == 0 {
		return vetClinicData, 200, nil
	}

	//conditional functional call, based on search operator
	if searchOperator == "and" {
		filteredClinicData = searchVetClinicsBasedOnAndCondition(vetClinicData, searchConditionKeys, onlyTimeConditionExists)
	} else {
		filteredClinicData = searchVetClinicsBasedOnOrCondition(vetClinicData, searchConditionKeys)
	}
	return filteredClinicData, 200, nil
}

/* [getVetClinicList] - Get list of all vet clinics from a url.*/

func getVetClinicList() ([]vetClinicInfo, error) {
	request, err := http.NewRequest("GET", "https://storage.googleapis.com/scratchpay-code-challenge/vet-clinics.json", nil)
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	// send HTTP request using request object
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println(err)
		err := errors.New("There is some issue.")
		return nil, err
	}

	// read response body
	dataByte, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body.Close()

	// convert byte array into json format
	responseData := make([]vetClinicInfo, 0)

	err = json.Unmarshal(dataByte, &responseData)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}

/* [searchClinicsBasedOnAndCondition] - This function is used to filter out clinics based on
search conditions and search operator = AND.*/

func searchVetClinicsBasedOnAndCondition(clinicsData []vetClinicInfo, searchConditionKeys searchConditions, onlyTimeConditionExists bool) []vetClinicInfo {
	filteredData := make([]vetClinicInfo, 0)

	for i := range clinicsData {

		// Search condition check for clininc name
		isSearchConditionMatched := false

		clinicNameLowerCase := strings.ToLower(clinicsData[i].Name)
		if searchConditionKeys.clinicNameSearchPhase != "" {
			if searchConditionKeys.clinicNameSearchPhase == clinicNameLowerCase {
				isSearchConditionMatched = true
			} else {
				clinicNameSubStrings := strings.Split(clinicNameLowerCase, " ")
				for j := range clinicNameSubStrings {
					if searchConditionKeys.clinicNameSearchPhase == clinicNameSubStrings[j] {
						isSearchConditionMatched = true
					}
				}
			}
		}

		// Search condition check for state
		if searchConditionKeys.stateSearchPhase != "" {
			if isSearchConditionMatched || searchConditionKeys.clinicNameSearchPhase == "" {
				if searchConditionKeys.stateSearchPhase == strings.ToLower(clinicsData[i].State) {
					isSearchConditionMatched = true
				} else {
					isSearchConditionMatched = false
				}
			}
		}

		// Search condition check for opening and closing time

		// convert time availability data from string to time format
		availableFrom, _ := time.Parse("15:04", clinicsData[i].Availablity.From)
		availableTo, _ := time.Parse("15:04", clinicsData[i].Availablity.To)

		if onlyTimeConditionExists || isSearchConditionMatched {
			if searchConditionKeys.timeFromStr != "" && searchConditionKeys.timeToStr != "" {
				if availableFrom.After(searchConditionKeys.timeFrom) &&
					availableTo.Before(searchConditionKeys.timeTo) {
					isSearchConditionMatched = true
				} else {
					isSearchConditionMatched = false
				}
			} else if searchConditionKeys.timeFromStr != "" {
				if availableFrom.After(searchConditionKeys.timeFrom) {
					isSearchConditionMatched = true
				} else {
					isSearchConditionMatched = false
				}
			} else if searchConditionKeys.timeToStr != "" {
				if availableTo.Before(searchConditionKeys.timeTo) {
					isSearchConditionMatched = true
				} else {
					isSearchConditionMatched = false
				}
			}
		}

		// Search condition check for time

		if isSearchConditionMatched {
			filteredData = append(filteredData, clinicsData[i])
		}
	}
	return filteredData
}

/* [searchVetClinicsBasedOnOrCondition] - This function is used to filter out VET clinics based on
search conditions and search operator = OR.*/
func searchVetClinicsBasedOnOrCondition(clinicsData []vetClinicInfo, searchConditionKeys searchConditions) []vetClinicInfo {
	filteredData := make([]vetClinicInfo, 0)

	for i := range clinicsData {

		// Search condition check for clininc name
		isSearchConditionMatched := false

		clinicNameLowerCase := strings.ToLower(clinicsData[i].Name)
		if searchConditionKeys.clinicNameSearchPhase == clinicNameLowerCase {
			isSearchConditionMatched = true
		} else {
			clinicNameSubStrings := strings.Split(clinicNameLowerCase, " ")
			for j := range clinicNameSubStrings {
				if searchConditionKeys.clinicNameSearchPhase == clinicNameSubStrings[j] {
					isSearchConditionMatched = true
				}
			}
		}

		// Search condition check for state
		if searchConditionKeys.stateSearchPhase == strings.ToLower(clinicsData[i].State) {
			isSearchConditionMatched = true
		}

		// Search condition check for time

		// convert time availability data from string to time format
		availableFrom, _ := time.Parse("15:04", clinicsData[i].Availablity.From)
		availableTo, _ := time.Parse("15:04", clinicsData[i].Availablity.To)

		if searchConditionKeys.timeFromStr != "" && searchConditionKeys.timeToStr != "" {
			if availableFrom.After(searchConditionKeys.timeFrom) &&
				availableTo.Before(searchConditionKeys.timeTo) {
				isSearchConditionMatched = true
			}
		} else if searchConditionKeys.timeFromStr != "" {
			if availableFrom.After(searchConditionKeys.timeFrom) {
				isSearchConditionMatched = true
			}
		} else if searchConditionKeys.timeToStr != "" {
			if availableTo.Before(searchConditionKeys.timeTo) {
				isSearchConditionMatched = true
			}
		}

		if isSearchConditionMatched {
			filteredData = append(filteredData, clinicsData[i])
		}
	}
	return filteredData
}
