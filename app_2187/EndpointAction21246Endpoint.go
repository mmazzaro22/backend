package main

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"net/http"
	"utils"
)

// Action21246Endpoint is "get" request handler for "action_21246" endpoint.
func Action21246Endpoint(w http.ResponseWriter, r *http.Request) {

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

	var CreatorId3Variable int

	var dec20525 = json.NewDecoder(r.Body)
	if err = dec20525.Decode(&CreatorId3Variable); err != nil {
		errStatus = http.StatusBadRequest
		errResponse = utils.InvalidJSONErrorResponse(err)
		return
	}

	var UserHasProperty1Variable bool
	err = tx.Get(&UserHasProperty1Variable, `
	SELECT EXISTS (
		SELECT 1
		FROM properties
		WHERE (creator_id = $1 OR $1 IS NULL)
	)
`, CreatorId3Variable)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(UserHasProperty1Variable))
}
