package transfer

import (
	"awesomeProject1/src/utils"
	"encoding/json"
	"net"
	"os"
	"runtime"
	"strconv"
)

// DiscoveryPack 发现设备的报文
type DiscoveryPack struct {
	Request    bool   `json:"request"`
	DeviceType string `json:"device_type"`
	DeviceName string `json:"device_name"`
	Port       uint16 `json:"port"`
}

// DiscoveryHost 发现的主机
type DiscoveryHost struct {
	Addr       string
	DeviceName string
	Port       uint16
}

type DiscoveryService struct {
	conn *net.UDPConn
	port int
}

func NewDiscoveryService(port int) *DiscoveryService {
	return &DiscoveryService{
		port: port,
	}
}

// RefreshWithBroadcast 刷新就是发送广播报文
func (d *DiscoveryService) RefreshWithBroadcast() {
	ips, _ := utils.GetBroadcastAddress()
	for _, ip := range ips {
		remoteAddr, _ := net.ResolveUDPAddr("udp", ip+":"+strconv.Itoa(DiscoveryPort))
		_, err := d.conn.WriteToUDP([]byte("{\"request\":true}"), remoteAddr)
		if err != nil {
			log.Error(err)
		}
	}
}

func (d *DiscoveryService) RefreshSingle(remoteAddr *net.UDPAddr) {
	_, err := d.conn.WriteToUDP([]byte("{\"request\":true}"), remoteAddr)
	if err != nil {
		log.Error(err)
	}
}

func (d *DiscoveryService) StartDiscoveryServer() {
	// 解析得到UDP地址
	addr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(d.port))
	// 在UDP地址上建立UDP监听,得到连接
	d.conn, err = net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
		return
	}
	// 发送一次广播报文
	go d.RefreshWithBroadcast()

	d.processConnection()
}

func (d *DiscoveryService) processConnection() {

	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(d.conn)

	// 建立缓冲区
	buffer := make([]byte, 1024)

	for {
		//从连接中读取内容,丢入缓冲区
		i, udpAddr, e := d.conn.ReadFromUDP(buffer)
		// 第一个是字节长度,第二个是udp的地址
		if e != nil {
			continue
		}

		if utils.IsLocalAddr(udpAddr) {
			continue
		}
		req := DiscoveryPack{}
		err := json.Unmarshal(buffer[:i], &req)
		if err != nil {
			continue
		}

		if !req.Request { //不是请求
			// 将发现结果抛出用于业务处理
			handleDiscoveryResult(&DiscoveryHost{
				DeviceName: req.DeviceName,
				Addr:       udpAddr.IP.String(),
				Port:       req.Port,
			})
			continue
		}

		// 准备响应发现请求
		jsonByte, _ := json.Marshal(localInfoPack())
		_, err = d.conn.WriteToUDP(jsonByte, udpAddr)

	}
}

func handleDiscoveryResult(newHost *DiscoveryHost) {
	//log.Info("--发现主机，待处理展示--")
	//log.Info(newHost.DeviceName)
}

// 生成本地数据包信息，用以响应UDP广播
func localInfoPack() DiscoveryPack {
	var pack = DiscoveryPack{DeviceName: runtime.GOOS, Port: TransferPort, Request: false}
	var err error
	pack.DeviceName, err = os.Hostname()
	if err != nil {
		pack.DeviceName = "未知"
	}
	pack.DeviceType = runtime.GOOS

	return pack
}
