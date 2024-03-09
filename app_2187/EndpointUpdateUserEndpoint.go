package main

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"utils"
)

// UpdateUserEndpoint is "put" request handler for "updateUser" endpoint.
func UpdateUserEndpoint(w http.ResponseWriter, r *http.Request) {

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

	var UserToUpdate2Variable UsersModel

	var dec18037 = json.NewDecoder(r.Body)
	if err = dec18037.Decode(&UserToUpdate2Variable); err != nil {
		errStatus = http.StatusBadRequest
		errResponse = utils.InvalidJSONErrorResponse(err)
		return
	}

	var nullableBio41725 *string

	if UserToUpdate2Variable.Bio != nil {
		nullableBio41725 = UserToUpdate2Variable.Bio

	}

	var nullableFirstName41725 *string

	if UserToUpdate2Variable.FirstName != nil {
		nullableFirstName41725 = UserToUpdate2Variable.FirstName

	}

	var nullablePhoneNumber41725 *string

	if UserToUpdate2Variable.PhoneNumber != nil {
		nullablePhoneNumber41725 = UserToUpdate2Variable.PhoneNumber

	}

	var nullablePartnerSiteId41725 *int

	if UserToUpdate2Variable.PartnerSiteId != nil {
		nullablePartnerSiteId41725 = UserToUpdate2Variable.PartnerSiteId

	}

	var nullableStripeCustomerId41725 *string

	if CurrentUser1Variable.StripeCustomerId != nil {
		nullableStripeCustomerId41725 = CurrentUser1Variable.StripeCustomerId

	}

	var nullableTypeOfUser41725 *string

	if CurrentUser1Variable.TypeOfUser != nil {
		nullableTypeOfUser41725 = CurrentUser1Variable.TypeOfUser

	}

	var nullableLastName41725 *string

	if UserToUpdate2Variable.LastName != nil {
		nullableLastName41725 = UserToUpdate2Variable.LastName

	}

	var hashedPassword41725 []byte

	hashedPassword41725, err = bcrypt.GenerateFromPassword([]byte(CurrentUser1Variable.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	_, err = tx.Exec(`
        UPDATE users
        SET bio = $1, first_name = $2, phone_number = $3, partner_site_id = $4, id = $5, email = $6, password = $7, stripe_customer_id = $8, last_log_in_date_time = $9, type_of_user = $10, last_name = $11
        WHERE id = $12`,
		nullableBio41725, nullableFirstName41725, nullablePhoneNumber41725, nullablePartnerSiteId41725, CurrentUser1Variable.Id, CurrentUser1Variable.Email, hashedPassword41725, nullableStripeCustomerId41725, CurrentUser1Variable.LastLogInDateTime, nullableTypeOfUser41725, nullableLastName41725, UserToUpdate2Variable.Id)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(UserToUpdate2Variable))
}
