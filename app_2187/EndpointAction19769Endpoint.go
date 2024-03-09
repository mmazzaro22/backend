package main

import (
	"github.com/jmoiron/sqlx"
	"net/http"
	"utils"
)

// Action19769Endpoint is "get" request handler for "action_19769" endpoint.
func Action19769Endpoint(w http.ResponseWriter, r *http.Request) {

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

	var Listings1Variable []GetListingsEPRConstruct

	err = tx.Select(&Listings1Variable, `
        SELECT a.*, b.*, c.first_name, c.last_name, c.email FROM listings a LEFT JOIN properties b ON a.property_id = b.id LEFT JOIN users c ON b.creator_id = c.id`,
	)

	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(Listings1Variable))
}
