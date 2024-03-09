package main

import (
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
	"utils"
)

// Action21272Endpoint is "get" request handler for "action_21272" endpoint.
func Action21272Endpoint(w http.ResponseWriter, r *http.Request) {

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

	var PropertyId9Variable int

	if r.URL != nil {
		var queries20550 = r.URL.Query()
		if queries20550 != nil && len(queries20550) > 0 {
			if vals20550, ok := queries20550["property_id"]; ok && len(vals20550) > 0 {
				PropertyId9Variable, err = strconv.Atoi(vals20550[0])
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

	var PropertyImages1Variable []PropertyImagesModel

	err = tx.Select(&PropertyImages1Variable, `
        SELECT id, file_name, property_id, file, user_id
        FROM property_images
        WHERE (property_id = $1 OR $1 IS NULL)`,
		PropertyId9Variable)

	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(PropertyImages1Variable))
}
