package main

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
	"utils"
)

// Action21335Endpoint is "put" request handler for "action_21335" endpoint.
func Action21335Endpoint(w http.ResponseWriter, r *http.Request) {

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

	var UserId2Variable int

	if r.URL != nil {
		var queries20611 = r.URL.Query()
		if queries20611 != nil && len(queries20611) > 0 {
			if vals20611, ok := queries20611["user_id"]; ok && len(vals20611) > 0 {
				UserId2Variable, err = strconv.Atoi(vals20611[0])
				if err != nil {
					errStatus = http.StatusBadRequest
					errResponse = utils.Response{Error: &utils.ValError{Message: "invalid user_id", Code: utils.ErrCodeBadInput}}
					return
				}
			}
		}
	}

	var validator20611 *validator.Validate
	validator20611 = validator.New()

	var Input9Variable UpdateUserEmailConstruct

	var dec20611 = json.NewDecoder(r.Body)
	if err = dec20611.Decode(&Input9Variable); err != nil {
		errStatus = http.StatusBadRequest
		errResponse = utils.InvalidJSONErrorResponse(err)
		return
	}
	err = validator20611.Struct(Input9Variable)
	var ve20611 validator.ValidationErrors
	if errors.As(err, &ve20611) {
		errStatus = http.StatusBadRequest
		errResponse = utils.Response{Error: &utils.ValError{Param: ve20611[0].Field(), Message: MsgToTag(ve20611[0]), Code: "invalid_json"}}
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

	var User2Variable UsersModel

	err = tx.Get(&User2Variable, `
        SELECT id, email, password, stripe_customer_id, profile_picture, type_of_user, last_log_in_date_time, bio, first_name, last_name, phone_number, partner_site_id
        FROM users
        WHERE (id = $1 OR $1 IS NULL)`,
		UserId2Variable)

	if err != nil {
		return
	}

	_, err = tx.Exec(`
        UPDATE users
        SET email = $1
        WHERE id = $2`,
		Input9Variable.Email, User2Variable.Id)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(Input9Variable))
}
