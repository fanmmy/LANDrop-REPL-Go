package utils

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

// UniqueFullPath 返回文件全路径，若文件已存在就重命名位后缀(n)
func UniqueFullPath(dir string, fileName string) string {
	newFilePath := filepath.Join(dir, fileName)
	ext := filepath.Ext(newFilePath)
	baseName := strings.TrimSuffix(filepath.Base(newFilePath), ext)
	// 循环检查是否存在重名文件，如果存在，则生成新的文件名
	i := 1
	for {
		// 检查文件是否存在
		_, err := os.Stat(newFilePath)
		if os.IsNotExist(err) {
			// 文件不存在，可以使用新的文件名
			break
		}
		// 生成新的文件名
		newBaseName := fmt.Sprintf("%s(%d)", baseName, i)
		newFilePath = filepath.Join(dir, newBaseName+ext)
		// 文件存在，继续尝试下一个文件名
		i++
	}
	return newFilePath
}

func Min[T int | int64](x, y T) T {
	if x < y {
		return x
	}
	return y
}

// IsLocalAddr 判断IP地址是不是本地IP
func IsLocalAddr(udpAddr *net.UDPAddr) bool {
	localAddrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
		return false
	}

	// 扫描所有IPV4的地址
	for _, addr := range localAddrs {

		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {

			if udpAddr.IP.String() == ipNet.IP.String() {
				return true
			}
		}
	}
	return false
}

// GetBroadcastAddress 返回广播地址列表
func GetBroadcastAddress() ([]string, error) {
	broadcastAddress := []string{}

	interfaces, err := net.Interfaces() // 获取所有网络接口
	if err != nil {
		return broadcastAddress, err
	}

	for _, face := range interfaces {
		// 选择 已启用的、能广播的、非回环 的接口
		if (face.Flags & (net.FlagUp | net.FlagBroadcast | net.FlagLoopback)) == (net.FlagBroadcast | net.FlagUp) {
			addrs, err := face.Addrs() // 获取该接口下IP地址
			if err != nil {
				return broadcastAddress, err
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok { // 转换成 IPNet { IP Mask } 形式
					if ipnet.IP.To4() != nil { // 只取IPv4的
						var fields net.IP // 用于存放广播地址字段（共4个字段）
						for i := 0; i < 4; i++ {
							fields = append(fields, (ipnet.IP.To4())[i]|(^ipnet.Mask[i])) // 计算广播地址各个字段
						}
						broadcastAddress = append(broadcastAddress, fields.String()) // 转换为字符串形式
					}
				}
			}
		}
	}

	return broadcastAddress, nil
}
