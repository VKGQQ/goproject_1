package message

// User 定义一个用户的结构体
type User struct {
	UserId     int    `json:"userId"`
	UserPwd    string `json:"userPwd"`
	UserName   string `json:"userName"`
	UserStatus int    `json:"userStatus"`
	Sex        string `json:"sex"`
}
