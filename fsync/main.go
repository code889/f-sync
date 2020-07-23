package fsync

import (
	"flag"
	"fmt"
	"os"

	"github.com/code889/f-conf/fconf"
)

func runClient() {
	c := &Client{}
	c.Run()
	// c.Test()
}
func runServer() {
	s := &Server{}
	s.run()
}
func initConf(cfp string) {
	if cfp != "" {
		_, err := os.Stat(cfp)
		if os.IsNotExist(err) {
			panic(fmt.Sprintf("config %v not exists", cfp))
		}
		i := fconf.Icf{}
		i.Load(cfp)
		gcf = Conf{
			LocalHome:  i.Opt("local", "home"),
			RemoteHome: i.Opt("remote", "home"),
			Host:       i.Opt("remote", "host"),
			Port:       i.Opt("remote", "port"),
		}
	} else {
		gcf = Conf{LocalHome: MyGetPwd(), Host: "127.0.0.1", Port: "8887"}
	}

}

func Setup() {

	var mod string
	var cfp string

	flag.StringVar(&mod, "m", "server", "模式，默认为服务器")
	flag.StringVar(&cfp, "c", "", "配置路径")

	flag.Parse()
	initConf(cfp)

	if mod == "server" {
		pcf.IsServer = true
		fmt.Println("running in server mode")
		runServer()
	} else {
		pcf.IsServer = false
		fmt.Println("running in client mode")
		runClient()
	}

}
func Main() {
	Setup()
}
