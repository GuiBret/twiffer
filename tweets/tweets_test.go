package tweets

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

type InvalidTweetPost struct {
	Msg string `json:"msg"`
}

func TestGetTweetShouldReturn404SinceTweetDoesNotExist(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost:4000/users/100000/tweet/2", nil)

	if err != nil {
		t.Fatal("Error")
	}

	vars := map[string]string{
		"id":      "1",
		"idtweet": "2000000",
	}
	r = mux.SetURLVars(r, vars)

	GetTweet(w, r)

	err_response := errhandling.ParseError(w.Result().Body)

	assert.Equal(t, "Tweet not found", err_response.Message)
	assert.Equal(t, 0, err_response.Code, "Should have a code 0")
	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode, "Should return 404")
}

// TODO: handle payload
func TestGetTweetShouldReturn200SinceTweetExists(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost:4000/users/100000/tweet/2", nil)

	if err != nil {
		t.Fatal("Error")
	}

	vars := map[string]string{
		"id":      "1",
		"idtweet": "5",
	}
	r = mux.SetURLVars(r, vars)

	GetTweet(w, r)

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(w.Result().Body)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode, "Should return 200")
}

func TestShouldReturn404SinceUserDoesNotExist(t *testing.T) {
	w := httptest.NewRecorder()

	tweet, err := json.Marshal(models.TweetPost{"This is my message"})

	if err != nil {
		t.Fatal("Error")
	}
	r, err := http.NewRequest(http.MethodPost, "http://localhost:4000/users/100000/tweet", bytes.NewReader(tweet))

	vars := map[string]string{
		"id": "400000",
	}
	r = mux.SetURLVars(r, vars)

	WriteTweet(w, r)

	err_response := errhandling.ParseError(w.Result().Body)

	assert.Equal(t, "User not found", err_response.Message)
	assert.Equal(t, 0, err_response.Code)
	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode, "Should return 404")

}

func TestShouldReturn400SinceInvalidPayload(t *testing.T) {
	w := httptest.NewRecorder()

	tweet, err := json.Marshal(InvalidTweetPost{"This is my message"})

	if err != nil {
		t.Fatal("Error")
	}
	r, err := http.NewRequest(http.MethodPost, "http://localhost:4000/users/100000/tweet", bytes.NewReader(tweet))

	vars := map[string]string{
		"id": "1",
	}
	r = mux.SetURLVars(r, vars)

	WriteTweet(w, r)

	err_response := errhandling.ParseError(w.Result().Body)

	assert.Equal(t, "Invalid payload", err_response.Message)
	assert.Equal(t, 0, err_response.Code)
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode, "Should return 400")

}

func TestShouldCreateTweet(t *testing.T) {

	w := httptest.NewRecorder()

	tweet, err := json.Marshal(models.TweetPost{"This is my message"})

	if err != nil {
		t.Fatal("Error")
	}
	r, err := http.NewRequest(http.MethodPost, "http://localhost:4000/users/100000/tweet", bytes.NewReader(tweet))

	vars := map[string]string{
		"id": "1",
	}
	r = mux.SetURLVars(r, vars)

	WriteTweet(w, r)

	var tweet_obj models.TweetGet

	ParseResponse(w.Result().Body, &tweet_obj)

	assert.Equal(t, http.StatusCreated, w.Result().StatusCode, "Should return 201")
	// assert.Equal(t, "This is my message", tweet_obj.Message, "This is my message")

	// TODO: handle payload later
	// assert.Equal(t, "This is my message", tweet_obj.Message)

}

func TestDeleteUserShouldReturn404SinceTweetDoesNotExist(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodDelete, "http://localhost:4000/users/100000/tweet/2", nil)

	if err != nil {
		t.Fatal("Error")
	}

	vars := map[string]string{
		"id":      "1",
		"idtweet": "2000000",
	}
	r = mux.SetURLVars(r, vars)

	DeleteTweet(w, r)

	err_response := errhandling.ParseError(w.Result().Body)

	assert.Equal(t, "Tweet not found", err_response.Message)
	assert.Equal(t, 0, err_response.Code)
	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode, "Should return 404")

}

func TestDeleteUserShouldReturn401SinceTweetDoesNotBelongToUser(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodDelete, "http://localhost:4000/users/100000/tweet/2", nil)

	if err != nil {
		t.Fatal("Error")
	}

	vars := map[string]string{
		"id":      "2",
		"idtweet": "5",
	}
	r = mux.SetURLVars(r, vars)

	DeleteTweet(w, r)

	err_response := errhandling.ParseError(w.Result().Body)

	assert.Equal(t, "User does not own this tweet", err_response.Message)
	assert.Equal(t, http.StatusForbidden, w.Result().StatusCode, "Should return 403")

}

func ParseResponse(payload io.ReadCloser, obj interface{}) {
	buffer := new(bytes.Buffer)
	buffer.ReadFrom(payload)
	user_response := buffer.Bytes()

	err := json.Unmarshal(user_response, obj)

	if err != nil {
		fmt.Print(err)
	}
}
