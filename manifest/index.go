package manifest

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/satori/go.uuid"
)

// WixManifest is the struct to decode a wix.json file.
type WixManifest struct {
	Product        string         `json:"product"`
	Company        string         `json:"company"`
	Version        string         `json:"version,omitempty"`
	VersionOk      string         `json:"-"`
	License        string         `json:"license,omitempty"`
	Banner         string         `json:"banner,omitempty"`
	Dialog         string         `json:"dialog,omitempty"`
	Icon           string         `json:"icon,omitempty"`
	UpgradeCode    string         `json:"upgrade-code"`
	Files          []File         `json:"files,omitempty"`
	Directories    []string       `json:"directories,omitempty"`
	DirNames       []string       `json:"-"`
	RelDirs        []string       `json:"-"`
	Env            WixEnvList     `json:"env"`
	Registries     []RegistryItem `json:"registries,omitempty"`
	Shortcuts      WixShortcuts   `json:"shortcuts"`
	Choco          ChocoSpec      `json:"choco"`
	Hooks          []Hook         `json:"hooks,omitempty"`
	InstallHooks   []Hook         `json:"-"`
	UninstallHooks []Hook         `json:"-"`
	Properties     []Property     `json:"properties,omitempty"`
	Conditions     []Condition    `json:"conditions,omitempty"`
}

// File is the struct to decode a file.
type File struct {
	Path    string   `json:"path"`
	Service *Service `json:"service,omitempty"`
}

// Service is the struct to decode a service.
type Service struct {
	Name        string `json:"name"`
	Bin         string `json:"-"`
	Start       string `json:"start"`
	DisplayName string `json:"display-name,omitempty"`
	Description string `json:"description,omitempty"`
	Arguments   string `json:"arguments,omitempty"`
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
	Return        string `json:"return,omitempty"`
}

// Property describes a property to initialize.
type Property struct {
	ID       string    `json:"id"`
	Registry *Registry `json:"registry,omitempty"`
	Value    *Value    `json:"value,omitempty"`
}

// Registry describes a registry entry.
type Registry struct {
	Path string `json:"path"`
	Root string `json:"-"`
	Key  string `json:"-"`
	Name string `json:"name,omitempty"`
}

// Value describes a simple string value
type Value string

// Condition describes a condition to check before installation.
type Condition struct {
	Condition string `json:"condition"`
	Message   string `json:"message"`
}

// WixEnvList is the struct to decode env key of the wix.json file.
type WixEnvList struct {
	GUID string   `json:"guid,omitempty"`
	Vars []WixEnv `json:"vars,omitempty"`
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

// RegistryItem is the struct to decode a registry item.
type RegistryItem struct {
	Registry
	GUID   string          `json:"guid,omitempty"`
	Values []RegistryValue `json:"values,omitempty"`
}

// RegistryValue is the struct to decode a registry value.
type RegistryValue struct {
	Name  string `json:"name"`
	Type  string `json:"type,omitempty"` // string (default if omitted), integer, ...
	Value string `json:"value"`
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
func (wixFile *WixManifest) SetGuids(force bool) bool {
	updated := false
	if wixFile.UpgradeCode == "" || force {
		wixFile.UpgradeCode = makeGUID()
		updated = true
	}
	if (wixFile.Env.GUID == "" || force) && len(wixFile.Env.Vars) > 0 {
		wixFile.Env.GUID = makeGUID()
		updated = true
	}
	for i, r := range wixFile.Registries {
		if r.GUID == "" || force {
			wixFile.Registries[i].GUID = makeGUID()
			updated = true
		}
	}
	if (wixFile.Shortcuts.GUID == "" || force) && len(wixFile.Shortcuts.Items) > 0 {
		wixFile.Shortcuts.GUID = makeGUID()
		updated = true
	}
	return updated
}

func makeGUID() string {
	return strings.ToUpper(uuid.NewV4().String())
}

// NeedGUID tells if the manifest json file is missing guid values.
func (wixFile *WixManifest) NeedGUID() bool {
	if wixFile.UpgradeCode == "" {
		return true
	}
	if wixFile.Env.GUID == "" && len(wixFile.Env.Vars) > 0 {
		return true
	}
	for _, r := range wixFile.Registries {
		if r.GUID == "" {
			return true
		}
	}
	if wixFile.Shortcuts.GUID == "" && len(wixFile.Shortcuts.Items) > 0 {
		return true
	}
	return false
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
	for i, file := range wixFile.Files {
		path, err := rewrite(out, file.Path)
		if err != nil {
			return err
		}
		wixFile.Files[i].Path = path
	}
	for _, dir := range wixFile.Directories {
		wixFile.DirNames = append(wixFile.DirNames, filepath.Base(dir))
		path, err := rewrite(out, dir)
		if err != nil {
			return err
		}
		wixFile.RelDirs = append(wixFile.RelDirs, path)
	}
	for i, s := range wixFile.Shortcuts.Items {
		if s.Icon != "" {
			path, err := rewrite(out, s.Icon)
			if err != nil {
				return err
			}
			wixFile.Shortcuts.Items[i].Icon = path
		}
	}
	return nil
}

func rewrite(out, path string) (string, error) {
	var err error
	path, err = filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return filepath.Rel(out, path)
}

// Normalize appropriately fixes some values within the decoded json.
// It applies defaults values on the wix/msi property generate the msi package.
// It applies defaults values on the choco property to generate a nuget package.
func (wixFile *WixManifest) Normalize() error {
	wixFile.VersionOk = wixFile.Version
	// Wix version Field of Product element
	// does not support semver strings
	// it supports only something like x.x.x.x
	// So, if the version has metadata/prerelease values,
	// lets get ride of those and save the workable version
	// into VersionOk field
	v, err := semver.NewVersion(wixFile.Version)
	if err == nil {
		wixFile.VersionOk = v.String()
	}

	if wixFile.Banner != "" {
		path, err := filepath.Abs(wixFile.Banner)
		if err != nil {
			return err
		}
		wixFile.Banner = path
	}
	if wixFile.Dialog != "" {
		path, err := filepath.Abs(wixFile.Dialog)
		if err != nil {
			return err
		}
		wixFile.Dialog = path
	}
	if wixFile.Icon != "" {
		path, err := filepath.Abs(wixFile.Icon)
		if err != nil {
			return err
		}
		wixFile.Icon = path
	}

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

	// Split registry path into root and key
	for _, prop := range wixFile.Properties {
		reg := prop.Registry
		if reg != nil {
			if reg.Root, reg.Key, err = extractRegistry(reg.Path); err != nil {
				return err
			}
		}
	}
	for i := range wixFile.Registries {
		r := &wixFile.Registries[i]
		if r.Root, r.Key, err = extractRegistry(r.Path); err != nil {
			return err
		}
		for j := range r.Values {
			v := &r.Values[j]
			if v.Type == "" {
				v.Type = "string"
			}
		}
	}

	// Bind services to their file component
	for i, file := range wixFile.Files {
		if file.Service != nil {
			wixFile.Files[i].Service.Bin = filepath.Base(file.Path)
		}
	}

	return nil
}

func extractRegistry(path string) (string, string, error) {
	p := strings.Split(path, `\`)
	if len(p) < 2 {
		return "", "", fmt.Errorf("invalid registry path %q", p)
	}
	return p[0], strings.Join(p[1:len(p)], `\`), nil
}
