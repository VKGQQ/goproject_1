package process

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/VKGQQ/goproject_1/chat/common/message"
	"github.com/VKGQQ/goproject_1/chat/server/utils"
)

type UserProcess struct {
	Conn net.Conn
}

func (this *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal err:", err)
		return
	}
	var resMes message.Message
	resMes.Type = message.LoginResMesType
	var loginResMes message.LoginResMes
	if loginMes.UserId == 123 && loginMes.UserPwd == "123123" {
		loginResMes.Code = 200
	} else {
		loginResMes.Code = 500
		loginResMes.Error = "该用户不存在，请注册后再使用......"
	}
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal err:", err)
		return
	}
	resMes.Data = string(data)
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal err:", err)
		return
	}
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return
}
