package process

import (
	"encoding/json"
	"fmt"

	"github.com/VKGQQ/goproject_1/chat/common/message"
)

type SmsProcess struct {
}

func (this *SmsProcess) SendGroupMes(mes *message.Message) {

	//遍历服务器端的onlineUsers map[int]*UserProcess,
	//将消息转发取出.
	//取出mes的内容 SmsMes
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	mes.Type = message.SmsTransferMesType
	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	// 基于快照在锁外遍历转发，避免持读锁做网络 IO
	for id, up := range userMgr.GetAllOnlineUsers() {
		//这里，还需要过滤到自己,即不要再发给自己
		if id == smsMes.UserId {
			continue
		}
		// 走收信人自己的写锁，保证帧完整
		if err := up.Send(data); err != nil {
			fmt.Println("转发消息失败 err=", err)
		}
	}
}
