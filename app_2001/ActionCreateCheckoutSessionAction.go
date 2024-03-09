package app_2001

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// createCheckoutSession
func CreateCheckoutSessionAction(CancelURL1Variable string, Currency2Variable string, Price1Variable int, ProductName1Variable string, Quantity2Variable int, SecretKey1Variable string, SuccessURL1Variable string) (Id1Variable string, Url1Variable string, err error) {

	defer func(e *error) {
		if e != nil && *e != nil {
			log.Println(*e)
		}
	}(&err)

	var ErrorMessage1Variable = "Error creating checkout session!"

	var Mode2Variable = "payment"

	var urlString22823 = "https://api.stripe.com/v1/checkout/sessions"

	var body22823 io.Reader

	var requestData22823 = make(url.Values, 0)

	requestData22823.Add("cancel_url", fmt.Sprint(CancelURL1Variable))

	requestData22823.Add("line_items[0][price_data][currency]", fmt.Sprint(Currency2Variable))

	requestData22823.Add("line_items[0][price_data][unit_amount]", fmt.Sprint(Price1Variable))

	requestData22823.Add("line_items[0][quantity]", fmt.Sprint(Quantity2Variable))

	requestData22823.Add("mode", fmt.Sprint(Mode2Variable))

	requestData22823.Add("line_items[0][price_data][product_data][name]", fmt.Sprint(ProductName1Variable))

	requestData22823.Add("success_url", fmt.Sprint(SuccessURL1Variable))

	body22823 = bytes.NewBufferString(requestData22823.Encode())

	var req22823 *http.Request
	req22823, err = http.NewRequest(strings.ToUpper("post"), urlString22823, body22823)
	if err != nil {
		return
	}

	var reqHeaders22823 = make(http.Header, 0)

	reqHeaders22823.Add("Authorization", fmt.Sprint(SecretKey1Variable))

	req22823.Header = reqHeaders22823

	var resp22823 *http.Response
	resp22823, err = http.DefaultClient.Do(req22823)
	if err != nil {
		return
	} else if resp22823 != nil && resp22823.Body != nil {
		defer resp22823.Body.Close()
	}

	var Response1Variable StripeCheckoutObjectConstruct

	dec22823 := json.NewDecoder(resp22823.Body)
	if err = dec22823.Decode(&Response1Variable); err != nil {
		return
	}
	CreateCheckoutSessionResponse1Variable := Response1Variable

	var CreateCheckoutSessionResponseCode1Variable = resp22823.StatusCode

	Url1Variable = CreateCheckoutSessionResponse1Variable.Url

	Id1Variable = CreateCheckoutSessionResponse1Variable.Id

	if CreateCheckoutSessionResponseCode1Variable != HttpSuccessCode1Variable || len(Url1Variable) == NumberZero1Variable && len(Id1Variable) == NumberZero1Variable {
		err = &EventError{
			StatusCode: CreateCheckoutSessionResponseCode1Variable,
			Message:    ErrorMessage1Variable,
			Data:       ErrorMessage1Variable,
		}
		return
	}

	return
}
