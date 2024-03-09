package main

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
	"utils"
)

// Get bookings
func GetBookingsEndpoint(w http.ResponseWriter, r *http.Request) {

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

	var CreatorId2Variable int

	var paths19125 = mux.Vars(r)
	if paths19125 != nil && len(paths19125) > 0 {
		if val19125, ok := paths19125["creator_id"]; ok {
			CreatorId2Variable, err = strconv.Atoi(val19125)
			if err != nil {
				errStatus = http.StatusBadRequest
				errResponse = utils.Response{Error: &utils.ValError{Message: "invalid creator_id", Code: utils.ErrCodeBadInput}}
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

	var Bookings1Variable []GetBookingsEPRConstruct

	err = tx.Select(&Bookings1Variable, `
        SELECT a.*, b.first_name, b.last_name, d.name, c.creator_id
FROM bookings a 
LEFT JOIN users b ON a.guest_id = b.id
LEFT JOIN listings c ON a.listing_id = c.id
LEFT JOIN partner_sites d ON c.partner_site_id = d.id 
WHERE c.creator_id = $1`,
		CreatorId2Variable)

	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(Bookings1Variable))
}
