package app_1995

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// driveListFiles
func DriveListFilesAction(AccessToken2Variable string) (DriveListFilesResponse2Variable GoogleDriveListFilesEPRConstruct, err error) {

	defer func(e *error) {
		if e != nil && *e != nil {
			log.Println(*e)
		}
	}(&err)

	var ErrorMessage3Variable = "Error getting files!"

	var Authorization2Variable = "Bearer " + AccessToken2Variable

	var urlString22873 = "https://www.googleapis.com/drive/v3/files"

	var body22873 io.Reader

	var req22873 *http.Request
	req22873, err = http.NewRequest(strings.ToUpper("get"), urlString22873, body22873)
	if err != nil {
		return
	}

	var reqHeaders22873 = make(http.Header, 0)

	reqHeaders22873.Add("Accept", fmt.Sprint(ApplicationJson1Variable))

	reqHeaders22873.Add("Authorization", fmt.Sprint(Authorization2Variable))

	reqHeaders22873.Add("Content-Type", fmt.Sprint(ApplicationJson1Variable))

	req22873.Header = reqHeaders22873

	var resp22873 *http.Response
	resp22873, err = http.DefaultClient.Do(req22873)
	if err != nil {
		return
	} else if resp22873 != nil && resp22873.Body != nil {
		defer resp22873.Body.Close()
	}

	var Response5Variable GoogleDriveListFilesEPRConstruct

	dec22873 := json.NewDecoder(resp22873.Body)
	if err = dec22873.Decode(&Response5Variable); err != nil {
		return
	}
	DriveListFilesResponse1Variable := Response5Variable

	var DriveListFilesResponseCode1Variable = resp22873.StatusCode

	DriveListFilesResponse2Variable = DriveListFilesResponse1Variable

	if len(DriveListFilesResponse2Variable.Kind) == NumberZero1Variable || DriveListFilesResponseCode1Variable != HttpSuccessCode1Variable {
		err = &EventError{
			StatusCode: DriveListFilesResponseCode1Variable,
			Message:    ErrorMessage3Variable,
			Data:       ErrorMessage3Variable,
		}
		return
	}

	return
}
