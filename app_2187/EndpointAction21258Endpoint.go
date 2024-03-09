package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"utils"
)

// Action21258Endpoint is "get" request handler for "action_21258" endpoint.
func Action21258Endpoint(w http.ResponseWriter, r *http.Request) {

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

	var Id17Variable int

	var paths20536 = mux.Vars(r)
	if paths20536 != nil && len(paths20536) > 0 {
		if val20536, ok := paths20536["id"]; ok {
			Id17Variable, err = strconv.Atoi(val20536)
			if err != nil {
				errStatus = http.StatusBadRequest
				errResponse = utils.Response{Error: &utils.ValError{Message: "invalid id", Code: utils.ErrCodeBadInput}}
				return
			}
		}
	}

	var Quantity2Variable int

	if r.URL != nil {
		var queries20536 = r.URL.Query()
		if queries20536 != nil && len(queries20536) > 0 {
			if vals20536, ok := queries20536["quantity"]; ok && len(vals20536) > 0 {
				Quantity2Variable, err = strconv.Atoi(vals20536[0])
				if err != nil {
					errStatus = http.StatusBadRequest
					errResponse = utils.Response{Error: &utils.ValError{Message: "invalid quantity", Code: utils.ErrCodeBadInput}}
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

	var Listing4Variable ListingsModel

	err = tx.Get(&Listing4Variable, `
        SELECT id, guest_limit, property_id, order_type_id, creator_id, price, partner_site_id, check_in_time_id, check_out_time_id
        FROM listings
        WHERE (id = $1 OR $1 IS NULL)`,
		Id17Variable)

	if err != nil {
		return
	}

	var Name2Variable = "purchase"

	var urlString41795 = "https://api.stripe.com/v1/checkout/sessions"

	var body41795 io.Reader

	var requestData41795 = make(url.Values, 0)

	requestData41795.Add("line_items[0][quantity]", fmt.Sprint(Quantity2Variable))

	requestData41795.Add("line_items[0][price_data][currency]", fmt.Sprint(Currency1Variable))

	requestData41795.Add("line_items[0][price_data][product_data][name]", fmt.Sprint(Name2Variable))

	requestData41795.Add("success_url", fmt.Sprint(ApplicationUrl1Variable))

	requestData41795.Add("mode", fmt.Sprint(StripePaymentMode1Variable))

	requestData41795.Add("line_items[0][price_data][unit_amount]", fmt.Sprint(Listing4Variable.Price))

	body41795 = bytes.NewBufferString(requestData41795.Encode())

	var req41795 *http.Request
	req41795, err = http.NewRequest(strings.ToUpper("post"), urlString41795, body41795)
	if err != nil {
		return
	}

	var reqHeaders41795 = make(http.Header, 0)

	reqHeaders41795.Add("Authorization", fmt.Sprint(StripeSecretKey1Variable))

	reqHeaders41795.Add("Content-Type", fmt.Sprint(ApplicationFormUrlEncoded1Variable))

	req41795.Header = reqHeaders41795

	var resp41795 *http.Response
	resp41795, err = http.DefaultClient.Do(req41795)
	if err != nil {
		return
	} else if resp41795 != nil && resp41795.Body != nil {
		defer resp41795.Body.Close()
	}

	var StripeCheckoutObject1Variable string

	dec41795 := json.NewDecoder(resp41795.Body)
	if err = dec41795.Decode(&StripeCheckoutObject1Variable); err != nil {
		return
	}
	CreateCheckoutResponse1Variable := StripeCheckoutObject1Variable

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(CreateCheckoutResponse1Variable))
}
