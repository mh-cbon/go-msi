package tpls

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/mattn/go-zglob"
	"github.com/mh-cbon/go-msi/manifest"
)

var funcMap = template.FuncMap{
	"dec": func(i int) int {
		return i - 1
	},
}

// Find all wxs fies in given directory
func Find(srcDir string, pattern string) ([]string, error) {
	glob := filepath.Join(srcDir, pattern)
	return zglob.Glob(glob)
}

// GenerateTemplate generates given src template to out file using given manifest
func GenerateTemplate(wixFile *manifest.WixManifest, src string, out string) error {
	tpl, err := template.New("").Funcs(funcMap).ParseFiles(src)
	if err != nil {
		return err
	}

	fileWriter, err := os.Create(out)
	if err != nil {
		return err
	}
	defer fileWriter.Close()
	err = tpl.ExecuteTemplate(fileWriter, filepath.Base(src), wixFile)
	if err != nil {
		return err
	}
	return nil
}
