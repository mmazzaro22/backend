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

// Action18697Endpoint is "post" request handler for "action_18697" endpoint.
func Action18697Endpoint(w http.ResponseWriter, r *http.Request) {

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

	var PropertyId2Variable int

	var paths18040 = mux.Vars(r)
	if paths18040 != nil && len(paths18040) > 0 {
		if val18040, ok := paths18040["property_id"]; ok {
			PropertyId2Variable, err = strconv.Atoi(val18040)
			if err != nil {
				errStatus = http.StatusBadRequest
				errResponse = utils.Response{Error: &utils.ValError{Message: "invalid property_id", Code: utils.ErrCodeBadInput}}
				return
			}
		}
	}

	var validator18040 *validator.Validate
	validator18040 = validator.New()

	var Input3Variable CreateListingEPIConstruct

	var dec18040 = json.NewDecoder(r.Body)
	if err = dec18040.Decode(&Input3Variable); err != nil {
		errStatus = http.StatusBadRequest
		errResponse = utils.InvalidJSONErrorResponse(err)
		return
	}
	err = validator18040.Struct(Input3Variable)
	var ve18040 validator.ValidationErrors
	if errors.As(err, &ve18040) {
		errStatus = http.StatusBadRequest
		errResponse = utils.Response{Error: &utils.ValError{Param: ve18040[0].Field(), Message: MsgToTag(ve18040[0]), Code: "invalid_json"}}
		return
	} else if err != nil {
		errStatus = http.StatusBadRequest
		errResponse = utils.InvalidJSONErrorResponse(err)
		return
	}

	var nullableOrderTypeId37378 *int

	if Input3Variable.OrderTypeId != nil {
		nullableOrderTypeId37378 = Input3Variable.OrderTypeId

	}

	var nullablePrice37378 *int

	if Input3Variable.Price != nil {
		nullablePrice37378 = Input3Variable.Price

	}

	var nullableCheckInTimeId37378 *int

	if Input3Variable.CheckInTimeId != nil {
		nullableCheckInTimeId37378 = Input3Variable.CheckInTimeId

	}

	var nullableCheckOutTimeId37378 *int

	if Input3Variable.CheckOutTimeId != nil {
		nullableCheckOutTimeId37378 = Input3Variable.CheckOutTimeId

	}

	_, err = tx.Exec(`
		INSERT INTO listings
		(guest_limit, property_id, order_type_id, creator_id, price, partner_site_id, check_in_time_id, check_out_time_id)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8)`,
		Input3Variable.GuestLimit, PropertyId2Variable, nullableOrderTypeId37378, CurrentUser1Variable.Id, nullablePrice37378, Input3Variable.PartnerSiteId, nullableCheckInTimeId37378, nullableCheckOutTimeId37378)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(nil))
}
