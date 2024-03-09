package main

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
	"utils"
)

// Action21415Endpoint is "get" request handler for "action_21415" endpoint.
func Action21415Endpoint(w http.ResponseWriter, r *http.Request) {

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

	var ReceiverId3Variable int

	var SenderId2Variable int

	var paths20692 = mux.Vars(r)
	if paths20692 != nil && len(paths20692) > 0 {
		if val20692, ok := paths20692["receiver_id"]; ok {
			ReceiverId3Variable, err = strconv.Atoi(val20692)
			if err != nil {
				errStatus = http.StatusBadRequest
				errResponse = utils.Response{Error: &utils.ValError{Message: "invalid receiver_id", Code: utils.ErrCodeBadInput}}
				return
			}
		}
		if val20692, ok := paths20692["sender_id"]; ok {
			SenderId2Variable, err = strconv.Atoi(val20692)
			if err != nil {
				errStatus = http.StatusBadRequest
				errResponse = utils.Response{Error: &utils.ValError{Message: "invalid sender_id", Code: utils.ErrCodeBadInput}}
				return
			}
		}
	}

	var tx *sqlx.Tx
	tx, err = pg.Beginx()
	if err != nil {
		return
	}
	defer tx.Rollback()

	var GetConversation1Variable []GetMessageEPRConstruct

	err = tx.Select(&GetConversation1Variable, `
        SELECT a.*, 
       b.id, 
       b.email, 
       b.first_name AS sender_first_name, 
       b.last_name AS sender_last_name, 
       pi1.file AS sender_profile_picture,
       b2.first_name AS receiver_first_name, 
       b2.last_name AS receiver_last_name, 
       pi2.file AS receiver_profile_picture,
       CASE WHEN a.sender_id = $1 THEN a.message ELSE NULL END AS sender_message, 
       CASE WHEN a.receiver_id = $1 THEN a.message ELSE NULL END AS receiver_message 
FROM conversation a 
LEFT JOIN users b ON a.sender_id = b.id 
LEFT JOIN users b2 ON a.receiver_id = b2.id
LEFT JOIN property_images pi1 ON a.sender_id = pi1.user_id
LEFT JOIN property_images pi2 ON a.receiver_id = pi2.user_id
WHERE (a.sender_id = $1 AND a.receiver_id = $2) OR (a.sender_id = $2 AND a.receiver_id = $1) ORDER BY a.sent_date DESC;`,
		ReceiverId3Variable, SenderId2Variable)

	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(GetConversation1Variable))
}
