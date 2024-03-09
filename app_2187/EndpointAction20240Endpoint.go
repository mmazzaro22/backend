package main

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
	"time"
	"utils"
)

// Action20240Endpoint is "post" request handler for "action_20240" endpoint.
func Action20240Endpoint(w http.ResponseWriter, r *http.Request) {

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

	var Message2Variable string

	var SenderId1Variable int

	var ReceiverId1Variable int

	var paths19542 = mux.Vars(r)
	if paths19542 != nil && len(paths19542) > 0 {
		if val19542, ok := paths19542["message"]; ok {
			Message2Variable = val19542
		}
		if val19542, ok := paths19542["sender_id"]; ok {
			SenderId1Variable, err = strconv.Atoi(val19542)
			if err != nil {
				errStatus = http.StatusBadRequest
				errResponse = utils.Response{Error: &utils.ValError{Message: "invalid sender_id", Code: utils.ErrCodeBadInput}}
				return
			}
		}
		if val19542, ok := paths19542["receiver_id"]; ok {
			ReceiverId1Variable, err = strconv.Atoi(val19542)
			if err != nil {
				errStatus = http.StatusBadRequest
				errResponse = utils.Response{Error: &utils.ValError{Message: "invalid receiver_id", Code: utils.ErrCodeBadInput}}
				return
			}
		}
	}

	var SentDate1Variable time.Time

	if r.URL != nil {
		var queries19542 = r.URL.Query()
		if queries19542 != nil && len(queries19542) > 0 {
			if vals19542, ok := queries19542["sent_date"]; ok && len(vals19542) > 0 {
				SentDate1Variable, err = time.Parse(time.RFC3339, vals19542[0])
				if err != nil {
					errStatus = http.StatusBadRequest
					errResponse = utils.Response{Error: &utils.ValError{Message: "invalid sent_date", Code: utils.ErrCodeBadInput}}
					return
				}
			}
		}
	}

	var tx *sqlx.Tx
	tx, err = pg.Beginx()
	if err != nil {
		return
	}
	defer tx.Rollback()

	var CreateMessageEPR1Variable ConversationModel
	err = tx.Get(&CreateMessageEPR1Variable, `
		INSERT INTO conversation
		(sender_id, receiver_id, message, sent_date)
		VALUES($1, $2, $3, $4)
		RETURNING *`,
		SenderId1Variable, ReceiverId1Variable, Message2Variable, SentDate1Variable)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(CreateMessageEPR1Variable))
}
