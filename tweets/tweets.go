package tweets

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"twitter-ripoff/dbmgmt"
	"twitter-ripoff/models"

	"github.com/gorilla/mux"
)

func main() {

}

func GetTweet(w http.ResponseWriter, r *http.Request) {
}

func WriteTweet(w http.ResponseWriter, r *http.Request) {
	var user models.User
	var payload models.TweetPost
	tweet := models.Tweet{}

	user_id, err := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "parsing error")
		return
	}

	db, err := dbmgmt.GetDBInstance()
	err = db.First(&user, user_id).Error

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "User not found")
		return
	}

	err = json.NewDecoder(r.Body).Decode(&payload)

	if err != nil || payload.Message == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid payload")
		return
	}

	tweet.Content = payload.Message
	tweet.FromID = int(user.ID)
	err = db.Save(&tweet).Error

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Print(err)
		return
	}

	res, err := json.Marshal(&tweet)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(res))
}

// TODO: finish this endpoint
func DeleteTweet(w http.ResponseWriter, r *http.Request) {

	var tweet models.Tweet
	id_user, err := strconv.Atoi(mux.Vars(r)["id"])
	id_tweet, err := strconv.Atoi(mux.Vars(r)["idtweet"])

	if err != nil {
		fmt.Print("Invalid number")
		return
	}

	db, err := dbmgmt.GetDBInstance()

	if err != nil {
		fmt.Print("DB error")
		return
	}

	err = db.First(&tweet, id_tweet).Error

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Tweet not found")
		return
	}

	if id_user != tweet.FromID {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "User does not own this tweet")

	}

}
