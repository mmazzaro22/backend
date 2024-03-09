package main

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
	"utils"
)

// Action20159Endpoint is "put" request handler for "action_20159" endpoint.
func Action20159Endpoint(w http.ResponseWriter, r *http.Request) {

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

	var Id15Variable int

	var paths19470 = mux.Vars(r)
	if paths19470 != nil && len(paths19470) > 0 {
		if val19470, ok := paths19470["id"]; ok {
			Id15Variable, err = strconv.Atoi(val19470)
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

	var BookingToReject1Variable BookingsModel

	err = tx.Get(&BookingToReject1Variable, `
        SELECT id, booking_date, checkout_session_id, guest_id, listing_id, status, request_detail
        FROM bookings
        WHERE (id = $1 OR $1 IS NULL)`,
		Id15Variable)

	if err != nil {
		return
	}

	var Rejected1Variable = "rejected"

	_, err = tx.Exec(`
        UPDATE bookings
        SET status = $1
        WHERE id = $2`,
		Rejected1Variable, BookingToReject1Variable.Id)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(nil))
}
