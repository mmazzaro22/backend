package main

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
	"utils"
)

// Action21475Endpoint is "get" request handler for "action_21475" endpoint.
func Action21475Endpoint(w http.ResponseWriter, r *http.Request) {

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

	var CreatorId4Variable int

	var paths20751 = mux.Vars(r)
	if paths20751 != nil && len(paths20751) > 0 {
		if val20751, ok := paths20751["creator_id"]; ok {
			CreatorId4Variable, err = strconv.Atoi(val20751)
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

	var GetPropertyCustom1Variable []GetPropertiesEPRConstruct

	err = tx.Select(&GetPropertyCustom1Variable, `
        SELECT p.*, pi.file
FROM properties p
LEFT JOIN (
    SELECT DISTINCT ON (property_id) property_id, file
    FROM property_images
) pi ON p.id = pi.property_id
WHERE p.creator_id = $1`,
		CreatorId4Variable)

	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(GetPropertyCustom1Variable))
}
