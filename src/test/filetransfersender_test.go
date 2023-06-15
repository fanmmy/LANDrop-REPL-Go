package test

import (
	"awesomeProject1/src/transfer"
	"fmt"
	"net"
	"testing"
)

func TestSendFiles(t *testing.T) {

	conn, err := net.Dial("tcp", "10.211.55.3:58313")
	if err != nil {
		fmt.Println("连接服务器失败:", err)
		return
	}
	defer conn.Close()
	s := transfer.NewFileSender(conn)

	s.SendFiles("/Users/fmy/Downloads/xx.xls", "/Users/fmy/Downloads/Excel公式.mp4")

}
