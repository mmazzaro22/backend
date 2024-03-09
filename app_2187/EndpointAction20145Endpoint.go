package main

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
	"utils"
)

// Action20145Endpoint is "put" request handler for "action_20145" endpoint.
func Action20145Endpoint(w http.ResponseWriter, r *http.Request) {

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

	var Id14Variable int

	var paths19457 = mux.Vars(r)
	if paths19457 != nil && len(paths19457) > 0 {
		if val19457, ok := paths19457["id"]; ok {
			Id14Variable, err = strconv.Atoi(val19457)
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

	var BookingToUpdate1Variable BookingsModel

	err = tx.Get(&BookingToUpdate1Variable, `
        SELECT id, booking_date, checkout_session_id, guest_id, listing_id, status, request_detail
        FROM bookings
        WHERE (id = $1 OR $1 IS NULL)`,
		Id14Variable)

	if err != nil {
		return
	}

	var Approved2Variable = "approved"

	_, err = tx.Exec(`
        UPDATE bookings
        SET status = $1
        WHERE id = $2`,
		Approved2Variable, BookingToUpdate1Variable.Id)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(nil))
}
