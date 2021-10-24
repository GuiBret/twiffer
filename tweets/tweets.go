package tweets

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"twitter-ripoff/dbmgmt"
	"twitter-ripoff/errhandling"
	"twitter-ripoff/models"

	"github.com/gorilla/mux"
)

func main() {

}

func GetTweet(w http.ResponseWriter, r *http.Request) {

	var tweet models.Tweet

	db, err := dbmgmt.GetDBInstance()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Print("DB error")
		return
	}

	id_tweet, err := strconv.Atoi(mux.Vars(r)["idtweet"])

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Print("Invalid tweet ID")
		return
	}

	err = db.First(&tweet, id_tweet).Error

	if err != nil {
		errhandling.HandleError(w, "Tweet not found", http.StatusNotFound, 0)
		return
	}

	res, err := json.Marshal(&tweet)

	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Cant parse tweet")
		return
	}

	fmt.Fprint(w, string(res))

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
		errhandling.HandleError(w, "User not found", http.StatusNotFound, 0)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&payload)

	if err != nil || payload.Message == "" {
		errhandling.HandleError(w, "Invalid payload", http.StatusBadRequest, 0)
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
		errhandling.HandleError(w, "Invalid number", http.StatusBadRequest, 0)
		return
	}

	db, err := dbmgmt.GetDBInstance()

	if err != nil {
		errhandling.HandleError(w, "DB error", http.StatusInternalServerError, 0)
		return
	}

	err = db.First(&tweet, id_tweet).Error

	if err != nil {
		errhandling.HandleError(w, "Tweet not found", http.StatusNotFound, 0)
		return
	}

	if id_user != tweet.FromID {
		errhandling.HandleError(w, "User does not own this tweet", http.StatusForbidden, 0)
		return

	}

}
