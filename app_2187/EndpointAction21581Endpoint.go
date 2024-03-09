package main

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"utils"
)

// Action21581Endpoint is "get" request handler for "action_21581" endpoint.
func Action21581Endpoint(w http.ResponseWriter, r *http.Request) {

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

	var Email3Variable string

	var Password2Variable string

	if r.URL != nil {
		var queries20856 = r.URL.Query()
		if queries20856 != nil && len(queries20856) > 0 {
			if vals20856, ok := queries20856["email"]; ok && len(vals20856) > 0 {
				Email3Variable = strings.Join(vals20856, "\n")
			}
			if vals20856, ok := queries20856["password"]; ok && len(vals20856) > 0 {
				Password2Variable = strings.Join(vals20856, "\n")
			}
		}
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
        WHERE (email = $1 OR $1 IS NULL)`,
		Email3Variable)

	if err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("username not found")
			errStatus = http.StatusForbidden
			errResponse = utils.Response{Error: &utils.ValError{Message: "invalid login credentials", Code: utils.ErrCodeUnAuthorized}}
		}
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(CurrentUser1Variable.Password), []byte(Password2Variable)); err != nil {
		err = fmt.Errorf("invalid password")
		errStatus = http.StatusForbidden
		errResponse = utils.Response{Error: &utils.ValError{Message: "invalid login credentials", Code: utils.ErrCodeUnAuthorized}}
		return
	}

	var NewlyLoggedInUser1Variable = CurrentUser1Variable

	err = addLoginSession(&w, r, CurrentUser1Variable.Id)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(NewlyLoggedInUser1Variable))
}
