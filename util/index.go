package util

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
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
