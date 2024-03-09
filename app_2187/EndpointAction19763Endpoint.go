package main

import (
	"github.com/jmoiron/sqlx"
	"net/http"
	"utils"
)

// Action19763Endpoint is "get" request handler for "action_19763" endpoint.
func Action19763Endpoint(w http.ResponseWriter, r *http.Request) {

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

	var tx *sqlx.Tx
	tx, err = pg.Beginx()
	if err != nil {
		return
	}
	defer tx.Rollback()

	var Properties2Variable []PropertiesModel

	err = tx.Select(&Properties2Variable, `
        SELECT id, name, creator_id, description, address, city, state, zip, images, land_type, amenity, status, acres
        FROM properties`,
	)

	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(Properties2Variable))
}
