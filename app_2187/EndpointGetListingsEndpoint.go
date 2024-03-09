package main

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
	"utils"
)

// Get listings
func GetListingsEndpoint(w http.ResponseWriter, r *http.Request) {

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

	var PropertyId1Variable int

	var paths18109 = mux.Vars(r)
	if paths18109 != nil && len(paths18109) > 0 {
		if val18109, ok := paths18109["property_id"]; ok {
			PropertyId1Variable, err = strconv.Atoi(val18109)
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

	var Lisings1Variable []GetPropertyListingEPRConstruct

	err = tx.Select(&Lisings1Variable, `
        SELECT b.*, a.*, c.address, c.city, c.state, d.hour
FROM listings a
LEFT JOIN partner_sites b ON a.partner_site_id = b.id
LEFT JOIN properties c ON a.property_id = c.id
LEFT JOIN check_in_out_time d ON (a.check_in_time_id = d.id) AND (a.check_out_time_id = d.id)
WHERE a.property_id = $1`,
		PropertyId1Variable)

	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(Lisings1Variable))
}
