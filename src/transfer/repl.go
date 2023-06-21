package transfer

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
	//rootCmd.AddCommand(rzCmd)
}

func tip(str string) {
	l.Clean()
	fmt.Println(str)
	l.Refresh()
}

func tipWithPrompt(str, prompt string) {
	l.Clean()
	l.SetPrompt(fmt.Sprintf("\033[31m%s»\033[0m ", prompt))
	fmt.Println(str)
	l.Refresh()
}

func resetTipPrompt() {
	l.SetPrompt("\033[31m»\033[0m ")
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
				tipWithPrompt(notify.Msg, "ack")
				subPromptChan <- "yes"
				confirmation := <-subPromptChan

				if confirmation == "y" {
					notifyChan <- Notify{NotifyType: RespAcceptFile, Msg: "y"}
					break
				} else if confirmation == "n" {
					notifyChan <- Notify{NotifyType: RespAcceptFile, Msg: "n"}
					l.Refresh()
					break
				} else {
					tip("非法输入. 请再试一次.")
				}
			}
		case Info: // 提示信息
			tip(notify.Msg)
		case Prompt: // 交互式主菜单
			line := notify.Msg
			mainMenu(line)

		}
		resetTipPrompt()
	}
}

func tipHelp() {
	tip(

		"sf:      DONE 上传文件\n" +
			"sfd:     DONE 上传文件，GUI提示\n" +
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
