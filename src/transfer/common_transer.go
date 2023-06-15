package transfer

import (
	"container/list"
	"fmt"
	"os"
)

// 下载路径
const downloadDir = "/Users/fmy/Downloads/LANDrop/"

// TransferQuanta 一次传输的数据量,单位字节数
const TransferQuanta = 64000

// DiscoveryPort 发现用的端口 TODO 启动时候发现端口占用应该提示出来
const DiscoveryPort = 52637

// TransferPort 发送文件用的端口，可变
const TransferPort = 7787

// State 定义自定义类型作为枚举类型
type State int

// 定义枚举值常量
const (
	HANDSHAKE1   State = iota // 0
	HANDSHAKE2                // 1
	TRANSFERRING              // 2
	FINISHED                  // 3
)

type HandShake2Resp struct {
	Response int `json:"response"`
}

func (r *HandShake2Resp) IsAccept() bool {
	return r.Response != 0
}

type HandShake2Pack struct {
	Files      []FileMetadata `json:"files"`
	DeviceName string         `json:"device_name"`
}

type FileMetadata struct {
	FileName string `json:"filename"`
	Size     int64  `json:"size"`
}

func (p *HandShake2Pack) FileList() *list.List {
	l := list.New()
	for i, _ := range p.Files {
		l.PushBack(&p.Files[i])
	}
	return l
}

// 返回值 metaQ ，fileQ
func getFiles(fileList ...string) (*list.List, *list.List) {

	metaQ := list.New()
	fileQ := list.New()
	// 遍历文件列表
	for _, filename := range fileList {
		// 打开文件
		file, err := os.Open(filename)
		if err != nil {
			fmt.Printf("Failed to open file %s: %s\n", filename, err)
			continue
		}
		// 获取文件信息
		fileInfo, err := file.Stat()
		if err != nil {
			fmt.Printf("Failed to get file info for %s: %s\n", filename, err)
			continue
		}

		// 获取文件名称和大小
		fileName := fileInfo.Name()
		fileSize := fileInfo.Size()

		metaQ.PushBack(&FileMetadata{fileName, fileSize})
		fileQ.PushBack(file)
		// 输出文件信息
		fmt.Printf("File: %s, Size: %d bytes\n", fileName, fileSize)
	}
	return metaQ, fileQ
}

// 链表转数组
func listToArray(l *list.List) []FileMetadata {
	var arr []FileMetadata
	// 遍历链表
	for e := l.Front(); e != nil; e = e.Next() {
		metadata := e.Value.(*FileMetadata)
		arr = append(arr, *metadata)
	}

	return arr
}






