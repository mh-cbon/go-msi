package util

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// GetBinPath Find path of the current binary file on the file system
func GetBinPath() (string, error) {
	var err error
	wd := ""
	if filepath.Base(os.Args[0]) == "main" { // go run ...
		wd, err = os.Getwd()
	} else {
		bin, err2 := exec.LookPath(os.Args[0])
		if err2 == nil {
			wd = filepath.Dir(bin)
		}
	}
	return wd, err
}

//CopyFile copy file src to dst.
func CopyFile(dst, src string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}

//ComputeSha256 computes the sha256 value of a file content.
func ComputeSha256(filepath string) (string, error) {
	hasher := sha256.New()
	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

//Exec computes the sha256 value of a file content.
func Exec(w string, args ...string) (string, error) {
	cmd := exec.Command(w, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}
