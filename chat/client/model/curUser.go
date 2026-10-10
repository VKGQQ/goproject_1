package model

import (
	"net"

	"github.com/VKGQQ/goproject_1/chat/common/message"
)

type CurUser struct {
	Conn net.Conn
	message.User
}
