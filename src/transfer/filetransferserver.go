package transfer

import (
	"fmt"
	"net"
)

// StartFileTransferServer 启动文件接受服务
func StartFileTransferServer(port int) error {
	listener, err := net.Listen("tcp", ":"+fmt.Sprintf("%d", port))
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println("新连接:", conn.RemoteAddr())
		go NewFileReceiver(conn).HandleRequest()
	}
}
