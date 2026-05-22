package users

type BaseField struct {
	UserName string `form:"user_name" json:"user_name"  binding:"required,min=1"`
	Pass     string `form:"pass" json:"pass" binding:"required,min=6,max=20"`
}

type Id struct {
	Id float64 `form:"id"  json:"id" binding:"required,min=1"`
}
