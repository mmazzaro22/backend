package main

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
	"utils"
)

// Update properties
func UpdatePropertiesEndpoint(w http.ResponseWriter, r *http.Request) {

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

	var Id4Variable int

	var paths18081 = mux.Vars(r)
	if paths18081 != nil && len(paths18081) > 0 {
		if val18081, ok := paths18081["id"]; ok {
			Id4Variable, err = strconv.Atoi(val18081)
			if err != nil {
				errStatus = http.StatusBadRequest
				errResponse = utils.Response{Error: &utils.ValError{Message: "invalid id", Code: utils.ErrCodeBadInput}}
				return
			}
		}
	}

	var validator18081 *validator.Validate
	validator18081 = validator.New()

	var Input5Variable UpdatePropertiesEPIConstruct

	var dec18081 = json.NewDecoder(r.Body)
	if err = dec18081.Decode(&Input5Variable); err != nil {
		errStatus = http.StatusBadRequest
		errResponse = utils.InvalidJSONErrorResponse(err)
		return
	}
	err = validator18081.Struct(Input5Variable)
	var ve18081 validator.ValidationErrors
	if errors.As(err, &ve18081) {
		errStatus = http.StatusBadRequest
		errResponse = utils.Response{Error: &utils.ValError{Param: ve18081[0].Field(), Message: MsgToTag(ve18081[0]), Code: "invalid_json"}}
		return
	} else if err != nil {
		errStatus = http.StatusBadRequest
		errResponse = utils.InvalidJSONErrorResponse(err)
		return
	}

	var tx *sqlx.Tx
	tx, err = pg.Beginx()
	if err != nil {
		return
	}
	defer tx.Rollback()

	var PropertyToUpdate1Variable PropertiesModel

	err = tx.Get(&PropertyToUpdate1Variable, `
        SELECT id, name, creator_id, description, address, city, state, zip, images, land_type, amenity, status, acres
        FROM properties
        WHERE (id = $1 OR $1 IS NULL)`,
		Id4Variable)

	if err != nil {
		return
	}

	var nullableDescription37415 *string

	if Input5Variable.Description != nil {
		nullableDescription37415 = Input5Variable.Description

	}

	var nullableState37415 *string

	if Input5Variable.State != nil {
		nullableState37415 = Input5Variable.State

	}

	var nullableLandType37415 *string

	if Input5Variable.LandType != nil {
		nullableLandType37415 = Input5Variable.LandType

	}

	var nullableAmenity37415 *string

	if Input5Variable.Amenity != nil {
		nullableAmenity37415 = new(string)
		*nullableAmenity37415 = Input5Variable.Amenity.string
	}

	var nullableName37415 *string

	if Input5Variable.Name != nil {
		nullableName37415 = Input5Variable.Name

	}

	var nullableAddress37415 *string

	if Input5Variable.Address != nil {
		nullableAddress37415 = Input5Variable.Address

	}

	var nullableCity37415 *string

	if Input5Variable.City != nil {
		nullableCity37415 = Input5Variable.City

	}

	var nullableZip37415 *int

	if Input5Variable.Zip != nil {
		nullableZip37415 = Input5Variable.Zip

	}

	var nullableStatus37415 *string

	if Input5Variable.Status != nil {
		nullableStatus37415 = Input5Variable.Status

	}

	_, err = tx.Exec(`
        UPDATE properties
        SET description = $1, state = $2, land_type = $3, amenity = $4, name = $5, address = $6, city = $7, zip = $8, status = $9
        WHERE id = $10`,
		nullableDescription37415, nullableState37415, nullableLandType37415, nullableAmenity37415, nullableName37415, nullableAddress37415, nullableCity37415, nullableZip37415, nullableStatus37415, PropertyToUpdate1Variable.Id)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(nil))
}
