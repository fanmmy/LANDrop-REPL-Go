package transfer

import (
	"awesomeProject1/src/utils"
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"net"
)

type FileTransferSession struct {
	State      State       //传输状态
	Crypto     *Crypto     //加密对象
	Conn       net.Conn    //网络连接
	NotifyChan chan Notify //通信通道

}

// EncryptAndSend 加密并传输
func (b *FileTransferSession) EncryptAndSend(data []byte) (int, error) {
	encryptedData, err := b.Crypto.Encrypt(data)
	if err != nil {
		return 0, err

	}
	size := len(encryptedData)
	sendData := make([]byte, size+2)
	sendData[0] = byte(size >> 8 & 0xFF)
	sendData[1] = byte(size & 0xFF)
	copy(sendData[2:], encryptedData)
	return b.Conn.Write(sendData)
}

func (b *FileTransferSession) Handshake1Process(reader *bufio.Reader) error {
	keySize := PublicKeySize

	buffer := make([]byte, keySize)
	n, _ := reader.Read(buffer)
	if n < keySize {
		return errors.New("HANDSHAKE1 可读 size < 32 ，第一次握手失败")
	}
	//设置远程公钥
	err := b.Crypto.SetRemotePublicKey(buffer[:keySize])
	if err != nil {
		return err
	}

	// 设置当前状态为第二次握手
	b.State = HANDSHAKE2
	return nil
}

func (b *FileTransferSession) ReadAndDecrypt(reader *bufio.Reader) ([]byte, error) {
	sizeBuffer := make([]byte, 2)
	crypto := b.Crypto
	n, _ := reader.Read(sizeBuffer)
	if n < 2 {
		log.Warn("数据长度小于2，跳过")
		return nil, io.EOF
	}
	// 解析消息体长度
	size := binary.BigEndian.Uint16(sizeBuffer)

	fileEntryData := make([]byte, size)

	//  以下处理拆包
	readLen := 0
	for int(size) > readLen {
		canReadLen := utils.Min(reader.Size(), int(size)-readLen)
		n, _ := reader.Read(fileEntryData[readLen : canReadLen+readLen])
		readLen += n
	}
	// 解密
	decrypt, err := crypto.Decrypt(fileEntryData[:])
	return decrypt, err
}
