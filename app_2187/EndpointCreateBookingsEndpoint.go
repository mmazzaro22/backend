package main

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"net/http"
	"utils"
)

// Create bookings
func CreateBookingsEndpoint(w http.ResponseWriter, r *http.Request) {

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

	var Input6Variable BookingsModel

	var dec19126 = json.NewDecoder(r.Body)
	if err = dec19126.Decode(&Input6Variable); err != nil {
		errStatus = http.StatusBadRequest
		errResponse = utils.InvalidJSONErrorResponse(err)
		return
	}

	var tx *sqlx.Tx
	tx, err = pg.Beginx()
	if err != nil {
		return
	}
	defer tx.Rollback()

	var nullableStatus39142 *string

	if Input6Variable.Status != nil {
		nullableStatus39142 = Input6Variable.Status

	}

	var nullableRequestDetail39142 *string

	if Input6Variable.RequestDetail != nil {
		nullableRequestDetail39142 = Input6Variable.RequestDetail

	}

	var nullableListingId39142 *int

	if Input6Variable.ListingId != nil {
		nullableListingId39142 = Input6Variable.ListingId

	}

	_, err = tx.Exec(`
		INSERT INTO bookings
		(status, request_detail, booking_date, checkout_session_id, guest_id, listing_id)
		VALUES($1, $2, $3, $4, $5, $6)`,
		nullableStatus39142, nullableRequestDetail39142, Input6Variable.BookingDate, Input6Variable.CheckoutSessionId, Input6Variable.GuestId, nullableListingId39142)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(nil))
}
