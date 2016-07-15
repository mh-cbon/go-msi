package manifest

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mh-cbon/go-msi/guid"
)

type WixManifest struct {
	Product     string       `json:"product"`
	Company     string       `json:"company"`
	Version     string       `json:"version,omitempty"`
	License     string       `json:"license,omitempty"`
	UpgradeCode string       `json:"upgrade-code"`
	Files       WixFiles     `json:"files,omitempty"`
	Directories []string     `json:"directories,omitempty"`
	RelDirs     []string     `json:"-"`
	Env         WixEnvList   `json:"env,omitempty"`
	Shortcuts   WixShortcuts `json:"shortcuts,omitempty"`
}

type WixFiles struct {
	Guid  string   `json:"guid"`
	Items []string `json:"items"`
}

type WixEnvList struct {
	Guid string   `json:"guid"`
	Vars []WixEnv `json:"vars"`
}
type WixEnv struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	Permanent string `json:"permanent"`
	System    string `json:"system"`
	Action    string `json:"action"`
	Part      string `json:"part"`
}
type WixShortcuts struct {
	Guid  string        `json:"guid,omitempty"`
	Items []WixShortcut `json:"items,omitempty"`
}
type WixShortcut struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Target      string `json:"target"`
	WDir        string `json:"wdir"`
	Arguments   string `json:"arguments"`
	Icon        string `json:"icon"`         // a path to the ico file, no space in it.
}

// Writes the manifest to the given file,
// if file is empty, writes to wix.json
func (wixFile *WixManifest) Write(p string) error {
	if p == "" {
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

// Load the manifest from given file path,
// if the file path is empty, reads from wix.json
func (wixFile *WixManifest) Load(p string) error {
	if p == "" {
		p = "wix.json"
	}
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return err
	}
	dat, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}
	err = json.Unmarshal(dat, &wixFile)
	if err != nil {
		return err
	}
	return nil
}

//SetGuids generates and apply guid values appropriately
func (wixFile *WixManifest) SetGuids() (bool, error) {
	var err error
	updated := false
	if wixFile.UpgradeCode == "" {
		wixFile.UpgradeCode, err = guid.Make()
		if err != nil {
			return false, err
		}
		updated = true
	}
	if wixFile.Files.Guid == "" {
		wixFile.Files.Guid, err = guid.Make()
		if err != nil {
			return false, err
		}
		updated = true
	}
	if wixFile.Env.Guid == "" && len(wixFile.Env.Vars) > 0 {
		wixFile.Env.Guid, err = guid.Make()
		if err != nil {
			return false, err
		}
		updated = true
	}
	if wixFile.Shortcuts.Guid == "" && len(wixFile.Shortcuts.Items) > 0 {
		wixFile.Shortcuts.Guid, err = guid.Make()
		if err != nil {
			return false, err
		}
		updated = true
	}
	return updated, nil
}

// Indicates if the manifest needs GUIDs to be set
func (wixFile *WixManifest) NeedGuid() bool {
	need := false
	if wixFile.UpgradeCode == "" {
		need = true
	}
	if wixFile.Files.Guid == "" {
		need = true
	}
	if wixFile.Env.Guid == "" && len(wixFile.Env.Vars) > 0 {
		need = true
	}
	if wixFile.Shortcuts.Guid == "" && len(wixFile.Shortcuts.Items) > 0 {
		need = true
	}
	return need
}

// Reads Files and Directories turn their values
// into a relative path to out(path to the wix templates files)
func (wixFile *WixManifest) RewriteFilePaths(out string) error {
	var err error
	for i, file := range wixFile.Files.Items {
		file, err = filepath.Abs(file)
		if err != nil {
			return err
		}
		wixFile.Files.Items[i], err = filepath.Rel(out, file)
		if err != nil {
			return err
		}
	}
	for _, d := range wixFile.Directories {
		d, err = filepath.Abs(d)
		if err != nil {
			return err
		}
		r, err := filepath.Rel(out, d)
		if err != nil {
			return err
		}
		wixFile.RelDirs = append(wixFile.RelDirs, r)
	}
	for i, s := range wixFile.Shortcuts.Items {
		file, err := filepath.Abs(s.Icon)
		if err != nil {
			return err
		}
		wixFile.Shortcuts.Items[i].Icon, err = filepath.Rel(out, file)
		if err != nil {
			return err
		}
	}
	return nil
}
