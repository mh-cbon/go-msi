package guid

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mh-cbon/go-msi/util"
)

// Make generates an uuid compatible with msi requirements.
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
			return "", fmt.Errorf("%s: %v: \n%s", cmd, err, out)
		}
		sout := string(out)
		sout = strings.TrimSpace(sout)
		sout = strings.ToUpper(sout)
		return sout, nil
	}

	cmd := "uuidgen"
	args := []string{}
	if runtime.GOOS == "linux" {
		args = append(args, "-r")
	}
	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s: %v: \n%s", cmd, err, out)
	}
	sout := string(out)
	sout = strings.TrimSpace(sout)
	sout = strings.ToUpper(sout)
	return sout, nil
}
