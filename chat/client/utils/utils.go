package utils

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"

	"github.com/VKGQQ/goproject_1/chat/common/message"
)

type Transfer struct {
	Conn net.Conn
	Buf  [8096]byte
}

func (this *Transfer) readPkg() (mes message.Message, err error) {
	_, err = this.Conn.Read(this.Buf[:4])
	if err != nil {
		return
	}
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(this.Buf[:4])
	n, err := this.Conn.Read(this.Buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		return
	}
	err = json.Unmarshal(this.Buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarshal err:", err)
		return
	}
	return
}

func (this *Transfer) WritePkg(data []byte) (err error) {
	var pkgLen uint32
	pkgLen = uint32(len(data))
	binary.BigEndian.PutUint32(this.Buf[0:4], pkgLen)
	n, err := this.Conn.Write(this.Buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write err:", err)
		return
	}
	n, err = this.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write err:", err)
		return
	}
	return
}
