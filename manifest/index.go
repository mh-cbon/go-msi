package manifest

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/satori/go.uuid"
)

// WixManifest is the struct to decode a wix.json file.
type WixManifest struct {
	Product        string       `json:"product"`
	Company        string       `json:"company"`
	Version        string       `json:"version,omitempty"`
	VersionOk      string       `json:"-"`
	License        string       `json:"license,omitempty"`
	UpgradeCode    string       `json:"upgrade-code"`
	Files          WixFiles     `json:"files,omitempty"`
	Directories    []string     `json:"directories,omitempty"`
	RelDirs        []string     `json:"-"`
	Env            WixEnvList   `json:"env,omitempty"`
	Shortcuts      WixShortcuts `json:"shortcuts,omitempty"`
	Choco          ChocoSpec    `json:"choco,omitempty"`
	Hooks          []Hook       `json:"hooks,omitempty"`
	InstallHooks   []Hook       `json:"-"`
	UninstallHooks []Hook       `json:"-"`
}

// ChocoSpec is the struct to decode the choco key of a wix.json file.
type ChocoSpec struct {
	ID             string `json:"id,omitempty"`
	Title          string `json:"title,omitempty"`
	Authors        string `json:"authors,omitempty"`
	Owners         string `json:"owners,omitempty"`
	Description    string `json:"description,omitempty"`
	ProjectURL     string `json:"project-url,omitempty"`
	Tags           string `json:"tags,omitempty"`
	LicenseURL     string `json:"license-url,omitempty"`
	IconURL        string `json:"icon-url,omitempty"`
	RequireLicense bool   `json:"require-license,omitempty"`
	MsiFile        string `json:"-"`
	MsiSum         string `json:"-"`
	BuildDir       string `json:"-"`
	ChangeLog      string `json:"-"`
}

const (
	whenInstall   = "install"
	whenUninstall = "uninstall"
)

// HookPhases describes known hook phases.
var HookPhases = map[string]bool{
	whenInstall:   true,
	whenUninstall: true,
}

// Hook describes a command to run on install / uninstall.
type Hook struct {
	Command       string `json:"command,omitempty"`
	CookedCommand string `json:"-"`
	When          string `json:"when,omitempty"`
}

// WixFiles is the struct to decode files key of the wix.json file.
type WixFiles struct {
	GUID  string   `json:"guid"`
	Items []string `json:"items"`
}

// WixEnvList is the struct to decode env key of the wix.json file.
type WixEnvList struct {
	GUID string   `json:"guid"`
	Vars []WixEnv `json:"vars"`
}

// WixEnv is the struct to decode env value of the wix.json file.
type WixEnv struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	Permanent string `json:"permanent"`
	System    string `json:"system"`
	Action    string `json:"action"`
	Part      string `json:"part"`
}

// WixShortcuts is the struct to decode shortcuts key of the wix.json file.
type WixShortcuts struct {
	GUID  string        `json:"guid,omitempty"`
	Items []WixShortcut `json:"items,omitempty"`
}

// WixShortcut is the struct to decode shortcut value of the wix.json file.
type WixShortcut struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Target      string `json:"target"`
	WDir        string `json:"wdir"`
	Arguments   string `json:"arguments"`
	Icon        string `json:"icon"` // a path to the ico file, no space in it.
}

// Write the manifest to the given file,
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
		return fmt.Errorf("JSON ReadFile failed with %v", err)
	}
	err = json.Unmarshal(dat, &wixFile)
	if err != nil {
		return fmt.Errorf("JSON Unmarshal failed with %v", err)
	}
	return nil
}

//SetGuids generates and apply guid values appropriately
func (wixFile *WixManifest) SetGuids(force bool) (bool, error) {
	updated := false
	if wixFile.UpgradeCode == "" || force {
		wixFile.UpgradeCode = uuid.NewV4().String()
		updated = true
	}
	if wixFile.Files.GUID == "" || force {
		wixFile.Files.GUID = uuid.NewV4().String()
		updated = true
	}
	if (wixFile.Env.GUID == "" || force) && len(wixFile.Env.Vars) > 0 {
		wixFile.Env.GUID = uuid.NewV4().String()
		updated = true
	}
	if (wixFile.Shortcuts.GUID == "" || force) && len(wixFile.Shortcuts.Items) > 0 {
		wixFile.Shortcuts.GUID = uuid.NewV4().String()
		updated = true
	}
	return updated, nil
}

