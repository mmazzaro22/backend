package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"utils"
)

// Action21102Endpoint is "get" request handler for "action_21102" endpoint.
func Action21102Endpoint(w http.ResponseWriter, r *http.Request) {

	var err error
	var errStatus int
	var errResponse utils.Response
	defer func(e *error) {
		if e != nil && *e != nil {
			log.Println(*e)
			if ee, ok := interface{}(*e).(*EventError); ok {
				errStatus = ee.StatusCode
				errResponse = utils.Response{Data: ee.Data, Error: &utils.ValError{Message: ee.Error(), Code: "event_error"}}
			}

			if errStatus == 0 {
				errResponse, errStatus = FormatErrorResponse(err)
			}
			utils.JSON(w, errStatus, errResponse)
		}
	}(&err)

	var AmenityName1Variable string

	var paths20387 = mux.Vars(r)
	if paths20387 != nil && len(paths20387) > 0 {
		if val20387, ok := paths20387["amenityName"]; ok {
			AmenityName1Variable = val20387
		}
	}

	var AmenityArray1Variable = []string{}

	AmenityArray1Variable = append(AmenityArray1Variable, AmenityName1Variable)

	utils.JSON(w, http.StatusOK, utils.OKResponse(AmenityArray1Variable))
}
