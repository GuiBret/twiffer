package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"twitter-ripoff/errhandling"
	"twitter-ripoff/models"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type InvalidUserPost struct {
	Tag_name     string `json:"tag_name"`
	Profile_name string `json:"profile_name"`
}

func TestShouldReturn400SincePayloadInvalid(t *testing.T) {

	mock_response_writer := httptest.NewRecorder()

	user, err := json.Marshal(&InvalidUserPost{"abc", "def"})

	if err != nil {
		t.Fatal("Oops")
	}

	request, err := http.NewRequest(http.MethodPost, "localhost:4000", bytes.NewReader(user))

	request.Header.Set("Content-Type", "application/json")

	CreateUser(mock_response_writer, request)

	assert.Equal(t, http.StatusBadRequest, mock_response_writer.Result().StatusCode, "Status code should be 400")

	err_response := errhandling.ParseError(mock_response_writer.Result().Body)

	assert.Equal(t, "Invalid payload", err_response.Message)
	assert.Equal(t, 0, err_response.Code)
	assert.Equal(t, http.StatusBadRequest, mock_response_writer.Result().StatusCode)

}

func TestShouldReturn201SincePayloadValid(t *testing.T) {

	mock_response_writer := httptest.NewRecorder()

	user, err := json.Marshal(&models.UserPost{"abc", "def"})

	if err != nil {
		t.Fatal("Oops")
	}

	request, err := http.NewRequest(http.MethodPost, "localhost:4000", bytes.NewReader(user))

	request.Header.Set("Content-Type", "application/json")

	CreateUser(mock_response_writer, request)

	assert.Equal(t, http.StatusCreated, mock_response_writer.Result().StatusCode, "Status code should be 201")

	var user_obj models.UserGet

	ParseResponse(mock_response_writer.Result().Body, &user_obj)

	assert.Equal(t, "abc", user_obj.Tagname, "Tagname should equal abc")
	assert.Equal(t, "def", user_obj.Profile_name, "Profile name should equal def")
}

// TODO: créer DB de test à lancer
func TestShouldReturn404SinceUserDoesNotExist(t *testing.T) {
	w := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "http://localhost:4000/users/400000", nil)

	vars := map[string]string{
		"id": "400000",
	}
	r = mux.SetURLVars(r, vars)
	if err != nil {
		t.Fatal("Error")
	}

	GetOneUser(w, r)

	err_response := errhandling.ParseError(w.Result().Body)

	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode, "Status code should be 404")
	assert.Equal(t, "User not found", err_response.Message, "Status code should be 404")
	assert.Equal(t, 0, err_response.Code, "Status code should be 404")

}

func TestShouldReturnUserSinceExists(t *testing.T) {
	w := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "http://localhost:4000/users/400000", nil)

	vars := map[string]string{
		"id": "1",
	}
	r = mux.SetURLVars(r, vars)
	if err != nil {
		t.Fatal("Error")
	}

	GetOneUser(w, r)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode, "Status code should be 200")

	var user_obj models.UserGet

	ParseResponse(w.Result().Body, &user_obj)

	assert.Equal(t, uint(1), user_obj.ID, "The users ID should be 1")
	assert.Equal(t, "Machin", user_obj.Tagname, "The users Tagname should be machin")
	assert.Equal(t, "bidule", user_obj.Profile_name, "The users Profile name should be bidule")

}

// TODO: mutualiser ces fonctions
func ParseResponse(payload io.ReadCloser, obj interface{}) {
	buffer := new(bytes.Buffer)
	buffer.ReadFrom(payload)
	user_response := buffer.Bytes()

	err := json.Unmarshal(user_response, obj)

	if err != nil {
		fmt.Print(err)
	}

	fmt.Print(obj)

}
