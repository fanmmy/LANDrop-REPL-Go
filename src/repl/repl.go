package repl

import (
	"fmt"
	"github.com/chzyer/readline"
	"strings"
)

type Notify struct {
	NotifyType NotifyType
	Msg        string
}

type NotifyType int

const (
	ReqAcceptFile  NotifyType = iota //请求是否接收文件
	RespAcceptFile                   // 响应是否接收文件
	Info                             // 信息通知
	Prompt                           // 输入提示
	ProcessBarDone                   // 进度条完成
)

var l *readline.Instance

var mainLineChan = make(chan Notify)  // 一级输入管道
var subPromptChan = make(chan string) // 二级输入管道

func init() {
	l, _ = readline.NewEx(&readline.Config{
		Prompt:          "\033[31m»\033[0m ",
		HistoryFile:     "/tmp/readline.tmp",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
}

func tip(str string) {
	l.Clean()
	fmt.Println(str)
	l.Refresh()
}

func LoopRepl(notifyChan chan Notify) {

	// 处理控制台输入
	go handleConsoleInput()
	// 处理及时通知消息
	go handleNotifyChan(notifyChan)

	for {
		notify := <-mainLineChan
		switch notify.NotifyType {
		case ReqAcceptFile: // 接收文件请求
			for {
				l.SetPrompt("\033[31m>>>\033[0m ")
				l.Refresh()
				tip(notify.Msg)

				subPromptChan <- "yes"
				confirmation := <-subPromptChan

				if confirmation == "y" {
					notifyChan <- Notify{NotifyType: RespAcceptFile, Msg: "y"}
					break
				} else if confirmation == "n" {
					notifyChan <- Notify{NotifyType: RespAcceptFile, Msg: "n"}
					break
				} else {
					tip("非法输入. 请再试一次.")
				}
			}
		case Info: // 提示信息
			tip(notify.Msg)
		case Prompt: // 交互式主菜单
			line := notify.Msg
			switch {
			case strings.HasPrefix(line, "rz"):
				switch strings.TrimSpace(line[2:]) {
				case "vi":
					l.SetVimMode(true)
				case "emacs":
					l.SetVimMode(false)
				default:
					println("invalid mode:", line[2:])
				}
			case line == "login":
				for {
					l.SetPrompt("\033[31m>>>\033[0m ")
					l.Refresh()
					tip("登录中，请稍后... (y/n)")
					subPromptChan <- "yes"
					confirmation := <-subPromptChan

					if confirmation == "y" {
						tip("login yes")
						break
					} else if confirmation == "n" {
						tip("login no")
						break
					} else {
						tip("非法输入. 请再试一次.")
					}
				}
			case line == "help":
				tipHelp()
			case line == "?":
				tipHelp()
			case line == "exit":
				tip("退出程序...")
				return

			default:
				tip("非法输入. 请再试一次.(help/?)查看帮助")
			}

		}
		l.SetPrompt("\033[31m»\033[0m ")
		l.Refresh()

	}
}
func tipHelp() {
	tip(

		"rz:      TODO 上传\n" +
			"trz:     TODO 上传，通过GUI文件选择框\n" +
			"cfg:     TODO 配置\n" +
			"info:    TODO 查看本机信息")
}
func handleNotifyChan(notifyChan chan Notify) {
	for {
		notify := <-notifyChan
		mainLineChan <- notify
	}
}

func handleConsoleInput() {
	for {
		l.Refresh()
		line, _ := l.Readline()
		line = strings.TrimSpace(line)
		select {
		case <-subPromptChan:
			subPromptChan <- line
		default:
			mainLineChan <- Notify{NotifyType: Prompt, Msg: line}
		}
	}
}
