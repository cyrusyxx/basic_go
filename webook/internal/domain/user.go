package domain

type User struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`

	NickName    string `json:"nick_name"`
	Birthday    string `json:"birthday"`
	Description string `json:"description"`

	WechatInfo
}
