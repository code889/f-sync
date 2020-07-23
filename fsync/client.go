package fsync

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Client struct {
}

func runCmd(cmd string, args ...string) (err error) {

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("%v", e))
		}
	}()

	var exists bool
	var f func(args ...string)
	var cmds map[string]func(args ...string)
	cmds = map[string]func(args ...string){
		"set":  act.Set,
		"show": act.Show,
		"push": act.Push,
		"pull": act.Pull,
	}
	f, exists = cmds[cmd]
	if !exists {
		panic(fmt.Sprintf("command %v not exists ", cmd))
	}

	f(args...)
	return
}

func loopCmd() {
	var raw string
	var cmd string
	var err error
	var args []string
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("请输入命令:")
		raw, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println(fmt.Sprintf("%v", err))
			break
		}
		fields := strings.Fields(strings.TrimSpace(raw))
		if len(fields) == 0 {
			continue
		}
		fmt.Println(fields)
		cmd, args = fields[0], fields[1:]
		fmt.Println(cmd)
		if cmd == "exit" {
			break
		}

		err = runCmd(cmd, args...)
		if err != nil {
			fmt.Println(fmt.Sprintf("命令 %v 错误: %v", raw, err))
		}
	}
}

func (c *Client) Run() {
	loopCmd()
}

func (c *Client) Test() {
	gcf.Host = "127.0.0.1"
	gcf.Port = "8887"
	gcf.LocalHome = MyGetPwd()
	gcf.RemoteHome = "/xdfapp/tmp"
	a := Action{}
	a.Push("fsync/request.go")
}
