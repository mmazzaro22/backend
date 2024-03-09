package main

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
	"utils"
)

// Get One properties
func GetOnePropertiesEndpoint(w http.ResponseWriter, r *http.Request) {

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

	var Id2Variable int

	var paths18078 = mux.Vars(r)
	if paths18078 != nil && len(paths18078) > 0 {
		if val18078, ok := paths18078["id"]; ok {
			Id2Variable, err = strconv.Atoi(val18078)
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

	var Data2Variable PropertiesModel

	err = tx.Get(&Data2Variable, `
        SELECT id, name, creator_id, description, address, city, state, zip, images, land_type, amenity, status, acres
        FROM properties
        WHERE (id = $1 OR $1 IS NULL)`,
		Id2Variable)

	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(Data2Variable))
}