// NeedGUID tells if the manifest json file is missing guid values.
func (wixFile *WixManifest) NeedGUID() bool {
	need := false
	if wixFile.UpgradeCode == "" {
		need = true
	}
	if wixFile.Files.GUID == "" {
		need = true
	}
	if wixFile.Env.GUID == "" && len(wixFile.Env.Vars) > 0 {
		need = true
	}
	if wixFile.Shortcuts.GUID == "" && len(wixFile.Shortcuts.Items) > 0 {
		need = true
	}
	return need
}

// RewriteFilePaths Reads Files and Directories of the wix.json file
// and turn their values into a relative path to out
// where out is the path to the wix templates files.
func (wixFile *WixManifest) RewriteFilePaths(out string) error {
	var err error
	out, err = filepath.Abs(out)
	if err != nil {
		return err
	}
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
		if s.Icon != "" {
			file, err := filepath.Abs(s.Icon)
			if err != nil {
				return err
			}
			wixFile.Shortcuts.Items[i].Icon, err = filepath.Rel(out, file)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Normalize Appropriately fixes some values within the decoded json
// It applies defaults values on the wix/msi property to
// to generate the msi package.
// It applies defaults values on the choco property to
// generate a nuget package
func (wixFile *WixManifest) Normalize() error {
	// Wix version Field of Product element
	// does not support semver strings
	// it supports only something like x.x.x.x
	// So, if the version has metadata/prerelease values,
	// lets get ride of those and save the workable version
	// into VersionOk field
	wixFile.VersionOk = wixFile.Version
	v, err := semver.NewVersion(wixFile.Version)
	if err != nil {
		return fmt.Errorf("Failed to parse version '%v': %v", wixFile.Version, err)
	}
	okVersion := ""
	okVersion += strconv.FormatInt(v.Major(), 10)
	okVersion += "." + strconv.FormatInt(v.Minor(), 10)
	okVersion += "." + strconv.FormatInt(v.Patch(), 10)
	wixFile.VersionOk = okVersion

	// choco fix
	if wixFile.Choco.ID == "" {
		wixFile.Choco.ID = wixFile.Product
	}
	if wixFile.Choco.Title == "" {
		wixFile.Choco.Title = wixFile.Product
	}
	if wixFile.Choco.Authors == "" {
		wixFile.Choco.Authors = wixFile.Company
	}
	if wixFile.Choco.Owners == "" {
		wixFile.Choco.Owners = wixFile.Company
	}
	if wixFile.Choco.Description == "" {
		wixFile.Choco.Description = wixFile.Product
	}
	wixFile.Choco.Tags += " admin" // required to pass chocolatey validation..

	// Escape hook commands and ensure the command name is enclosed in quotes (needed by wix)
	for i, hook := range wixFile.Hooks {
		cmd := strings.Trim(hook.Command, " ")
		if len(cmd) > 0 && cmd[0] != '"' {
			words := strings.Split(cmd, " ")
			cmd = `"` + words[0] + `"` + cmd[len(words[0]):]
		}
		buf := &bytes.Buffer{}
		if err := xml.EscapeText(buf, []byte(cmd)); err != nil {
			return err
		}
		wixFile.Hooks[i].CookedCommand = buf.String()
	}

	// Separate install and uninstall hooks to simplify templating
	for _, hook := range wixFile.Hooks {
		switch hook.When {
		case whenInstall:
			wixFile.InstallHooks = append(wixFile.InstallHooks, hook)
		case whenUninstall:
			wixFile.UninstallHooks = append(wixFile.UninstallHooks, hook)
		}
	}

	return nil
}
