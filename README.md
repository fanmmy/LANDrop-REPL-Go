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
- [ ] REPL环境完善
  - [ ] info - 查看基础信息
  - [ ] rz - 传输文件
  - [ ] help/? - 帮助信息

目前实现的功能还比较糙，一些端口和配置都是硬编码的，后续将支持命令行传参并结合配置文件进行配置。

## Building

1. Clone this repository

```bash
git clone https://github.com/fanmmy/LANDrop-Core-Go.git
```

2. Generate vendor

```bash
go mod vendor
```

3. Run main.go

```bash
go run src/main.go
```

