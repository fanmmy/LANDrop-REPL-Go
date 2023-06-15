package transfer

import (
	"awesomeProject1/src/utils"
	"bufio"
	"container/list"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type FileReceiver struct {
	*FileTransferSession
	transferQ   *list.List //接受队列 是*FileMetadata类型
	writingFile *os.File   //正在写的文件
}

func NewFileReceiver(conn net.Conn) *FileReceiver {
	crypto := &FileReceiver{
		FileTransferSession: &FileTransferSession{
			State:  HANDSHAKE1,
			Crypto: NewCrypto(),
			Conn:   conn,
		},
	}
	return crypto
}

func (r *FileReceiver) HandleRequest() {
	conn := r.Conn
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		// 第一次握手，响应本地公钥并将对方公钥设置到对象中
		if r.State == HANDSHAKE1 {
			err := r.Handshake1Process(reader)
			if err != nil {
				return
			}

		} else {
			decrypt, err := r.ReadAndDecrypt(reader)
			if err != nil {
				fmt.Println("数据解密失败", err)
				return
			}
			r.processData(decrypt)
		}

	}

}

func (r *FileReceiver) processData(data []byte) {
	if r.State == HANDSHAKE2 {
		fmt.Printf("第二次握手解密成功 %s \n", string(data))
		shake2Packs := HandShake2Pack{}
		err := json.Unmarshal(data, &shake2Packs)
		if err != nil {
			fmt.Println("文件元信息JSON解析失败")
			return
		}
		digest := r.Crypto.SessionKeyDigest()
		fmt.Println("验证吗是", digest)
		// 将这些文件添加到transferQ俩面
		r.transferQ = shake2Packs.FileList()

		// 响应确认消息
		_, err = r.EncryptAndSend([]byte("{\"response\":1}"))
		if err != nil {
			fmt.Println("文件元信息JSON解析失败")
			return
		}
		r.State = TRANSFERRING
		// 创建第一个文件
		r.createNextFile()

	} else if r.State == TRANSFERRING {
		// 开始传输了，就处理文件
		queue := r.transferQ
		if queue.Len() > 0 {
			// 出队一个元素
			currFileMeta := queue.Front().Value.(*FileMetadata)
			writeSize := utils.Min(currFileMeta.Size, int64(len(data)))
			// 已经读了多少字节了
			currFileMeta.Size -= writeSize
			//写文件
			_, err := r.writingFile.Write(data[:writeSize])
			if err != nil {
				return
			}
			data = data[writeSize:]
			if currFileMeta.Size == 0 { //当前文件传输完了
				fmt.Println("传输完成-", currFileMeta)
				// 出队一个元素
				queue.Remove(queue.Front())
				r.createNextFile()
			}
		}
	}
}

func (r *FileReceiver) createNextFile() {
	conn := r.Conn
	queue := r.transferQ
	if queue.Len() > 0 {
		currMeta := queue.Front().Value.(*FileMetadata)
		file, err := os.OpenFile(utils.UniqueFullPath(downloadDir, currMeta.FileName), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			return
		}
		if r.writingFile != nil { //关闭文件句柄
			err = r.writingFile.Close()
			if err != nil {
				return
			}
		}
		r.writingFile = file

	}
	if queue.Len() == 0 {
		if r.writingFile != nil { //关闭文件句柄
			err := r.writingFile.Close()
			if err != nil {
				return
			}
		}
		fmt.Println("所有文件都传输完了")
		r.State = FINISHED
		err := conn.Close()
		if err != nil {
			return
		}
	}
}