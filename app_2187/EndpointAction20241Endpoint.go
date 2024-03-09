package main

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
	"utils"
)

// Action20241Endpoint is "get" request handler for "action_20241" endpoint.
func Action20241Endpoint(w http.ResponseWriter, r *http.Request) {

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

	var ReceiverId2Variable int

	var paths19543 = mux.Vars(r)
	if paths19543 != nil && len(paths19543) > 0 {
		if val19543, ok := paths19543["receiver_id"]; ok {
			ReceiverId2Variable, err = strconv.Atoi(val19543)
			if err != nil {
				errStatus = http.StatusBadRequest
				errResponse = utils.Response{Error: &utils.ValError{Message: "invalid receiver_id", Code: utils.ErrCodeBadInput}}
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

	var GetMessageEPR1Variable []GetMessageEPRConstruct

	err = tx.Select(&GetMessageEPR1Variable, `
        SELECT a.*, 
       b.id, 
       b.email, 
       b.first_name AS sender_first_name, 
       b.last_name AS sender_last_name, 
       b.profile_picture AS sender_profile_picture, 
       b2.first_name AS receiver_first_name, 
       b2.last_name AS receiver_last_name, 
       b2.profile_picture AS receiver_profile_picture
FROM conversation a 
LEFT JOIN users b ON a.sender_id = b.id 
LEFT JOIN users b2 ON a.receiver_id = b2.id 
WHERE a.receiver_id = $1;`,
		ReceiverId2Variable)

	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(GetMessageEPR1Variable))
}
