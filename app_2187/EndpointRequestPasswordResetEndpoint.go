package main

import (
	"bytes"
	"crypto/tls"
	"github.com/go-gomail/gomail"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strings"
	"utils"
)

// Request Password Reset
func RequestPasswordResetEndpoint(w http.ResponseWriter, r *http.Request) {

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

	var Email1Variable string

	if r.URL != nil {
		var queries18036 = r.URL.Query()
		if queries18036 != nil && len(queries18036) > 0 {
			if vals18036, ok := queries18036["email"]; ok && len(vals18036) > 0 {
				Email1Variable = strings.Join(vals18036, "\n")
			}
		}
	}

	var tx *sqlx.Tx
	tx, err = pg.Beginx()
	if err != nil {
		return
	}
	defer tx.Rollback()

	var FromEmail1Variable = "test@test.dittofi.com"

	var User1Variable UsersModel

	err = tx.Get(&User1Variable, `
        SELECT id, email, password, stripe_customer_id, profile_picture, type_of_user, last_log_in_date_time, bio, first_name, last_name, phone_number, partner_site_id
        FROM users
        WHERE (email = $1 OR $1 IS NULL)`,
		Email1Variable)

	if err != nil {
		return
	}

	var Subject1Variable = "Reset Password"

	var Token1Variable string

	Token1Variable, err = EncodeJWTCustomFn(User1Variable.Id)
	if err != nil {
		return
	}

	var message37367 = gomail.NewMessage()
	message37367.SetAddressHeader("To", Email1Variable, Email1Variable)
	message37367.SetAddressHeader("From", FromEmail1Variable, FromEmail1Variable)
	message37367.SetAddressHeader("Cc", Email1Variable, Email1Variable)
	message37367.SetAddressHeader("Bcc", Email1Variable, Email1Variable)
	message37367.SetHeader("Subject", Subject1Variable)

	var parameters37367 = struct {
		Token interface{}
	}{
		Token: Token1Variable,
	}

	var emailBodyBuf37367 = new(bytes.Buffer)
	err = templates.ExecuteTemplate(emailBodyBuf37367, "Reset Password Email.tmp", parameters37367)
	if err != nil {
		return
	}

	message37367.SetBody("text/html", emailBodyBuf37367.String())

	var dialer37367 *gomail.Dialer
	dialer37367 = gomail.NewDialer(SMTPHost, SMTPPort, SMTPUsername, SMTPPassword)
	dialer37367.SSL = true
	dialer37367.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	var sendCloser37367 gomail.SendCloser
	sendCloser37367, err = dialer37367.Dial()
	if err != nil {
		return
	}

	err = gomail.Send(sendCloser37367, message37367)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(nil))
}
