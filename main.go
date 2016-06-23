package main

import (
  "encoding/json"
  "io/ioutil"
  "path/filepath"
  "text/template"
  "os"

  "golang.org/x/text/encoding/charmap"
  "golang.org/x/text/transform"
)


type WixManifest struct {
    Product       string          `json:"product"`
    Company       string          `json:"company"`
    Version       string          `json:"version"`
    License       string          `json:"license"`
    UpgradeCode   string          `json:"upgrade-code"`
    Files         WixFiles        `json:"files"`
    Directories   []string        `json:"directories"`
    Env           []WixEnvList    `json:"env"`
    Shortcuts     []WixShortcuts  `json:"shortcuts"`
}

type WixFiles {
    Guid      string      `json:"guid"`
    Dir       string      `json:"dir"`
    Items     []string    `json:"items"`
}

type WixEnvList {
    Guid       string      `json:"guid"`
    Vars       []WixEnv    `json:"vars"`
}
type WixEnv {
    Name       string      `json:"name"`
    Value      string      `json:"value"`
    Permanent  string      `json:"permanent"`
    System     string      `json:"system"`
    Action     string      `json:"action"`
    Part       string      `json:"part"`
}
type WixShortcuts {
    Dir       string        `json:"guid"`
    Items     []WixShortcut `json:"name"`
}
type WixShortcut {
    Guid          string      `json:"guid"`
    Name          string      `json:"name"`
    Description   string      `json:"description"`
    Target        string      `json:"target"`
    WDir          string      `json:"wdir"`
    Arguments     string      `json:"arguments"`
}

func main () {

  var wixFile WixManifest
  dat, err := ioutil.ReadFile("wix.json")
  if err != nil {
    panic(err)
  }
  err = json.Unmarshal(dat, &wixFile);

  wixFile, err := LoadWixManifest("")
  if err != nil {
    panic(err)
  }

  updated, err := CheckWixManifest(wixFile)
  if err != nil {
    panic(err)
  }
  if updated {
    err = WriteWixManifest(wixFile, "wix.json")
    if err != nil {
      panic(err)
    }
  }

  wixFile.License, err = RewriteLicenseFile(wixFile.License, "builder")
  if err != nil {
    panic(err)
  }

  RewriteWixFilePaths(wixFile, "builder")

  err = GenerateTemplates(wixFile, "templates", "builder")
  if err != nil {
    panic(err)
  }

  cmdStr = GenerateCmd(wixFile, "builder", "test.msi")
  if err != nil {
    panic(err)
  }

  fmt.Println(cmdStr)
  fmt.Println("All done!")
}

func RewriteWixFilePaths(wixFile *wixFile, o string) {
  for i, file := range wixFile.Files.Items {
    wixFile.Files.Items[i], _ := filepath.Rel(o, file)
  }
  for i, d := range wixFile.Directories {
    wixFile.Directories[i], _ := filepath.Rel(o, d)
  }
}

func GenerateCmd (wixFile *WixManifest, odir string, output string) string {
  targetFile := filepath.Join(odir, "build.cmd")

  cmd := ""

  for i, dir := range wixFile.Directories {
    cmd += "heat dir "+dir+" -nologo -cg AppFiles"+i+" -gg -g1 -srd -sfrag -template fragment -dr APPDIR"+i+" -var var.SourceDir"+i+" -out "+odir+"/AppFiles"+i+".wxs"
  }
  cmd += "\n"
  cmd += "candle -o "+odir+"\\"
  for i, dir := range wixFile.Directories {
    cmd += " -dSourceDir"+i+"="+dir
  }
  for i, dir := range wixFile.Directories {
    cmd += " "+odir+"\\AppFiles"+i+".wxs"
  }
  cmd += " "+odir+"\\LicenseAgreementDlg_HK.wxs "+odir+"\\WixUI_HK.wxs "+odir+"\\product.wxs"
  cmd += "\n"
  cmd += "light -ext WixUIExtension -ext WixUtilExtension -sacl "
  cmd += " -out "+output
  for i, dir := range wixFile.Directories {
    cmd += " "+odir+"\\AppFiles"+i+".wixobj"
  }
  cmd += " "+odir+"\\LicenseAgreementDlg_HK.wixobj "+odir+"\\WixUI_HK.wixobj "+odir+"\\product.wixobj"
  cmd += "\n"
  cmd += "@pause"
  cmd += "\n"

  return cmd
}

