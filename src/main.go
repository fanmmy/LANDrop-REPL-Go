package main

import (
	"awesomeProject1/src/repl"
	"awesomeProject1/src/transfer"
	"log"
)

func main() {

	notifyChan := make(chan repl.Notify) // 及时通知

	go func() {
		err := transfer.StartFileTransferServer(transfer.TransferPort, notifyChan)
		if err != nil {
			log.Fatal("启动文件接受服务失败", err)
		}
	}()
	// 启动自动发现服务
	go transfer.NewDiscoveryService(transfer.DiscoveryPort).StartDiscoveryServer()
	// 主线程循环交互式环境
	repl.LoopRepl(notifyChan)

}
