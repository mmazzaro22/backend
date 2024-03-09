package main

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
	"utils"
)

// Get One listings
func GetOneListingsEndpoint(w http.ResponseWriter, r *http.Request) {

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

	var Id5Variable int

	var paths18108 = mux.Vars(r)
	if paths18108 != nil && len(paths18108) > 0 {
		if val18108, ok := paths18108["id"]; ok {
			Id5Variable, err = strconv.Atoi(val18108)
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

	var Listing5Variable GetListingsEPRConstruct

	err = tx.Get(&Listing5Variable, `
        SELECT a.*, b.*, c.first_name, c.last_name, c.email, c.profile_picture, 
       d.id AS order_type_id, 
       in_time.hour AS in_time_hour, out_time.hour AS out_time_hour,
       f.name AS order_type_name,
       b.id AS property_id, a.id AS id
FROM listings a 
LEFT JOIN properties b ON a.property_id = b.id 
LEFT JOIN users c ON b.creator_id = c.id 
LEFT JOIN order_types d ON a.order_type_id = d.id 
LEFT JOIN check_in_out_time in_time ON a.check_in_time_id = in_time.id
LEFT JOIN check_in_out_time out_time ON a.check_out_time_id = out_time.id
LEFT JOIN order_types f ON f.id = d.id
WHERE a.id = $1`,
		Id5Variable)

	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(Listing5Variable))
}