func RewriteLicenseFile (s string, odir string) (string, error) {
  f, err := os.Open(s)
  if err != nil {
      return s, err
  }
  defer f.Close()

  targetFile := filepath.Join(odir, "license.rtf")
  out, err := os.Create(targetFile)
  if err != nil {
      return s, err
  }
  defer out.Close()

  wInUTF8 := transform.NewWriter(out, charmap.Windows1252.NewEncoder())

  _, err = io.Copy(wInUTF8, f)
  if err != nil {
      return s, err
  }

  return targetFile, nil
}

func GenerateTemplates (wixFile *WixManifest, p string, o string) error {
  os.MkdirAll(o, 0644)

  var err error

  product := template.ParseFiles(filepath.Join(p, "product.wxs"))
  licenseDlg := template.ParseFiles(filepath.Join(p, "LicenseAgreementDlg_HK.wxs"))
  ui := template.ParseFiles(filepath.Join(p, "WixUI_HK.wxs"))

  wProduct := os.Create(filepath.Join(o, "product.wxs"))
  defer wProduct.Close()
  err = t.Execute(wProduct, wixFile)
  if err!=nil {
    return error
  }

  wLicenseDlg := os.Create(filepath.Join(o, "LicenseAgreementDlg_HK.wxs"))
  defer wLicenseDlg.Close()
  err = t.Execute(wLicenseDlg, wixFile)
  if err!=nil {
    return error
  }

  wUi := os.Create(filepath.Join(o, "WixUI_HK.wxs"))
  defer wUi.Close()
  err = t.Execute(wUi, wixFile)
  if err!=nil {
    return error
  }

  return nil
}

func WriteWixManifest(wixFile *WixManifest, p string) error {
  byt, err := json.MarshalIndent(wixFile, "", "  ")
  if err != nil {
    return err
  }
  err = ioutil.WriteFile(p, byt, 0644)
  if err != nil {
    return err
  }
  return nil
}

func LoadWixManifest (p string) (*WixManifest, err){
  if p=="" {
    p = "wix.json"
  }
  var wixFile WixManifest
  dat, err := ioutil.ReadFile(p)
  if err != nil {
    return nil, err
  }
  err = json.Unmarshal(dat, &wixFile);
  if  err != nil {
    return nil, err
  }
  return &wixFile, err
}

func CheckWixManifest (wixFile *WixManifest) (bool, error) {
  var err error
  updated := false
  if wixFile.UpgradeCode=="" {
    wixFile.UpgradeCode, err = MakeNewGuid()
    if err!=nil {
      return false, err
    }
    updated = true
  }
  if wixFile.Files!=nil {
    if wixFile.Files.Guid=="" {
      wixFile.Files.Guid, err = MakeNewGuid()
      if err!=nil {
        return false, err
      }
      updated = true
    }
  }
  for _, env := range wixFile.Env {
    if env.Guid=="" && len(env.Vars)>0 {
      env.Guid, err = MakeNewGuid()
      if err!=nil {
        return false, err
      }
      updated = true
    }
  }
  for _, shortcut := range wixFile.Shortcuts.Items {
    if shortcut.Guid=="" {
      shortcut.Guid, err = MakeNewGuid()
      if err!=nil {
        return false, err
      }
      updated = true
    }
  }
  return updated, nil
}

func MakeNewGuid () (string, error) {
  if runtime.GOOS == "windows" {
    cmd := "cscript.exe"
  	args := []string{filepath.Join(filepath.Base(os.Args[0]), "utils", "myuuid.vbs")}
    out, err := exec.Command(cmd, args...).CombinedOutput();
  	if  err != nil {
  		return "", err
  	}
		return string(out), nil
  } else {
    cmd := "uuidgen"
  	args := []string{"-t"}
    out, err := exec.Command(cmd, args...).CombinedOutput();
  	if  err != nil {
  		return "", err
  	}
		return string(out), nil
  }
}
