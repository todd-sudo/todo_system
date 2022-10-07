package dto

type UpdateUserDTO struct {
	Username  string `json:"username" form:"username"`
	FirstName string `json:"first_name" form:"first_name"`
	LastName  string `json:"last_name" form:"last_name"`
	Avatar    string `json:"avatar" form:"avatar"`
}

type InsertUserDTO struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type UsernameUserDTO struct {
	Username string `json:"username" form:"username"`
}
