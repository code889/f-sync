package fsync

import (
	"net"
	"os"
	"path/filepath"
)

var MyClose = func(c net.Conn) {
	_ = c.Close()
}

var MyGetPwd = func() string {
	t, _ := os.Getwd()
	return t
}

var MyCreateDir = func(path string) {

	path = filepath.Clean(path)
	_, err := os.Stat(filepath.Dir(path))
	if os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			panic(err)
		}
	}

}
