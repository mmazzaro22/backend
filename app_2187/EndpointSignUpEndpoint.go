package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"utils"
)

// SignUp
func SignUpEndpoint(w http.ResponseWriter, r *http.Request) {

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

	var validator18032 *validator.Validate
	validator18032 = validator.New()

	var Input1Variable SignUpEPIConstruct

	var dec18032 = json.NewDecoder(r.Body)
	if err = dec18032.Decode(&Input1Variable); err != nil {
		errStatus = http.StatusBadRequest
		errResponse = utils.InvalidJSONErrorResponse(err)
		return
	}
	err = validator18032.Struct(Input1Variable)
	var ve18032 validator.ValidationErrors
	if errors.As(err, &ve18032) {
		errStatus = http.StatusBadRequest
		errResponse = utils.Response{Error: &utils.ValError{Param: ve18032[0].Field(), Message: MsgToTag(ve18032[0]), Code: "invalid_json"}}
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

	var hashedPassword37364 []byte

	hashedPassword37364, err = bcrypt.GenerateFromPassword([]byte(Input1Variable.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	var defaultValueLastLogInDateTime37364 = "now()"

	_, err = tx.Exec(`
		INSERT INTO users
		(last_name, email, password, first_name, last_log_in_date_time)
		VALUES($1, $2, $3, $4, $5)`,
		Input1Variable.LastName, Input1Variable.Email, hashedPassword37364, Input1Variable.FirstName, defaultValueLastLogInDateTime37364)
	if err != nil {
		return
	}

	var CurrentUser1Variable UsersModel

	err = tx.Get(&CurrentUser1Variable, `
        SELECT id, email, password, stripe_customer_id, profile_picture, type_of_user, last_log_in_date_time, bio, first_name, last_name, phone_number, partner_site_id
        FROM users
        WHERE (email = $1 OR $1 IS NULL)`,
		Input1Variable.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("username not found")
			errStatus = http.StatusForbidden
			errResponse = utils.Response{Error: &utils.ValError{Message: "invalid login credentials", Code: utils.ErrCodeUnAuthorized}}
		}
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(CurrentUser1Variable.Password), []byte(Input1Variable.Password)); err != nil {
		err = fmt.Errorf("invalid password")
		errStatus = http.StatusForbidden
		errResponse = utils.Response{Error: &utils.ValError{Message: "invalid login credentials", Code: utils.ErrCodeUnAuthorized}}
		return
	}

	err = addLoginSession(&w, r, CurrentUser1Variable.Id)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(nil))
}
