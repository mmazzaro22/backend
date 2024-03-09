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

// createSubscription
func CreateSubscriptionAction(CancelURL2Variable string, ClientSecret1Variable string, PriceID2Variable string, Quantity5Variable int, SuccessURL2Variable string) (Id2Variable string, Url2Variable string, err error) {

	defer func(e *error) {
		if e != nil && *e != nil {
			log.Println(*e)
		}
	}(&err)

	var ModeSubscription1Variable = "subscription"

	var ErrorMessage3Variable = "Error creating subscription!"

	var urlString33464 = "https://api.stripe.com/v1/checkout/sessions"

	var body33464 io.Reader

	var requestData33464 = make(url.Values, 0)

	requestData33464.Add("success_url", fmt.Sprint(SuccessURL2Variable))

	requestData33464.Add("cancel_url", fmt.Sprint(CancelURL2Variable))

	requestData33464.Add("line_items[0][price]", fmt.Sprint(PriceID2Variable))

	requestData33464.Add("line_items[0][quantity]", fmt.Sprint(Quantity5Variable))

	requestData33464.Add("mode", fmt.Sprint(ModeSubscription1Variable))

	body33464 = bytes.NewBufferString(requestData33464.Encode())

	var req33464 *http.Request
	req33464, err = http.NewRequest(strings.ToUpper("post"), urlString33464, body33464)
	if err != nil {
		return
	}

	var reqHeaders33464 = make(http.Header, 0)

	reqHeaders33464.Add("Authorization", fmt.Sprint(ClientSecret1Variable))

	req33464.Header = reqHeaders33464

	var resp33464 *http.Response
	resp33464, err = http.DefaultClient.Do(req33464)
	if err != nil {
		return
	} else if resp33464 != nil && resp33464.Body != nil {
		defer resp33464.Body.Close()
	}

	var Response3Variable StripeCheckoutObjectConstruct

	dec33464 := json.NewDecoder(resp33464.Body)
	if err = dec33464.Decode(&Response3Variable); err != nil {
		return
	}
	CreateSubscriptionCheckoutResponse1Variable := Response3Variable

	var CreateSubscriptionCheckoutResponseCode1Variable = resp33464.StatusCode

	Url2Variable = CreateSubscriptionCheckoutResponse1Variable.Url

	Id2Variable = CreateSubscriptionCheckoutResponse1Variable.Id

	if len(Id2Variable) == NumberZero1Variable || len(Url2Variable) == NumberZero1Variable {
		err = &EventError{
			StatusCode: CreateSubscriptionCheckoutResponseCode1Variable,
			Message:    ErrorMessage3Variable,
			Data:       ErrorMessage3Variable,
		}
		return
	}

	return
}
