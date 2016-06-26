package guid

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mh-cbon/go-msi/util"
)

func Make() (string, error) {
	if runtime.GOOS == "windows" {

		b, err := util.GetBinPath()
		if err != nil {
			return "", err
		}

		cmd := "cscript.exe"
		args := []string{filepath.Join(filepath.Base(b), "utils", "myuuid.vbs")}
		out, err := exec.Command(cmd, args...).CombinedOutput()
		if err != nil {
			return "", err
		}
		sout := string(out)
		sout = strings.TrimSpace(sout)
		sout = strings.ToUpper(sout)
		return sout, nil
	} else {
		cmd := "uuidgen"
		args := []string{"-r"}
		out, err := exec.Command(cmd, args...).CombinedOutput()
		if err != nil {
			return "", err
		}
		sout := string(out)
		sout = strings.TrimSpace(sout)
		sout = strings.ToUpper(sout)
		return sout, nil
	}
}
