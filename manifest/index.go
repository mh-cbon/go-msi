package manifest

import (
  "encoding/json"
  "io/ioutil"

  "github.com/mh-cbon/go-msi/guid"
)

type WixManifest struct {
    Product       string          `json:"product"`
    Company       string          `json:"company"`
    Version       string          `json:"version"`
    License       string          `json:"license"`
    UpgradeCode   string          `json:"upgrade-code"`
    Files         WixFiles        `json:"files"`
    Directories   []string        `json:"directories"`
    RelDirs       []string        `json:"-"`
    Env           WixEnvList      `json:"env"`
    Shortcuts     WixShortcuts    `json:"shortcuts"`
}

type WixFiles struct {
    Guid      string      `json:"guid"`
    Items     []string    `json:"items"`
}

type WixEnvList struct {
    Guid       string      `json:"guid"`
    Vars       []WixEnv    `json:"vars"`
}
type WixEnv struct{
    Name       string      `json:"name"`
    Value      string      `json:"value"`
    Permanent  string      `json:"permanent"`
    System     string      `json:"system"`
    Action     string      `json:"action"`
    Part       string      `json:"part"`
}
type WixShortcuts struct{
    Guid      string        `json:"guid"`
    Items     []WixShortcut `json:"items"`
}
type WixShortcut struct{
    Name          string      `json:"name"`
    Description   string      `json:"description"`
    Target        string      `json:"target"`
    WDir          string      `json:"wdir"`
    Arguments     string      `json:"arguments"`
}

func (wixFile *WixManifest) Write(p string) error {
  if p=="" {
    p = "wix.json"
  }
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

func (wixFile *WixManifest) Load(p string) error {
  if p=="" {
    p = "wix.json"
  }
  if _, err := os.Stat(p); os.IsNotExist(err) {
    return err
  }
  dat, err := ioutil.ReadFile(p)
  if err != nil {
    return err
  }
  err = json.Unmarshal(dat, wixFile);
  if  err != nil {
    return err
  }
  return err
}

func (wixFile *WixManifest) SetGuids () (bool, error) {
  var err error
  updated := false
  if wixFile.UpgradeCode=="" {
    wixFile.UpgradeCode, err = guid.Make()
    if err!=nil {
      return false, err
    }
    updated = true
  }
  if wixFile.Files.Guid=="" {
    wixFile.Files.Guid, err = guid.Make()
    if err!=nil {
      return false, err
    }
    updated = true
  }
  if wixFile.Env.Guid=="" && len(wixFile.Env.Vars)>0 {
    wixFile.Env.Guid, err = guid.Make()
    if err!=nil {
      return false, err
    }
    updated = true
  }
  if wixFile.Shortcuts.Guid=="" && len(wixFile.Shortcuts.Items)>0 {
    wixFile.Shortcuts.Guid, err = guid.Make()
    if err!=nil {
      return false, err
    }
    updated = true
  }
  return updated, nil
}

func (wixFile *WixManifest) NeedGuid () bool {
  var err error
  need := false
  if wixFile.UpgradeCode=="" {
    need = true
  }
  if wixFile.Files.Guid=="" {
    need = true
  }
  if wixFile.Env.Guid=="" && len(wixFile.Env.Vars)>0 {
    need = true
  }
  if wixFile.Shortcuts.Guid=="" && len(wixFile.Shortcuts.Items)>0 {
    need = true
  }
  return need
}

func (wixFile *WixManifest) RewriteFilePaths(o string) {
  for i, file := range wixFile.Files.Items {
    wixFile.Files.Items[i], _ = filepath.Rel(o, file)
  }
  for _, d := range wixFile.Directories {
    r, _ := filepath.Rel(o, d)
    wixFile.RelDirs = append(wixFile.RelDirs, r)
  }
}
