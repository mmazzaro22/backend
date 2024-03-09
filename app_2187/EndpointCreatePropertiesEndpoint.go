package main

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"net/http"
	"utils"
)

// Create properties
func CreatePropertiesEndpoint(w http.ResponseWriter, r *http.Request) {

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

	var iLoginSessionUserID int
	iLoginSessionUserID, err = getLoginSessionUserID(r)
	if err != nil {
		errStatus = http.StatusUnauthorized
		errResponse = utils.Response{Error: &utils.ValError{Message: "login required", Code: utils.ErrCodeUnAuthorized}}
		return
	}

	var tx *sqlx.Tx
	tx, err = pg.Beginx()
	if err != nil {
		return
	}
	defer tx.Rollback()

	var CurrentUser1Variable UsersModel

	err = tx.Get(&CurrentUser1Variable, `
        SELECT id, email, password, stripe_customer_id, profile_picture, type_of_user, last_log_in_date_time, bio, first_name, last_name, phone_number, partner_site_id
        FROM users
        WHERE (id = $1 OR $1 IS NULL)`,
		iLoginSessionUserID)

	if err != nil {
		return
	}

	var validator18080 *validator.Validate
	validator18080 = validator.New()

	var Input7Variable CreatePropertiesEPIConstruct

	var dec18080 = json.NewDecoder(r.Body)
	if err = dec18080.Decode(&Input7Variable); err != nil {
		errStatus = http.StatusBadRequest
		errResponse = utils.InvalidJSONErrorResponse(err)
		return
	}
	err = validator18080.Struct(Input7Variable)
	var ve18080 validator.ValidationErrors
	if errors.As(err, &ve18080) {
		errStatus = http.StatusBadRequest
		errResponse = utils.Response{Error: &utils.ValError{Param: ve18080[0].Field(), Message: MsgToTag(ve18080[0]), Code: "invalid_json"}}
		return
	} else if err != nil {
		errStatus = http.StatusBadRequest
		errResponse = utils.InvalidJSONErrorResponse(err)
		return
	}

	var PendingStatus1Variable = "Pending"

	var nullableState37414 *string

	if Input7Variable.State != nil {
		nullableState37414 = new(string)
		*nullableState37414 = Input7Variable.State.string
	}

	var nullableLandType37414 *string

	if Input7Variable.LandType != nil {
		nullableLandType37414 = new(string)
		*nullableLandType37414 = Input7Variable.LandType.string
	}

	var nullableAmenity37414 *string

	if Input7Variable.Amenity != nil {
		nullableAmenity37414 = new(string)
		*nullableAmenity37414 = Input7Variable.Amenity.string
	}

	var NewRecord7Variable PropertiesModel
	err = tx.Get(&NewRecord7Variable, `
		INSERT INTO properties
		(address, city, state, land_type, name, description, creator_id, acres, zip, amenity, status)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING *`,
		Input7Variable.Address, Input7Variable.City, nullableState37414, nullableLandType37414, Input7Variable.Name, Input7Variable.Description, CurrentUser1Variable.Id, Input7Variable.Acres, Input7Variable.Zip, nullableAmenity37414, PendingStatus1Variable)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(NewRecord7Variable))
}
