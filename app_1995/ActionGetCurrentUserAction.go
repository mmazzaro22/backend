package app_1995

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// getCurrentUser
func GetCurrentUserAction(AccessToken1Variable string) (Email1Variable string, err error) {

	defer func(e *error) {
		if e != nil && *e != nil {
			log.Println(*e)
		}
	}(&err)

	var ErrorMessage2Variable = "Error getting user!"

	var Authorization1Variable = "Bearer " + AccessToken1Variable

	var urlString22803 = "https://www.googleapis.com/oauth2/v1/userinfo"

	var body22803 io.Reader

	var req22803 *http.Request
	req22803, err = http.NewRequest(strings.ToUpper("get"), urlString22803, body22803)
	if err != nil {
		return
	}

	var reqHeaders22803 = make(http.Header, 0)

	reqHeaders22803.Add("Accept", fmt.Sprint(ApplicationJson1Variable))

	reqHeaders22803.Add("Authorization", fmt.Sprint(Authorization1Variable))

	reqHeaders22803.Add("Content-Type", fmt.Sprint(ApplicationJson1Variable))

	req22803.Header = reqHeaders22803

	var resp22803 *http.Response
	resp22803, err = http.DefaultClient.Do(req22803)
	if err != nil {
		return
	} else if resp22803 != nil && resp22803.Body != nil {
		defer resp22803.Body.Close()
	}

	var Response2Variable GoogleGetCurrentUserConstruct

	dec22803 := json.NewDecoder(resp22803.Body)
	if err = dec22803.Decode(&Response2Variable); err != nil {
		return
	}
	GetCurrentUserResponse1Variable := Response2Variable

	var GetCurrentUserResponseCode1Variable = resp22803.StatusCode

	Email1Variable = GetCurrentUserResponse1Variable.Email

	if GetCurrentUserResponseCode1Variable != HttpSuccessCode1Variable || len(Email1Variable) == NumberZero1Variable {
		err = &EventError{
			StatusCode: HttpSuccessCode1Variable,
			Message:    ErrorMessage2Variable,
			Data:       ErrorMessage2Variable,
		}
		return
	}

	return
}
