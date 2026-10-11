package process

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/VKGQQ/goproject_1/chat/common/message"
	"github.com/VKGQQ/goproject_1/chat/server/model"
)

// connWriteTimeout 单次发包（含帧头和帧体）的写超时时间。
// 必须设置：否则对端网络卡住时 Conn.Write 会无限期阻塞，
// 持有的 writeMu 也会一直不释放，导致其他 goroutine 无法给该用户发包。
const connWriteTimeout = 5 * time.Second

type UserProcess struct {
	Conn   net.Conn
	UserId int
	// writeMu 是“每连接一把”的写锁：同一连接的帧头+帧体必须在
	// 持锁期间连续写出，保证多个 goroutine（本人响应、群发转发、上线通知）
	// 并发发包时不会出现字节交错导致的帧失步。
	// 不同连接各有一把锁，互不影响、可完全并行。
	writeMu sync.Mutex
}

// Send 是向该用户连接写入一个完整协议帧的唯一入口。
// 帧格式：4 字节大端长度 + JSON 包体。
func (this *UserProcess) Send(data []byte) error {
	this.writeMu.Lock()
	defer this.writeMu.Unlock()
	// 每次发包前设置绝对截止时间，避免慢客户端把锁长期占住
	if err := this.Conn.SetWriteDeadline(time.Now().Add(connWriteTimeout)); err != nil {
		return err
	}
	var header [4]byte
	binary.BigEndian.PutUint32(header[:], uint32(len(data)))
	if _, err := this.Conn.Write(header[:]); err != nil {
		return err
	}
	if _, err := this.Conn.Write(data); err != nil {
		return err
	}
	return nil
}

func (this *UserProcess) NotifyOtherOnlineUser(userId int) {
	// 基于快照在锁外遍历并发送通知，避免持读锁做网络 IO
	for id, up := range userMgr.GetAllOnlineUsers() {
		if id == userId {
			continue
		}
		up.NotifyMeOnline(userId)
	}
}

func (this *UserProcess) NotifyMeOnline(userId int) {
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType
	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = message.UserOnline
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Marshal err:", err)
		return
	}
	mes.Data = string(data)
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err:", err)
		return
	}
	if err = this.Send(data); err != nil {
		fmt.Println("NotifyMeOnline err:", err)
	}
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
	user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)
	if err != nil {
		if errors.Is(err, model.ErrorUserNotExists) {
			loginResMes.Code = 500
			loginResMes.Error = err.Error()
		} else if errors.Is(err, model.ErrorUserPwd) {
			loginResMes.Code = 403
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "服务器内部错误"
		}
	} else {
		loginResMes.Code = 200
		this.UserId = user.UserId
		userMgr.AddOnlineUser(this)
		this.NotifyOtherOnlineUser(this.UserId)
		loginResMes.UsersId = userMgr.OnlineUserIds()
		fmt.Println(user, "登录成功")
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
	err = this.Send(data)
	return
}

func (this *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.Unmarshal err:", err)
		return
	}
	var resMes message.Message
	resMes.Type = message.RegisterResMesType
	var registerResMes message.RegisterResMes
	err = model.MyUserDao.Register(&registerMes.User)
	if err != nil {
		if errors.Is(err, model.ErrorUserExists) {
			registerResMes.Code = 505
			registerResMes.Error = err.Error()
		} else {
			registerResMes.Code = 506
			registerResMes.Error = "注册时发生未知错误"
		}
	} else {
		registerResMes.Code = 200
	}
	data, err := json.Marshal(registerResMes)
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
	err = this.Send(data)
	return
}
