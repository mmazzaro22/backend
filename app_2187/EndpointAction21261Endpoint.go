package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"utils"
)

// Action21261Endpoint is "get" request handler for "action_21261" endpoint.
func Action21261Endpoint(w http.ResponseWriter, r *http.Request) {

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

	var Price1Variable = 100

	var Quantity4Variable = 1

	var Currency3Variable = "eur"

	var Name4Variable = "name"

	var urlString41798 = "https://api.stripe.com/v1/checkout/sessions"

	var body41798 io.Reader

	var requestData41798 = make(url.Values, 0)

	requestData41798.Add("mode", fmt.Sprint(StripePaymentMode1Variable))

	requestData41798.Add("line_items[0][price_data][unit_amount]", fmt.Sprint(Price1Variable))

	requestData41798.Add("line_items[0][quantity]", fmt.Sprint(Quantity4Variable))

	requestData41798.Add("line_items[0][price_data][currency]", fmt.Sprint(Currency3Variable))

	requestData41798.Add("line_items[0][price_data][product_data][name]", fmt.Sprint(Name4Variable))

	requestData41798.Add("success_url", fmt.Sprint(ApplicationUrl1Variable))

	body41798 = bytes.NewBufferString(requestData41798.Encode())

	var req41798 *http.Request
	req41798, err = http.NewRequest(strings.ToUpper("post"), urlString41798, body41798)
	if err != nil {
		return
	}

	var reqHeaders41798 = make(http.Header, 0)

	reqHeaders41798.Add("Authorization", fmt.Sprint(StripeSecretKey1Variable))

	reqHeaders41798.Add("Content-Type", fmt.Sprint(ApplicationFormUrlEncoded1Variable))

	req41798.Header = reqHeaders41798

	var resp41798 *http.Response
	resp41798, err = http.DefaultClient.Do(req41798)
	if err != nil {
		return
	} else if resp41798 != nil && resp41798.Body != nil {
		defer resp41798.Body.Close()
	}

	err = removeLoginSession(&w, r)
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(nil))
}
