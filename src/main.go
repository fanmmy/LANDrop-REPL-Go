package main

import (
	"awesomeProject1/src/transfer"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	go func() {
		err := transfer.StartFileTransferServer(transfer.TransferPort)
		if err != nil {
			fmt.Println("启动文件接受服务失败", err)
			os.Exit(-1)
		}
	}()

	// 启动自动发现服务
	discovery := transfer.NewDiscoveryService(transfer.DiscoveryPort)
	// 创建一个通道用于同步服务启动状态
	serverReady := make(chan bool, 1)
	go discovery.StartDiscoveryServer(serverReady)
	<-serverReady

	// 在合适的时机发送广播报文，用于发现局域网内的服务
	go func() {
		for {
			//n, _ := net.ResolveUDPAddr("udp", "10.211.55.3:"+fmt.Sprintf("%d", transfer.DiscoveryPort))
			//discovery.RefreshSingle(n)
			discovery.RefreshWithBroadcast()
			time.Sleep(3 * time.Second)
		}
	}()

	// 接收退出信号 同时用于阻塞主线程
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	func() {
		<-quit
		// 收到退出信号后的清理操作
		// ...
		fmt.Println("收到退出信号")
		os.Exit(0)
	}()

}
