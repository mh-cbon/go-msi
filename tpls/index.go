package tpls

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/mattn/go-zglob"
	"github.com/mh-cbon/go-msi/manifest"
)

// find all wxs fies in given directory
func Find(srcDir string, pattern string) ([]string, error) {
	glob := filepath.Join(srcDir, pattern)
	return zglob.Glob(glob)
}

// Generate given src file to out file using given manifest
func GenerateTemplate(wixFile *manifest.WixManifest, src string, out string) error {
	tpl, err := template.ParseFiles(src)
	if err != nil {
		return err
	}

	fileWriter, err := os.Create(out)
	if err != nil {
		return err
	}
	defer fileWriter.Close()
	err = tpl.Execute(fileWriter, wixFile)
	if err != nil {
		return err
	}
	return nil
}
