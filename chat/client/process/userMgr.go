package process

import (
	"fmt"

	"github.com/VKGQQ/goproject_1/chat/client/model"
	"github.com/VKGQQ/goproject_1/chat/common/message"
)

var onlineUsers = make(map[int]*message.User, 20)
var CurUser model.CurUser

func outputOnlineUser() {
	fmt.Println("当前在线用户列表：")
	for id := range onlineUsers {
		fmt.Println("用户ID：\t", id)
	}
	fmt.Println()
}

func updateUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes) {
	user, ok := onlineUsers[notifyUserStatusMes.UserId]
	if !ok {
		user = &message.User{
			UserId: notifyUserStatusMes.UserId,
		}
	}
	user.UserStatus = notifyUserStatusMes.Status
	onlineUsers[notifyUserStatusMes.UserId] = user
}
