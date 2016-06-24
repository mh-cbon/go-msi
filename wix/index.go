package wix

func GenerateCmd (wixFile *WixManifest, templates []string, msiOutFile string, arch string) string {

  cmd := ""

  for i, dir := range wixFile.RelDirs {
    sI := strconv.Itoa(i)
    cmd += "heat dir "+dir+" -nologo -cg AppFiles"+sI
    cmd += " -gg -g1 -srd -sfrag -template fragment -dr APPDIR"+sI
    cmd += " -var var.SourceDir"+sI
    cmd += " -out AppFiles"+sI+".wxs"
    cmd += "\r\n"
  }
  cmd += "candle -arch "+arch
  for i, dir := range wixFile.RelDirs {
    sI := strconv.Itoa(i)
    cmd += " -dSourceDir"+sI+"="+dir
  }
  for i, _ := range wixFile.Directories {
    sI := strconv.Itoa(i)
    cmd += " AppFiles"+sI+".wxs"
  }
  for _, tpl := range templates {
    cmd += " "+tpl
  }
  cmd += "\r\n"
  cmd += "light -ext WixUIExtension -ext WixUtilExtension -sacl -spdb "
  cmd += " -out "+msiOutFile
  for i, _ := range wixFile.Directories {
    sI := strconv.Itoa(i)
    cmd += " AppFiles"+sI+".wixobj"
  }
  for _, tpl := range templates {
    cmd += " "+strings.Replace(tpl, ".wxs", ".wixobj")
  }
  cmd += "\r\n"

  return cmd
}
