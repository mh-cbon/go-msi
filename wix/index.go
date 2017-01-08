package wix

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mh-cbon/go-msi/manifest"
)

var eol = "\r\n"

// GenerateCmd generates required command lines to produce an msi package,
func GenerateCmd(wixFile *manifest.WixManifest, templates []string, msiOutFile string, arch string) string {

	cmd := ""

	for i, dir := range wixFile.RelDirs {
		sI := strconv.Itoa(i)
		cmd += "heat dir " + dir + " -nologo -cg AppFiles" + sI
		cmd += " -gg -g1 -srd -sfrag -template fragment -dr APPDIR" + sI
		cmd += " -var var.SourceDir" + sI
		cmd += " -out AppFiles" + sI + ".wxs"
		cmd += eol
	}
	cmd += "candle"
	if arch != "" {
		if arch == "386" {
			arch = "x86"
		} else if arch == "amd64" {
			arch = "x64"
		}
		cmd += " -arch " + arch
	}
	for i, dir := range wixFile.RelDirs {
		sI := strconv.Itoa(i)
		cmd += " -dSourceDir" + sI + "=" + dir
	}
	for i := range wixFile.Directories {
		sI := strconv.Itoa(i)
		cmd += " AppFiles" + sI + ".wxs"
	}
	for _, tpl := range templates {
		cmd += " " + filepath.Base(tpl)
	}
	cmd += eol
	cmd += "light -ext WixUIExtension -ext WixUtilExtension -sacl -spdb "
	cmd += " -out " + msiOutFile
	for i := range wixFile.Directories {
		sI := strconv.Itoa(i)
		cmd += " AppFiles" + sI + ".wixobj"
	}
	for _, tpl := range templates {
		cmd += " " + strings.Replace(filepath.Base(tpl), ".wxs", ".wixobj", -1)
	}
	cmd += eol

	return cmd
}
