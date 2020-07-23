package fsync

type Action struct {
}

type Conf struct {
	LocalHome  string
	RemoteHome string
	Host       string
	Port       string
}

type ProcConf struct {
	IsServer bool
}

var pcf ProcConf
var gcf Conf
var act Action

func init() {
	pcf = ProcConf{}
	act = Action{}
}

func init() {

}
