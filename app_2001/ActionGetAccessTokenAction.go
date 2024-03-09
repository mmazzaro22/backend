package app_2001

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// getAccessToken
func GetAccessTokenAction(AccessToken1Variable string, Code2Variable string) (AccessToken2Variable StripeAccessTokenConstruct, err error) {

	defer func(e *error) {
		if e != nil && *e != nil {
			log.Println(*e)
		}
	}(&err)

	var Token1Variable = "Bearer " + AccessToken1Variable

	var ErrorMessage2Variable = "Error getting access token!"

	var urlString31989 = "https://connect.stripe.com/oauth/token"

	var body31989 io.Reader

	var requestData31989 = make(url.Values, 0)

	requestData31989.Add("grant_type", fmt.Sprint(GrantTypeAuthorizationCode1Variable))

	requestData31989.Add("code", fmt.Sprint(Code2Variable))

	body31989 = bytes.NewBufferString(requestData31989.Encode())

	var req31989 *http.Request
	req31989, err = http.NewRequest(strings.ToUpper("post"), urlString31989, body31989)
	if err != nil {
		return
	}

	var reqHeaders31989 = make(http.Header, 0)

	reqHeaders31989.Add("Authorization", fmt.Sprint(Token1Variable))

	reqHeaders31989.Add("Content-Type", fmt.Sprint(GrantTypeAuthorizationCode1Variable))

	req31989.Header = reqHeaders31989

	var resp31989 *http.Response
	resp31989, err = http.DefaultClient.Do(req31989)
	if err != nil {
		return
	} else if resp31989 != nil && resp31989.Body != nil {
		defer resp31989.Body.Close()
	}

	var GetAccessTokenReponseCode1Variable = resp31989.StatusCode

	if AccessToken2Variable.AccessToken != AccessToken2Variable.AccessToken {
		err = &EventError{
			StatusCode: GetAccessTokenReponseCode1Variable,
			Message:    ErrorMessage2Variable,
			Data:       ErrorMessage2Variable,
		}
		return
	}

	return
}
