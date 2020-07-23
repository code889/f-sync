package fsync

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func walk(path string) []string {

	path = filepath.Clean(path)
	var names []string
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		if f.IsDir() {
			tnames := walk(path + "/" + f.Name())
			names = append(names, tnames...)
		} else {
			names = append(names, path+"/"+f.Name())
		}
	}
	return names
}

func printWalk() {
	names := walk(".")
	for _, name := range names {
		fmt.Println(name)
	}

}
