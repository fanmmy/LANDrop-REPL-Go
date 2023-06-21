# LANDrop-Core-Go

基于Golang实现的LANDrop核心，参考 https://github.com/LANDrop/LANDrop 实现。

## Features

- [x] 完全兼容LANDrop协议，可与原生LANDrop混用通信
- [x] 通过UDP广播发现局域网内的其他设备
- [x] 通过TCP连接实现文件传输，借助`golang.org/x/crypto`实现传输消息加密
- [x] 重名文件接收时自动重命名
- [x] 实现文件传输进度显示
- [x] 提供REPL环境，实现接收文件确认
- [ ] 支持配置文件和main命令传参
- [x] REPL交互式环境
  - [x] sf - 发送文件，根据提示探出文件选择框并输入IP地址和端口号
  - [x] sfd - 直接发送文件
  - [x] help - 帮助信息
  - [ ] 发现局域网内其他Landrop服务并使用
  - [ ] info - 查看基础信息

目前实现的功能还比较糙，一些端口和配置都是硬编码的，后续将支持命令行传参并结合配置文件进行配置。

## Usage
### 等待接收文件

接收目录目前硬编码在`src/transfer/common_transer.go`中的`downloadDir`全局变量，后续将支持命令行传参和配置文件配置。

![等待接收文件](/Users/fmy/coding/Golang/LANDrop-Core-Go/doc/imgs/receive.gif)
### 主动发送文件

#### sfd命令（Send files with directory）

![等待接收文件](/Users/fmy/coding/Golang/LANDrop-Core-Go/doc/imgs/sfd.png)

#### sf命令（GUI文件选择框）

<video id="video" controls="" preload="none" poster="作者(图片地址)">
<source id="mp4" src="视频地址" type="video/mp4">
</video>


## Building

1. Clone this repository

```bash
git clone https://github.com/fanmmy/LANDrop-Core-Go.git
```

2. Generate vendor

```bash
go mod vendor
```

3. build main.go

```bash
go build -o landrop main.go 
```


