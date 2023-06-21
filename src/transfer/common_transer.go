package transfer

import (
	"container/list"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
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

var log = logrus.New()

func init() {
	file, err := os.OpenFile("transfer.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
		log.SetReportCaller(true)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
}

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

func (p *HandShake2Pack) StringPrintFile() string {
	// 找到最长的 FileName 长度
	maxLen := 0
	for _, data := range p.Files {
		if len(data.FileName) > maxLen {
			maxLen = len(data.FileName) + 6
		}
	}
	// 构建字符串缓冲区
	var builder strings.Builder

	// 输出表头
	builder.WriteString("Filename")
	builder.WriteString(strings.Repeat(" ", maxLen-8))
	builder.WriteString("Size\n")

	// 输出每个 FileMetadata
	for _, data := range p.Files {
		builder.WriteString(data.FileName)
		builder.WriteString(strings.Repeat(" ", maxLen-len(data.FileName)))
		builder.WriteString(fmt.Sprintf("%s\n", humanize.Bytes(uint64(data.Size))))
	}

	// 获取最终拼接好的字符串
	return builder.String()

}

func (p *HandShake2Pack) FileList() *list.List {
	l := list.New()
	for i, _ := range p.Files {
		l.PushBack(&p.Files[i])
	}
	return l
}

func (p *HandShake2Pack) TotalSize() int64 {
	var totalSize int64 = 0
	for _, meta := range p.Files {
		totalSize += meta.Size
	}
	return totalSize
}

// 返回值 metaQ ，fileQ,totalSize
func getFiles(fileList []string) (*list.List, *list.List, int64, error) {

	metaQ := list.New()
	fileQ := list.New()
	var totalSize int64 = 0
	// 遍历文件列表
	for _, filename := range fileList {
		// 打开文件
		file, err := os.Open(filename)
		if err != nil {
			log.Errorf("Failed to open file %s: %s\n", filename, err)
			return nil, nil, 0, err
		}
		// 获取文件信息
		fileInfo, err := file.Stat()
		if err != nil {
			log.Errorf("Failed to get file info for %s: %s\n", filename, err)
			return nil, nil, 0, err
		}

		// 获取文件名称和大小
		fileName := fileInfo.Name()
		fileSize := fileInfo.Size()
		totalSize += fileSize
		metaQ.PushBack(&FileMetadata{fileName, fileSize})
		fileQ.PushBack(file)

	}
	return metaQ, fileQ, totalSize, nil
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
