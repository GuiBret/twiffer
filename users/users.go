package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"twitter-ripoff/dbmgmt"
	"twitter-ripoff/errhandling"
	"twitter-ripoff/models"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {

}

func GetOneUser(w http.ResponseWriter, r *http.Request) {

	db, err := dbmgmt.GetDBInstance()

	user_id, err := strconv.Atoi(mux.Vars(r)["id"])
	user := &models.User{}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = db.Select([]string{"ID", "Tagname", "Profile_name"}).First(user, user_id).Error

	if err != nil {

		errhandling.HandleError(w, "User not found", http.StatusNotFound, 0)

		return
	}

	var user_parsed models.UserGet

	// TODO: gérer ça autrement
	user_parsed.ID = user.ID
	user_parsed.Tagname = user.Tagname
	user_parsed.Profile_name = user.Profile_name

	res, err := json.Marshal(&user_parsed)
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(res))
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	var u models.UserPost
	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil || u.Profile_name == "" || u.Tagname == "" {
		errhandling.HandleError(w, "Invalid payload", http.StatusBadRequest, 0)
		return
	}

	db, err := dbmgmt.GetDBInstance()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user := models.User{}

	user.Tagname = u.Tagname
	user.Profile_name = u.Profile_name
	db.Save(&user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res, err := json.Marshal(&user)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(res))
}
