package clinics

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type dentalClinicInfo struct {
	Name        string  `json:"name"`
	State       string  `json:"stateName"`
	Availablity timings `json:"availability"`
}

type timings struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type searchConditions struct {
	clinicNameSearchPhase string
	stateSearchPhase      string
	timeFromStr           string
	timeToStr             string
	timeFrom              time.Time
	timeTo                time.Time
}

/*================================================================================================
			[SearchDentalClinics] - Search Dental Clinics
	1) Fetch all clinics if no search condition is provided
	2) Fetch clinincs which satisfy search conditions
================================================================================================*/
func SearchDentalClinics(r *http.Request) ([]dentalClinicInfo, int, error) {
	var filteredClinicData = []dentalClinicInfo{}

	queryParams := r.URL.Query()
	searchConditionKeys, searchOperator, onlyTimeConditionExists, err := validateQueryParams(queryParams)
	if err != nil {
		return nil, 400, err
	}

	dentalClinicData, err := getDentalClinicList()
	if err != nil {
		return nil, 500, err
	}

	// return all clinics if there is no search condition
	if len(r.URL.Query()) == 0 {
		return dentalClinicData, 200, nil
	}

	//conditional functional call, based on search operator
	if searchOperator == "and" {
		filteredClinicData = searchClinicsBasedOnAndCondition(dentalClinicData, searchConditionKeys, onlyTimeConditionExists)
	} else {
		filteredClinicData = searchClinicsBasedOnOrCondition(dentalClinicData, searchConditionKeys)
	}
	return filteredClinicData, 200, nil
}

/* [validateQueryParams] -  It is a Common Function called from both dental-service and
vet-service to validate query params.*/

func validateQueryParams(queryParams url.Values) (searchConditions, string, bool, error) {
	var searchConditionKeys = searchConditions{}
	searchOperator := "or"
	onlyTimeConditionExists := true

	//Query params for data to be searched
	if keys, ok := queryParams["clinicName"]; ok {
		if len(keys[0]) >= 1 {
			onlyTimeConditionExists = false
			searchConditionKeys.clinicNameSearchPhase = strings.ToLower(keys[0])
		} else {
			err := errors.New("Please provide clinic name for search.")
			return searchConditions{}, "", false, err
		}
	}

	if keys, ok := queryParams["state"]; ok {
		if len(keys[0]) >= 1 {
			onlyTimeConditionExists = false
			searchConditionKeys.stateSearchPhase = strings.ToLower(keys[0])
		} else {
			err := errors.New("Please provide state for search.")
			return searchConditions{}, "", false, err
		}
	}

	if keys, ok := queryParams["openFrom"]; ok {
		if len(keys[0]) >= 1 {
			var err error
			searchConditionKeys.timeFromStr = keys[0]
			// convert string to time format
			searchConditionKeys.timeFrom, err = time.Parse("15:04", keys[0])
			if err != nil {
				fmt.Println(err)
				err := errors.New("Please provide time in hour and minute format.")
				return searchConditions{}, "", false, err
			}
			//Subtract 1 sec
			searchConditionKeys.timeFrom = searchConditionKeys.timeFrom.Add(-time.Second * 1)
		} else {
			err := errors.New("Please provide open from for search.")
			return searchConditions{}, "", false, err
		}
	}

	if keys, ok := queryParams["openTo"]; ok {
		if len(keys[0]) >= 1 {
			var err error
			searchConditionKeys.timeToStr = keys[0]
			// convert string to time format
			searchConditionKeys.timeTo, err = time.Parse("15:04", keys[0])
			if err != nil {
				err := errors.New("Please provide time in hour and minute format.")
				return searchConditions{}, "", false, err
			}
			//Add 1 sec
			searchConditionKeys.timeTo = searchConditionKeys.timeTo.Add(time.Second * 1)
		} else {
			err := errors.New("Please provide open to for search.")
			return searchConditions{}, "", false, err
		}
	}

	if keys, ok := queryParams["condition"]; ok {
		if (searchConditions{} != searchConditionKeys) { // check if the stuct is empty
			if len(keys[0]) > 0 {
				if len(queryParams) == 1 {
					searchOperator = "or"
				} else if len(queryParams) >= 2 {
					if keys[0] == "or" || keys[0] == "and" {
						searchOperator = keys[0] // this condition can either be or/and
					} else {
						err := errors.New("Please provide a valid value for condition.")
						return searchConditions{}, "", false, err
					}
				}
			} else {
				err := errors.New("Please provide condition for search.")
				return searchConditions{}, "", false, err
			}
		} else {
			err := errors.New("Please provide atleast one key for search.")
			return searchConditions{}, "", false, err
		}
	}

	return searchConditionKeys, searchOperator, onlyTimeConditionExists, nil
}

/* [getDentalClinicList] - Get list of all dental clinics from a url.*/

func getDentalClinicList() ([]dentalClinicInfo, error) {
	request, err := http.NewRequest("GET", "https://storage.googleapis.com/scratchpay-code-challenge/dental-clinics.json", nil)
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
	responseData := make([]dentalClinicInfo, 0)
	err = json.Unmarshal(dataByte, &responseData)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}

/* [searchClinicsBasedOnAndCondition] - This function is used to filter out clinics based on
search conditions and search operator = AND.*/

func searchClinicsBasedOnAndCondition(clinicsData []dentalClinicInfo, searchConditionKeys searchConditions, onlyTimeConditionExists bool) []dentalClinicInfo {
	filteredData := make([]dentalClinicInfo, 0)

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

/* [searchClinicsBasedOnOrCondition] - This function is used to filter out clinics based on
search conditions and search operator = OR.*/

func searchClinicsBasedOnOrCondition(clinicsData []dentalClinicInfo, searchConditionKeys searchConditions) []dentalClinicInfo {
	filteredData := make([]dentalClinicInfo, 0)

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
