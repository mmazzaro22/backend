package main

import (
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
	"utils"
)

// Action21369Endpoint is "delete" request handler for "action_21369" endpoint.
func Action21369Endpoint(w http.ResponseWriter, r *http.Request) {

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

	var UserId5Variable int

	var PropertyId10Variable int

	if r.URL != nil {
		var queries20645 = r.URL.Query()
		if queries20645 != nil && len(queries20645) > 0 {
			if vals20645, ok := queries20645["user_id"]; ok && len(vals20645) > 0 {
				UserId5Variable, err = strconv.Atoi(vals20645[0])
				if err != nil {
					errStatus = http.StatusBadRequest
					errResponse = utils.Response{Error: &utils.ValError{Message: "invalid user_id", Code: utils.ErrCodeBadInput}}
					return
				}
			}
			if vals20645, ok := queries20645["property_id"]; ok && len(vals20645) > 0 {
				PropertyId10Variable, err = strconv.Atoi(vals20645[0])
				if err != nil {
					errStatus = http.StatusBadRequest
					errResponse = utils.Response{Error: &utils.ValError{Message: "invalid property_id", Code: utils.ErrCodeBadInput}}
					return
				}
			}
		}
	}

	var tx *sqlx.Tx
	tx, err = pg.Beginx()
	if err != nil {
		return
	}
	defer tx.Rollback()

	var Image1Variable PropertyImagesModel

	err = tx.Get(&Image1Variable, `
        SELECT id, file_name, property_id, file, user_id
        FROM property_images
        WHERE (user_id = $1 OR $1 IS NULL) OR (property_id = $2 OR $2 IS NULL)`,
		UserId5Variable, PropertyId10Variable)

	if err != nil {
		return
	}

	_, err = tx.Exec(`
        DELETE FROM property_images
        WHERE id = $1`,
		Image1Variable.Id)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(nil))
}
