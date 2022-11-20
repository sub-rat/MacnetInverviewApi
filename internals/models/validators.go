package models

type UserModelValidator struct {
	UserName  string `form:"username" json:"username,omitempty"`
	Email     string `form:"email" json:"email" binding:"required,email"`
	Password  string `form:"password" json:"password" binding:"required,min=8,max=255"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type LoginValidator struct {
	Email    string `form:"email" json:"email" binding:"email,min=8,max=255"`
	UserName string `form:"username" json:"username"`
	Password string `form:"password" json:"password" binding:"required,min=8,max=255"`
}

type PostValidator struct {
	Message string `json:"message" binding:"required"`
}
