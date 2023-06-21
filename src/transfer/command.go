package transfer

import (
	"fmt"
	"github.com/ncruces/zenity"
	"github.com/spf13/cobra"
	"net"
	"os"
	"strings"
)

var cmdMap = initCmdMap()

// 返回值是否退出
func mainMenu(line string) {
	args := strings.Fields(line)
	if len(args) == 0 {
		return
	}
	command := args[0]
	subArgs := args[1:]

	switch command {
	case "sfd":
		sfdCmd := cmdMap["sfd"]
		sfdCmd.SetArgs(subArgs[:])
		if err := sfdCmd.Execute(); err != nil {
			tip(err.Error())
			return
		}
		return

	case "sf":
		sfCmd := cmdMap["sf"]

		tipWithPrompt("请选择要发送的文件: ", "sf")

		multiple, err := zenity.SelectFileMultiple()
		if err != nil {
			tip(err.Error())
			return
		}
		tip(strings.Join(multiple, "\n"))

		tip("请输入IP和端口示例(192.168.3.1:54321): ")

		for {
			subPromptChan <- "yes"
			ipPort := <-subPromptChan
			addr, err := net.ResolveTCPAddr("tcp", ipPort)
			if err != nil {
				tip(err.Error())
				continue
			}

			args := make([]string, 0)
			args = append(args, multiple...)
			args = append(args, "--host")
			args = append(args, addr.IP.String())
			args = append(args, "--port")
			args = append(args, fmt.Sprintf("%d", addr.Port))

			sfCmd.SetArgs(args[:])

			if err = sfCmd.Execute(); err != nil {
				tip(err.Error())
				return
			}
			break
		}

	case "help":
		tipHelp()
	case "?":
		tipHelp()
	case "exit":
		tip("退出程序...")
		os.Exit(0)

	default:
		tip("非法输入,Usage:\n")
		tipHelp()
	}
	return
}

func initCmdMap() map[string]*cobra.Command {
	// 创建 map
	commands := make(map[string]*cobra.Command)

	var sfCmd = &cobra.Command{
		Use:   "sf",
		Short: "Send files step by step",
		Args:  cobra.MinimumNArgs(1),
		RunE:  handleSf,
	}
	var sfdCmd = &cobra.Command{
		Use:   "sfd",
		Short: "Send files with directory",
		Args:  cobra.MinimumNArgs(1),
		RunE:  handleSf,
	}
	var rootCmd = &cobra.Command{
		Use:   "repl",
		Short: "Repl"}
	commands["sf"] = sfCmd
	commands["sfd"] = sfdCmd
	commands["root"] = rootCmd

	// 添加命令标志
	sfCmd.Flags().StringP("host", "H", "", "Host address")
	sfCmd.Flags().Uint16P("port", "P", 0, "Port")
	// 设置 host 参数为必需参数
	sfCmd.MarkFlagRequired("host")
	sfCmd.MarkFlagRequired("port")
	// 添加命令标志
	sfdCmd.Flags().StringP("host", "H", "", "Host address")
	sfdCmd.Flags().Uint16P("port", "P", 0, "Port")
	// 设置 host 参数为必需参数
	sfdCmd.MarkFlagRequired("host")
	sfdCmd.MarkFlagRequired("port")
	return commands
}

func handleSf(cmd *cobra.Command, args []string) error {

	ip, err := cmd.Flags().GetString("host")

	port, err := cmd.Flags().GetUint16("port")
	if err != nil {
		fmt.Println(err)
		cmd.Usage()
		return err
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		//fmt.Println("连接服务器失败:", err)
		return err
	}
	//defer conn.Close()
	go func() {
		err := NewFileSender(conn).SendFiles(args[:]...)
		if err != nil {
			tip(err.Error())
		}
	}()
	return nil
}
