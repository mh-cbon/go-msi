package util

import (
	"os"
	"os/exec"
	"path/filepath"
)

// Find path of the current binary file on the file system
func GetBinPath() (string, error) {
	var err error
	wd := ""
	if filepath.Base(os.Args[0]) == "main" { // go run ...
		wd, err = os.Getwd()
	} else {
		bin := ""
		bin, err = exec.LookPath(os.Args[0])
		if err == nil {
			wd = filepath.Dir(bin)
		}
	}
	return wd, err
}
