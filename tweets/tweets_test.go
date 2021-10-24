package tweets

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

type InvalidTweetPost struct {
	Msg string `json:"msg"`
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

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(w.Result().Body)

	assert.Equal(t, "User not found", buffer.String())
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

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(w.Result().Body)

	assert.Equal(t, "Invalid payload", buffer.String())
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

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(w.Result().Body)
	tweet_response := buffer.Bytes()

	var tweet_obj models.TweetGet
	json.Unmarshal(tweet_response, &tweet_obj)

	assert.Equal(t, http.StatusCreated, w.Result().StatusCode, "Should return 201")

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

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(w.Result().Body)

	assert.Equal(t, "Tweet not found", buffer.String())
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

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(w.Result().Body)

	assert.Equal(t, "User does not own this tweet", buffer.String())
	assert.Equal(t, http.StatusForbidden, w.Result().StatusCode, "Should return 403")

}

// func TestShouldReturn404SinceTweetDoesNotExist(t *testing.T) {
// 	w := httptest.NewRecorder()

// 	r, err := http.NewRequest(http.MethodDelete, "http://localhost:4000/users/1/tweet/2")
// }
