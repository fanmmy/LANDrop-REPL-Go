package transfer

import (
	"bufio"
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vbauerster/mpb/v8"
	"io"
	"net"
	"os"
	"runtime"
)

type FileSender struct {
	*FileTransferSession
	transferQ *list.List //传输文件的源信息队列
	filesQ    *list.List //传输文件的原始句柄队列
	Process
	bar *mpb.Bar
}

func NewFileSender(conn net.Conn) *FileSender {
	crypto := &FileSender{
		FileTransferSession: &FileTransferSession{
			State:  HANDSHAKE1,
			Crypto: NewCrypto(),
			Conn:   conn,
		},
	}
	return crypto
}

func (s *FileSender) StartSession() {
	s.Conn.Write(s.Crypto.LocalPublicKey())
}

// handshake1Finished 第一次握手结束就开始发送文件元信息,此为第二次握手
func (s *FileSender) handshake1Finished() {
	transferMetaDataArr := listToArray(s.transferQ)
	pack := HandShake2Pack{
		Files:      transferMetaDataArr,
		DeviceName: runtime.GOOS,
	}
	marshal, err := json.Marshal(pack)
	fmt.Println("第二次握手发送数据:", string(marshal))
	if err != nil {
		return
	}
	_, err = s.EncryptAndSend(marshal)
	if err != nil {
		log.Error(err)
		return
	}

}

func (s *FileSender) SendFiles(fileList ...string) {
	// 启动一个会话，即开始第一次握手
	var totalSize int64
	s.transferQ, s.filesQ, totalSize = getFiles(fileList...)
	conn := s.Conn
	s.StartSession()
	reader := bufio.NewReader(conn)
	for {
		if s.State == HANDSHAKE1 {
			err := s.Handshake1Process(reader)
			if err != nil {
				return
			}
			// 第一次握手成功 发送JSON数据，开始第二次握手
			s.handshake1Finished()
			fmt.Println("第一次握手成功")

		} else if s.State == HANDSHAKE2 {
			decrypt, err := s.ReadAndDecrypt(reader)
			if err != nil {
				//fmt.Println("第二次握手接收数据失败:", err)
				return
			}
			//开始传输文件，初始化进度条
			s.Process = Process{
				Receive: false,
				AllNum:  s.transferQ.Len(),
			}
			_, s.bar = NewProcessBar(totalSize, &s.Process)
			err = s.processData(decrypt)
			if err != nil {
				log.Error(err)
				return
			}

		}
	}
}

func (s *FileSender) processData(data []byte) error {
	resp := HandShake2Resp{}
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return errors.New("第二次握手接收端报文反序列化失败")
	}
	log.Info("第二次握手服务器响应:--", resp)
	if !resp.IsAccept() {
		s.State = FINISHED
		err := s.Conn.Close()
		if err != nil {
			return err
		}
		return errors.New("文件接收端不接受文件列表")
	}
	//  开始传输文件
	fileQuantaBuffer := make([]byte, TransferQuanta)
	for s.filesQ.Len() > 0 {
		currFile := s.filesQ.Front().Value.(*os.File)
		n, err := currFile.Read(fileQuantaBuffer)
		//说明一下当前正在传输谁
		s.Process.DoingFile = s.transferQ.Front().Value.(*FileMetadata).FileName
		// 更新进度条
		s.bar.IncrInt64(int64(n))
		// 到达文件结尾，发送完毕
		if err == io.EOF {
			log.Info("发送完毕", s.transferQ.Front().Value.(*FileMetadata).FileName)
			//发送完毕需要出队列并且更新进度信息
			s.Process.DoneNum++
			s.filesQ.Remove(s.filesQ.Front())
			s.transferQ.Remove(s.transferQ.Front())
			err := currFile.Close()
			if err != nil {
				return err
			}
			continue
		}
		n, err = s.EncryptAndSend(fileQuantaBuffer[:n])

		if err != nil {
			return err
		}
	}
	if s.filesQ.Len() == 0 {
		log.Info("所有文件都发送完毕")
		fmt.Println("完成")
	}
	return nil
}
