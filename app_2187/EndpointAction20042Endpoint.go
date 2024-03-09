package main

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
	"utils"
)

// Action20042Endpoint is "get" request handler for "action_20042" endpoint.
func Action20042Endpoint(w http.ResponseWriter, r *http.Request) {

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

	var PropertyId6Variable int

	var paths19358 = mux.Vars(r)
	if paths19358 != nil && len(paths19358) > 0 {
		if val19358, ok := paths19358["property_id"]; ok {
			PropertyId6Variable, err = strconv.Atoi(val19358)
			if err != nil {
				errStatus = http.StatusBadRequest
				errResponse = utils.Response{Error: &utils.ValError{Message: "invalid property_id", Code: utils.ErrCodeBadInput}}
				return
			}
		}
	}

	var tx *sqlx.Tx
	tx, err = pg.Beginx()
	if err != nil {
		return
	}
	defer tx.Rollback()

	var PartnerListings1Variable []GetPropertyListingDataEPRConstruct

	err = tx.Select(&PartnerListings1Variable, `
        SELECT b.* FROM listings a LEFT JOIN partner_sites b ON a. partner_site_id = b.id LEFT JOIN properties c ON a.property_id = c.id WHERE a.property_id = $1`,
		PropertyId6Variable)

	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(PartnerListings1Variable))
}
