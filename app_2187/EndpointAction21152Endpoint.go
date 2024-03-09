package main

import (
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"image"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"utils"
)

// Action21152Endpoint is "post" request handler for "action_21152" endpoint.
func Action21152Endpoint(w http.ResponseWriter, r *http.Request) {

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

	var FilePath1Variable string

	var PropertyId7Variable int

	var UserId3Variable int

	if r.URL != nil {
		var queries20435 = r.URL.Query()
		if queries20435 != nil && len(queries20435) > 0 {
			if vals20435, ok := queries20435["filePath"]; ok && len(vals20435) > 0 {
				FilePath1Variable = strings.Join(vals20435, "\n")
			}
			if vals20435, ok := queries20435["property_id"]; ok && len(vals20435) > 0 {
				PropertyId7Variable, err = strconv.Atoi(vals20435[0])
				if err != nil {
					errStatus = http.StatusBadRequest
					errResponse = utils.Response{Error: &utils.ValError{Message: "invalid property_id", Code: utils.ErrCodeBadInput}}
					return
				}
			}
			if vals20435, ok := queries20435["user_id"]; ok && len(vals20435) > 0 {
				UserId3Variable, err = strconv.Atoi(vals20435[0])
				if err != nil {
					errStatus = http.StatusBadRequest
					errResponse = utils.Response{Error: &utils.ValError{Message: "invalid user_id", Code: utils.ErrCodeBadInput}}
					return
				}
			}
		}
	}

	var Input8Variable File

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		return
	}

	var file20435 multipart.File
	file20435, _, err = r.FormFile("input")
	if err != nil {
		return
	}

	Input8Variable, err = fs.SetTempFile(file20435)
	if err != nil {
		return
	}
	defer Input8Variable.Close()

	var tx *sqlx.Tx
	tx, err = pg.Beginx()
	if err != nil {
		return
	}
	defer tx.Rollback()

	var FilePathName1Variable PropertyImagesModel
	err = tx.Get(&FilePathName1Variable, `
		INSERT INTO property_images
		(file_name, property_id, user_id)
		VALUES($1, $2, $3)
		RETURNING *`,
		FilePath1Variable, PropertyId7Variable, UserId3Variable)
	if err != nil {
		return
	}

	if len(FilePath1Variable) == 0 {
		err = fmt.Errorf("missing file path")
		return
	}

	// generate attachment
	var newAttachment41580 = Attachment{Filename: FilePath1Variable}
	var fileStat41580 os.FileInfo
	fileStat41580, err = Input8Variable.Stat()
	if err != nil {
		err = fmt.Errorf("get file stat: %w", err)
		return
	} else {
		newAttachment41580.Size = int(fileStat41580.Size())
	}
	var mime41580 *mimetype.MIME
	mime41580, err = mimetype.DetectReader(Input8Variable)
	if err != nil {
		err = fmt.Errorf("get file type: %w", err)
		return
	} else {
		newAttachment41580.Type = mime41580.String()
	}

	_, err = Input8Variable.Seek(0, io.SeekStart)
	if err != nil {
		err = fmt.Errorf("seek to start of file: %w", err)
		return
	}
	var config41580 image.Config
	config41580, _, err = image.DecodeConfig(Input8Variable)
	if err != nil {
		err = nil
	} else {
		newAttachment41580.Height = &config41580.Height
		newAttachment41580.Width = &config41580.Width
	}

	_, err = Input8Variable.Seek(0, io.SeekStart)
	if err != nil {
		err = fmt.Errorf("seek to start of file: %w", err)
		return
	}

	// get current attachments
	var currentAttachmentsRow41580 PropertyImagesModel

	err = tx.Get(&currentAttachmentsRow41580, `
        SELECT file FROM property_images WHERE id = $1`,
		FilePathName1Variable.Id)

	if err != nil {
		return
	}

	// update current attachments
	var addedAttachment41580 bool
	if currentAttachmentsRow41580.File != nil {
		for i, attachment := range *currentAttachmentsRow41580.File {
			if attachment.Filename == newAttachment41580.Filename {
				err = fs.SetFile(Input8Variable, attachment.Filepath, true)
				if err != nil {
					return
				}

				newAttachment41580.Id = attachment.Id
				newAttachment41580.Filepath = attachment.Filepath
				newAttachment41580.Url = attachment.Url
				(*currentAttachmentsRow41580.File)[i] = newAttachment41580
				addedAttachment41580 = true
			}
		}
	}

	if !addedAttachment41580 {
		// store new attachment
		var attachmentID uuid.UUID
		attachmentID, err = uuid.NewV4()
		if err != nil {
			err = fmt.Errorf("generate attachment id: %w", err)
			return
		} else {
			newAttachment41580.Id = attachmentID.String()
		}

		var dir, filename string
		dir, filename = filepath.Split(newAttachment41580.Filename)
		newAttachment41580.Filepath = filepath.Join(dir, fmt.Sprintf("%s-%s", newAttachment41580.Id, filename))
		newAttachment41580.Url = filepath.Join(UploadsBaseURL, newAttachment41580.Filepath)

		err = fs.SetFile(Input8Variable, newAttachment41580.Filepath, true)
		if err != nil {
			return
		}
		if currentAttachmentsRow41580.File == nil {
			currentAttachmentsRow41580.File = new(AttachmentArray)
		}

		*currentAttachmentsRow41580.File = append(*currentAttachmentsRow41580.File, newAttachment41580)
	}

	_, err = tx.Exec(`
        UPDATE property_images
        SET file = $1
        WHERE id = $2`,
		currentAttachmentsRow41580.File, FilePathName1Variable.Id)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.JSON(w, http.StatusOK, utils.OKResponse(FilePathName1Variable))
}
