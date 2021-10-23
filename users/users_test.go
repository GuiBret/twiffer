package users

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"twitter-ripoff/models"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type InvalidUserPost struct {
	Tag_name     string `json:"tag_name"`
	Profile_name string `json:"profile_name"`
}

// TODO: error handling
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

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(mock_response_writer.Result().Body)

	assert.Equal(t, "Error", buffer.String())

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

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(mock_response_writer.Result().Body)
	user_response := buffer.Bytes()

	var user_obj models.UserGet
	json.Unmarshal(user_response, &user_obj)

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

	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode, "Status code should be 404")

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

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(w.Result().Body)
	user_response := buffer.Bytes()

	var user_obj models.UserGet
	json.Unmarshal(user_response, &user_obj)

	assert.Equal(t, uint(1), user_obj.ID, "The users ID should be 1")
	assert.Equal(t, "Machin", user_obj.Tagname, "The users Tagname should be machin")
	assert.Equal(t, "bidule", user_obj.Profile_name, "The users Profile name should be bidule")

}
