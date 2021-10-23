package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Tagname      string
	Profile_name string
	Tweets       []Tweet `gorm:"foreignKey:From"`
}

type UserGet struct {
	ID           uint   `json:"id"`
	Tagname      string `json:"tagname"`
	Profile_name string `json:"profile_name"`
}

type Tweet struct {
	gorm.Model
	From      User
	FromID    int
	Content   string
	Timestamp time.Time
}

type TweetPost struct {
	Message string `json:"message"`
}

type Mention struct {
	gorm.Model
	Tweet   Tweet
	TweetID int
	User    User
	UserID  int
}

type Like struct {
	gorm.Model
	Tweet     Tweet
	TweetID   int
	User      User
	UserID    int
	Timestamp time.Time
}

type Notification struct {
	gorm.Model
	Type      string
	Timestamp time.Time
}

type UserPost struct {
	Tagname      string `json:"tagname"`
	Profile_name string `json:"profile_name"`
}
