package app_1995

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// getAccessToken
func GetAccessTokenAction(ClientId1Variable string, ClientSecret1Variable string, Code1Variable string, RedirectUri1Variable string) (AccessToken1Variable string, err error) {

	defer func(e *error) {
		if e != nil && *e != nil {
			log.Println(*e)
		}
	}(&err)

	var ErrorMessage1Variable = "Error getting access token!"

	var GrantType1Variable = "authorization_code"

	var urlString22787 = "https://www.googleapis.com/oauth2/v4/token"

	var body22787 io.Reader

	var requestData22787 = make(url.Values, 0)

	requestData22787.Add("client_id", fmt.Sprint(ClientId1Variable))

	requestData22787.Add("client_secret", fmt.Sprint(ClientSecret1Variable))

	requestData22787.Add("code", fmt.Sprint(Code1Variable))

	requestData22787.Add("grant_type", fmt.Sprint(GrantType1Variable))

	requestData22787.Add("redirect_uri", fmt.Sprint(RedirectUri1Variable))

	body22787 = bytes.NewBufferString(requestData22787.Encode())

	var req22787 *http.Request
	req22787, err = http.NewRequest(strings.ToUpper("post"), urlString22787, body22787)
	if err != nil {
		return
	}

	var reqHeaders22787 = make(http.Header, 0)

	reqHeaders22787.Add("Accept", fmt.Sprint(ApplicationJson1Variable))

	reqHeaders22787.Add("Content-Type", fmt.Sprint(ApplicationFormUrlEncoded1Variable))

	req22787.Header = reqHeaders22787

	var resp22787 *http.Response
	resp22787, err = http.DefaultClient.Do(req22787)
	if err != nil {
		return
	} else if resp22787 != nil && resp22787.Body != nil {
		defer resp22787.Body.Close()
	}

	var Response1Variable GoogleGetAccessTokenEPRConstruct

	dec22787 := json.NewDecoder(resp22787.Body)
	if err = dec22787.Decode(&Response1Variable); err != nil {
		return
	}
	GetAccessTokenReponse1Variable := Response1Variable

	var GetAccessTokenReponseCode1Variable = resp22787.StatusCode

	AccessToken1Variable = GetAccessTokenReponse1Variable.AccessToken

	if len(AccessToken1Variable) == NumberZero1Variable || GetAccessTokenReponseCode1Variable != HttpSuccessCode1Variable {
		err = &EventError{
			StatusCode: GetAccessTokenReponseCode1Variable,
			Message:    ErrorMessage1Variable,
			Data:       ErrorMessage1Variable,
		}
		return
	}

	return
}
