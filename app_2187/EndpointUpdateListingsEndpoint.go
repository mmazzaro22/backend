package main

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
	"utils"
)

// Update listings
func UpdateListingsEndpoint(w http.ResponseWriter, r *http.Request) {

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

	var Id16Variable int

	var paths19482 = mux.Vars(r)
	if paths19482 != nil && len(paths19482) > 0 {
		if val19482, ok := paths19482["id"]; ok {
			Id16Variable, err = strconv.Atoi(val19482)
			if err != nil {
				errStatus = http.StatusBadRequest
				errResponse = utils.Response{Error: &utils.ValError{Message: "invalid id", Code: utils.ErrCodeBadInput}}
				return
			}
		}
	}

	var validator19482 *validator.Validate
	validator19482 = validator.New()

	var UpdateListingEPI1Variable UpdateListingsEPIConstruct

	var dec19482 = json.NewDecoder(r.Body)
	if err = dec19482.Decode(&UpdateListingEPI1Variable); err != nil {
		errStatus = http.StatusBadRequest
		errResponse = utils.InvalidJSONErrorResponse(err)
		return
	}
	err = validator19482.Struct(UpdateListingEPI1Variable)
	var ve19482 validator.ValidationErrors
	if errors.As(err, &ve19482) {
		errStatus = http.StatusBadRequest
		errResponse = utils.Response{Error: &utils.ValError{Param: ve19482[0].Field(), Message: MsgToTag(ve19482[0]), Code: "invalid_json"}}
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

	var ListingToUpdate1Variable ListingsModel

	err = tx.Get(&ListingToUpdate1Variable, `
        SELECT id, guest_limit, property_id, order_type_id, creator_id, price, partner_site_id, check_in_time_id, check_out_time_id
        FROM listings
        WHERE (id = $1 OR $1 IS NULL)`,
		Id16Variable)

	if err != nil {
		return
	}

	var nullableGuestLimit39657 *int

	if UpdateListingEPI1Variable.GuestLimit != nil {
		nullableGuestLimit39657 = new(int)
		*nullableGuestLimit39657 = UpdateListingEPI1Variable.GuestLimit.int
	}

	var nullableOrderTypeId39657 *int

	if UpdateListingEPI1Variable.OrderTypeId != nil {
		nullableOrderTypeId39657 = new(int)
		*nullableOrderTypeId39657 = UpdateListingEPI1Variable.OrderTypeId.int
	}

	var nullablePrice39657 *int

	if UpdateListingEPI1Variable.Price != nil {
		nullablePrice39657 = new(int)
		*nullablePrice39657 = UpdateListingEPI1Variable.Price.int
	}

	_, err = tx.Exec(`
        UPDATE listings
        SET guest_limit = $1, order_type_id = $2, price = $3, check_in_time_id = $4, check_out_time_id = $5
        WHERE id = $6`,
		nullableGuestLimit39657, nullableOrderTypeId39657, nullablePrice39657, UpdateListingEPI1Variable.CheckInId, UpdateListingEPI1Variable.CheckOutId, ListingToUpdate1Variable.Id)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(UpdateListingEPI1Variable))
}
