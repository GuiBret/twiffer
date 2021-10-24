package errhandling

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type StandardError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func HandleError(w http.ResponseWriter, message string, status_code int, code int) {

	var err StandardError

	w.WriteHeader(status_code)

	err.Code = code
	err.Message = message

	err_parsed, parsing_error := json.Marshal(&err)

	if parsing_error != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Fprint(w, string(err_parsed))

}

func ParseError(payload io.ReadCloser) StandardError {

	var err_parsed StandardError

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(payload)
	user_error := buffer.Bytes()

	json.Unmarshal(user_error, &err_parsed)

	return err_parsed

}
