package main

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
	"utils"
)

// Action20255Endpoint is "get" request handler for "action_20255" endpoint.
func Action20255Endpoint(w http.ResponseWriter, r *http.Request) {

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

	var PartnerSiteId1Variable int

	var paths19556 = mux.Vars(r)
	if paths19556 != nil && len(paths19556) > 0 {
		if val19556, ok := paths19556["partner_site_id"]; ok {
			PartnerSiteId1Variable, err = strconv.Atoi(val19556)
			if err != nil {
				errStatus = http.StatusBadRequest
				errResponse = utils.Response{Error: &utils.ValError{Message: "invalid partner_site_id", Code: utils.ErrCodeBadInput}}
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

	var PartnerSiteListings1Variable []GetListingsEPRConstruct

	err = tx.Select(&PartnerSiteListings1Variable, `
        SELECT a.*, b.*, c.first_name, c.last_name, c.email, c.profile_picture, 
       d.id AS order_type_id, 
       in_time.hour AS in_time_hour, out_time.hour AS out_time_hour,
       f.name AS order_type_name,
       pi.file AS file,
       b.id AS property_id, a.id AS id
FROM listings a 
LEFT JOIN properties b ON a.property_id = b.id 
LEFT JOIN users c ON b.creator_id = c.id 
LEFT JOIN order_types d ON a.order_type_id = d.id 
LEFT JOIN check_in_out_time in_time ON a.check_in_time_id = in_time.id
LEFT JOIN check_in_out_time out_time ON a.check_out_time_id = out_time.id
LEFT JOIN order_types f ON f.id = d.id
LEFT JOIN (
  SELECT DISTINCT ON (property_id) property_id, file
  FROM property_images
) pi ON pi.property_id = b.id
WHERE a.partner_site_id = $1`,
		PartnerSiteId1Variable)

	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(PartnerSiteListings1Variable))
}
