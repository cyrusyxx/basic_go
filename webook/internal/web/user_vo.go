package web

type UserVo struct {
	Id    int64  `json:"id"`
	Email string `json:"email"`
	Phone string `json:"phone"`

	NickName    string `json:"nickname"`
	Birthday    string `json:"birthday"`
	Description string `json:"description"`
}
