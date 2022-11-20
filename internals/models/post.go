package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	UserId    int
	User      *User       `json:"user"`
	Message   string      `json:"message"`
	PostLike  []PostLike  `json:"post_like"`
	PostShare []PostShare `json:"post_share"`
}

type PostLike struct {
	gorm.Model
	PostId int
	UserId int
	User   *User `json:"user"`
}

type PostShare struct {
	gorm.Model
	PostId int
	Post   *Post `json:"post"`
	UserId int
	User   *User `json:"user"`
}
