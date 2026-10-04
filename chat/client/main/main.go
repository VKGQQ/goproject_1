package main

import (
	"fmt"
	"os"

	"github.com/VKGQQ/goproject_1/chat/client/process"
)

var userId int
var userPwd string

func main() {
	var key int
	for {
		fmt.Println("----------欢迎登录多人聊天系统----------")
		fmt.Println("\t\t\t 1 登录聊天室")
		fmt.Println("\t\t\t 2 注册用户")
		fmt.Println("\t\t\t 3 退出系统")
		fmt.Println("\t\t\t 请选择(1-3):")
		fmt.Scanln(&key)
		switch key {
		case 1:
			fmt.Println("请输入用户的id")
			fmt.Scanln(&userId)
			fmt.Println("请输入用户的密码")
			fmt.Scanln(&userPwd)
			up := &process.UserProcess{}
			err := up.Login(userId, userPwd)
			if err != nil {
				return
			}
		case 2:
			fmt.Println("注册用户")
		case 3:
			os.Exit(0)
		default:
			fmt.Println("您的输入有误请重新输入")
		}
	}
}
