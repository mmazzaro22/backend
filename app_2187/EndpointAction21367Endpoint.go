package main

import (
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
	"utils"
)

// Action21367Endpoint is "get" request handler for "action_21367" endpoint.
func Action21367Endpoint(w http.ResponseWriter, r *http.Request) {

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

	var UserId4Variable int

	if r.URL != nil {
		var queries20643 = r.URL.Query()
		if queries20643 != nil && len(queries20643) > 0 {
			if vals20643, ok := queries20643["user_id"]; ok && len(vals20643) > 0 {
				UserId4Variable, err = strconv.Atoi(vals20643[0])
				if err != nil {
					errStatus = http.StatusBadRequest
					errResponse = utils.Response{Error: &utils.ValError{Message: "invalid user_id", Code: utils.ErrCodeBadInput}}
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

	var ProfilePicture1Variable PropertyImagesModel

	err = tx.Get(&ProfilePicture1Variable, `
        SELECT id, file_name, property_id, file, user_id
        FROM property_images
        WHERE (user_id = $1 OR $1 IS NULL)`,
		UserId4Variable)

	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(ProfilePicture1Variable))
}
