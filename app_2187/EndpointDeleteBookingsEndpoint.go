package main

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
	"utils"
)

// Delete bookings
func DeleteBookingsEndpoint(w http.ResponseWriter, r *http.Request) {

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

	var Id10Variable int

	var paths19128 = mux.Vars(r)
	if paths19128 != nil && len(paths19128) > 0 {
		if val19128, ok := paths19128["id"]; ok {
			Id10Variable, err = strconv.Atoi(val19128)
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

	var ModelToDelete1Variable BookingsModel

	err = tx.Get(&ModelToDelete1Variable, `
        SELECT id, booking_date, checkout_session_id, guest_id, listing_id, status, request_detail
        FROM bookings
        WHERE (id = $1 OR $1 IS NULL)`,
		Id10Variable)

	if err != nil {
		return
	}

	_, err = tx.Exec(`
        DELETE FROM bookings
        WHERE id = $1`,
		ModelToDelete1Variable.Id)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(nil))
}
