package process

import (
	"fmt"

	"encoding/json"

	"github.com/VKGQQ/goproject_1/chat/common/message"
)

func outputGroupMes(mes *message.Message) {
	//反序列化mes.Data
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err.Error())
		return
	}

	//显示信息
	info := fmt.Sprintf("接收到用户id:\t%d 发来的信息:\t%s",
		smsMes.UserId, smsMes.Content)
	fmt.Println(info)
	fmt.Println()
}
