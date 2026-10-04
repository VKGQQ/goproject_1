package process

import (
	"fmt"
	"net"
	"os"

	"github.com/VKGQQ/goproject_1/chat/client/utils"
)

func ShowMenu() {
	for {
		fmt.Println("----------恭喜登录成功----------")
		fmt.Println("----------1.显示在线用户列表-----")
		fmt.Println("----------2.发送信息----------")
		fmt.Println("----------3.信息列表----------")
		fmt.Println("----------4.退出系统----------")
		fmt.Println("请选择(1-4):")
		var key int
		fmt.Scanln(&key)
		switch key {
		case 1:
			fmt.Println("显示在线用户列表")
		case 2:
			fmt.Println("发送信息")
		case 3:
			fmt.Println("信息列表")
		case 4:
			fmt.Println("您选择提出系统")
			os.Exit(0)
		default:
			fmt.Println("您的输入有误请重新输入")
		}
	}
}

func serverProcessMes(Conn net.Conn) {
	tf := &utils.Transfer{
		Conn: Conn,
	}
	for {
		fmt.Println("客户端正在等待服务器发送的信息......")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("tf.ReadPkg err:", err)
			return
		}
		fmt.Println("mes:", mes)
	}
}
