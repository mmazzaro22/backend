package main

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
	"utils"
)

// Delete properties
func DeletePropertiesEndpoint(w http.ResponseWriter, r *http.Request) {

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

	var Id3Variable int

	var paths18082 = mux.Vars(r)
	if paths18082 != nil && len(paths18082) > 0 {
		if val18082, ok := paths18082["id"]; ok {
			Id3Variable, err = strconv.Atoi(val18082)
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

	var PropertyToDelete1Variable PropertiesModel

	err = tx.Get(&PropertyToDelete1Variable, `
        SELECT id, name, creator_id, description, address, city, state, zip, images, land_type, amenity, status, acres
        FROM properties
        WHERE (id = $1 OR $1 IS NULL)`,
		Id3Variable)

	if err != nil {
		return
	}

	_, err = tx.Exec(`
        DELETE FROM properties
        WHERE id = $1`,
		PropertyToDelete1Variable.Id)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(nil))
}
