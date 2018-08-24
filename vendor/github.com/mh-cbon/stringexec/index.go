package stringexec

import (
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
)

// Return a new exec.Cmd object for the given command string
func Command(cmd string) (*exec.Cmd, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	if runtime.GOOS == "windows" {
		return WindowsCommand(cwd, cmd)
	}
	return FriendlyUnixCommand(cwd, cmd)
}

func WindowsCommand(cwd string, cmd string) (*exec.Cmd, error) {
	dir, err := ioutil.TempDir("", "stringexec")
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(dir+"/some.bat", []byte(cmd), 0766)
	if err != nil {
		return nil, err
	}

	oCmd := exec.Command("cmd", []string{"/C", dir + "/some.bat"}...)
	oCmd.Dir = cwd
	// defer os.Remove(tmpfile.Name()) // clean up // not sure how to clean it :x
	return oCmd, nil
}

func FriendlyUnixCommand(cwd string, cmd string) (*exec.Cmd, error) {
	oCmd := exec.Command("sh", []string{"-c", cmd}...)
	oCmd.Dir = cwd
	return oCmd, nil
}
