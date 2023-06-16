package test

import (
	"awesomeProject1/src/transfer"
	"fmt"
	"net"
	"testing"
)

func TestSendFiles(t *testing.T) {

	conn, err := net.Dial("tcp", "10.211.55.3:64636")
	if err != nil {
		fmt.Println("连接服务器失败:", err)
		return
	}
	defer conn.Close()
	s := transfer.NewFileSender(conn)

	s.SendFiles("/Users/fmy/Downloads/SunloginClient_12.5.2.46788.dmg")

}
