package domain

type User struct {
	Id       int64
	Email    string
	Password string
	Phone    string

	NickName    string
	Birthday    string
	Description string

	WechatInfo
}
