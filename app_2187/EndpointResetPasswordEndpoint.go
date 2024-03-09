package main

import (
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"utils"
)

// Reset Password
func ResetPasswordEndpoint(w http.ResponseWriter, r *http.Request) {

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

	var Password1Variable string

	var Token2Variable string

	if r.URL != nil {
		var queries18035 = r.URL.Query()
		if queries18035 != nil && len(queries18035) > 0 {
			if vals18035, ok := queries18035["password"]; ok && len(vals18035) > 0 {
				Password1Variable = strings.Join(vals18035, "\n")
			}
			if vals18035, ok := queries18035["token"]; ok && len(vals18035) > 0 {
				Token2Variable = strings.Join(vals18035, "\n")
			}
		}
	}

	var tx *sqlx.Tx
	tx, err = pg.Beginx()
	if err != nil {
		return
	}
	defer tx.Rollback()

	var True1Variable = true

	var Success1Variable bool

	var UserId1Variable int

	Success1Variable, UserId1Variable, err = DecodeJWTCustomFn(Token2Variable)
	if err != nil {
		return
	}

	if Success1Variable == True1Variable {
		var UserToUpdate1Variable UsersModel

		err = tx.Get(&UserToUpdate1Variable, `
        SELECT id, email, password, stripe_customer_id, profile_picture, type_of_user, last_log_in_date_time, bio, first_name, last_name, phone_number, partner_site_id
        FROM users
        WHERE (id = $1 OR $1 IS NULL)`,
			UserId1Variable)

		if err != nil {
			return
		}

		var hashedPassword37370 []byte

		hashedPassword37370, err = bcrypt.GenerateFromPassword([]byte(Password1Variable), bcrypt.DefaultCost)
		if err != nil {
			return
		}

		_, err = tx.Exec(`
        UPDATE users
        SET password = $1
        WHERE id = $2`,
			hashedPassword37370, UserToUpdate1Variable.Id)
		if err != nil {
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(nil))
}
