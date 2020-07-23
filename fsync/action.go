package fsync

import (
	"fmt"
	"os"
	"path/filepath"
)

func (a *Action) Set(args ...string) {
	var key, val string
	if len(args) != 2 {
		panic("set command num error")
	}
	key, val = args[0], args[1]
	switch key {
	case "host":
		gcf.Host = val
	case "port":
		gcf.Port = val
	case "remote":
		gcf.RemoteHome = val
	case "local":
		gcf.LocalHome = val

	default:
		panic("set command error")

	}
}

func (a *Action) Show(args ...string) {

	var info string
	info = `
当前目录: %v
目标目录: %v
目标主机: %v
目标端口: %v
`
	info = fmt.Sprintf(info, gcf.LocalHome, gcf.RemoteHome, gcf.Host, gcf.Port)
	fmt.Println(info)

}
func (a *Action) Push(args ...string) {

	if len(args) != 1 {
		panic("push param error")
	}
	path := args[0]
	t := Task{}
	local := filepath.Clean(gcf.LocalHome+"/"+path)
	f, err := os.Stat(local)
	if err != nil {
		panic(err)
	}
	remote := filepath.Clean(gcf.RemoteHome+"/"+path)
	if f.IsDir() {
		t.SendDir(gcf.Host, gcf.Port, local, remote)
	} else {
		t.SendFile(gcf.Host, gcf.Port, local, remote)
	}

}

func (a *Action) Pull(args ...string) {

	if len(args) != 1 {
		panic("push param error")
	}
	path := args[0]
	t := Task{}

	local := filepath.Clean(gcf.LocalHome+"/"+path)
	remote := filepath.Clean(gcf.RemoteHome+"/"+path)
	MyCreateDir(local)
	t.RecvFile(gcf.Host, gcf.Port, remote, local)


}