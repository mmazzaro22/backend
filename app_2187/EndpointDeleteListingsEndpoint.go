package main

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
	"utils"
)

// Delete listings
func DeleteListingsEndpoint(w http.ResponseWriter, r *http.Request) {

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

	var Id6Variable int

	var paths18110 = mux.Vars(r)
	if paths18110 != nil && len(paths18110) > 0 {
		if val18110, ok := paths18110["id"]; ok {
			Id6Variable, err = strconv.Atoi(val18110)
			if err != nil {
				errStatus = http.StatusBadRequest
				errResponse = utils.Response{Error: &utils.ValError{Message: "invalid id", Code: utils.ErrCodeBadInput}}
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

	var ListingToDelete1Variable ListingsModel

	err = tx.Get(&ListingToDelete1Variable, `
        SELECT id, guest_limit, property_id, order_type_id, creator_id, price, partner_site_id, check_in_time_id, check_out_time_id
        FROM listings
        WHERE (id = $1 OR $1 IS NULL)`,
		Id6Variable)

	if err != nil {
		return
	}

	_, err = tx.Exec(`
        DELETE FROM listings
        WHERE id = $1`,
		ListingToDelete1Variable.Id)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(nil))
}
