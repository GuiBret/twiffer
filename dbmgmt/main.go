package dbmgmt

import (
	"fmt"
	"log"
	"twitter-ripoff/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	InitDB()
}

func GetDBInstance() (*gorm.DB, error) {
	return gorm.Open(mysql.Open("myuser:mypassword@tcp(127.0.0.1:3307)/twiffer"))
}

func InitDB() {
	db, err := GetDBInstance()

	if err != nil {
		fmt.Print(err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Tweet{}, &models.Mention{}, &models.Like{}, &models.Notification{})

	if err != nil {
		log.Fatal(err)
		return
	}

}
