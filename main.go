package main

import (
  "encoding/json"
  "io/ioutil"
  "path/filepath"
  "text/template"
  "io"
  "os"
  "fmt"
  "strings"
  "strconv"
  "runtime"
  "os/exec"

  "golang.org/x/text/encoding/charmap"
  "golang.org/x/text/transform"
  "github.com/mh-cbon/go-msi/manifest"
  "github.com/mh-cbon/go-msi/tpls"
  "github.com/mh-cbon/go-msi/rtf"
)

var VERSION = "0.0.0"

func main() {

	app := cli.NewApp()
	app.Name = "go-msi"
	app.Version = VERSION
	app.Usage = "Easy msi pakage for Go"
	app.UsageText = "go-msi <cmd> <options>"
	app.Commands = []cli.Command{
		{
			Name:   "check-json",
			Usage:  "Check the JSON wix manifest",
			Action: checkJson,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "path, p",
					Value: "wix.json",
					Usage: "Path to the wix manifest file",
				},
			},
		},
		{
			Name:   "set-guid",
			Usage:  "Sets appropriate guids in your wix manifest",
			Action: setGuid,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "path, p",
					Value: "wix.json",
					Usage: "Path to the wix manifest file",
				},
			},
		},
		{
			Name:   "generate-templates",
			Usage:  "Generate wix templates",
			Action: generateTemplates,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "path, p",
					Value: "wix.json",
					Usage: "Path to the wix manifest file",
				},
				cli.StringFlag{
					Name:  "src, s",
					Value: "templates",
					Usage: "Diretory path to the wix templates files",
				},
				cli.StringFlag{
					Name:  "out, o",
					Value: "builder",
					Usage: "Directory path to the generated wix templates files",
				},
				cli.StringFlag{
					Name:  "version",
					Value: "",
					Usage: "The version of your program",
				},
				cli.StringFlag{
					Name:  "license, l",
					Value: "",
					Usage: "Path to the license file",
				},
			},
		},
		{
			Name:   "to-windows",
			Usage:  "Write Windows1252 encoded file",
			Action: toWindows1252,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "src, s",
					Value: "",
					Usage: "Path to an UTF-8 encoded file",
				},
				cli.StringFlag{
					Name:  "out, o",
					Value: "",
					Usage: "Path to the ANSI generated file",
				},
			},
		},
		{
			Name:   "to-rtf",
			Usage:  "Write RTF formatted file",
			Action: toRtf,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "src, s",
					Value: "",
					Usage: "Path to a text file",
				},
				cli.StringFlag{
					Name:  "out, o",
					Value: "",
					Usage: "Path to the RTF generated file",
				},
				cli.BoolFlag{
					Name:  "reencode, e",
					Usage: "Also re encode UTF-8 to Windows1252 charset",
				},
			},
		},
		{
			Name:   "gen-wix-cmd",
			Usage:  "Generate a batch file of Wix commands to run",
			Action: generateWixCommands,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "path, p",
					Value: "wix.json",
					Usage: "Path to the wix manifest file",
				},
				cli.StringFlag{
					Name:  "src, s",
					Value: "templates",
					Usage: "Diretory path to the wix templates files",
				},
				cli.StringFlag{
					Name:  "out, o",
					Value: "builder",
					Usage: "Directory path to the generated wix cmd file",
				},
				cli.StringFlag{
					Name:  "arch, a",
					Value: "",
					Usage: "A target architecture , x64 or x86 (ia64 is not handled)",
				},
				cli.StringFlag{
					Name:  "msi, m",
					Value: "",
					Usage: "Path to write resulting msi file to",
				},
			},
		},
		{
			Name:   "run-wix-cmd",
			Usage:  "Run the batch file of Wix commands",
			Action: runWixCommands,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "out, o",
					Value: "builder",
					Usage: "Directory path to the generated wix cmd file",
				},
			},
		},
		{
			Name:   "make",
			Usage:  "All-in-one command to make MSI files",
			Action: make,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "out, o",
					Value: "builder",
					Usage: "Directory path to the generated wix cmd file",
				},
			},
		},
	}

	app.Run(os.Args)
}

func checkJson(c *cli.Context) error {
  path := c.String("path")

  var wixFile *manifest.WixManifest
  err = wixFile.Load(path)
  if err!=nil {
    return cli.NewExitError(err.Error(), 1)
  }

  fmt.Println("The manifest is syntaxically correct !")

  if wixFile.NeedGuid() {
    fmt.Println("The manifest needs Guid")
    fmt.Println("To update your file automatically run:")
    fmt.Println("     go-msi set-guid")
    return cli.NewExitError("Incomplete manifest file detected", 1)
  }
  return nil
}

func setGuid(c *cli.Context) error {
  path := c.String("path")

  var wixFile *manifest.WixManifest
  err = wixFile.Load(path)
  if err!=nil {
    return cli.NewExitError(err.Error(), 1)
  }

  updated, err := wixFile.SetGuids()
  if err!=nil {
    return cli.NewExitError(err.Error(), 1)
  }

  if updated {
    fmt.Println("The manifest was updated")
  } else {
    fmt.Println("The manifest was not updated")
  }

  err := wixFile.Write(path)
  if err!=nil {
    return cli.NewExitError(err.Error(), 1)
  }
  fmt.Println("The file is saved on disk")

  return nil
}

func generateTemplates(c *cli.Context) error {
  path := c.String("path")
  src := c.String("src")
  out := c.String("out")
  version := c.String("version")
  license := c.String("license")

  var wixFile *manifest.WixManifest
  err = wixFile.Load(path)
  if err!=nil {
    return cli.NewExitError(err.Error(), 1)
  }

  if wixFile.NeedGuid() {
    fmt.Println("The manifest needs Guid")
    fmt.Println("To update your file automatically run:")
    fmt.Println("     go-msi set-guid")
    return cli.NewExitError("Cannot proceed, manifest file is incomplete", 1)
  }

  wixFile.RewriteFilePaths(out)

  if version!="" {
    wixFile.Version = version
  }

  if license!="" {
    wixFile.License = license
  }

  templates, err := tpls.Find(src)
  if err!=nil {
    return cli.NewExitError(err.Error(), 1)
  }
  if len(templates)==0 {
    return cli.NewExitError("No templates *.wxs found in this directory", 1)
  }

  err = os.MkdirAll(out, 0744)
  if err!=nil {
    return cli.NewExitError(err.Error(), 1)
  }

  for _, tpl := range templates {
    dst := filepath.Join(out, filepath.Base(tpl))
    err = tpls.GenerateTemplate(wixFile, tpl, dst)
    if err!=nil {
      return cli.NewExitError(err.Error(), 1)
    }
  }

  fmt.Printf("Generated %s templates\n", len(templates))
  for _, tpl := range templates {
    dst := filepath.Join(out, filepath.Base(tpl))
    fmt.Println("- %s", dst)
  }

  return nil
}

func toWindows1252(c *cli.Context) error {
  src := c.String("src")
  out := c.String("out")

  if src=="" {
    return cli.NewExitError("--src argument is required", 1)
  }
  if out=="" {
    return cli.NewExitError("--out argument is required", 1)
  }
  if _, err := os.Stat(src); os.IsNotExist(err) {
    return cli.NewExitError(err.Error(), 1)
  }
  os.MkdirAll(filepath.Dir(out), 0744)
  err := rtf.WriteAsWindows1252(src, out)
  if err!=nil {
    return cli.NewExitError(err.Error(), 1)
  }
  return nil
}

func toRtf(c *cli.Context) error {
  src := c.String("src")
  out := c.String("out")
  reencode := c.Bool("reencode")

  if src=="" {
    return cli.NewExitError("--src argument is required", 1)
  }
  if out=="" {
    return cli.NewExitError("--out argument is required", 1)
  }
  if _, err := os.Stat(src); os.IsNotExist(err) {
    return cli.NewExitError(err.Error(), 1)
  }

  os.MkdirAll(filepath.Dir(out), 0744)

  err := rtf.WriteAsRtf(src, out, reencode)
  if err!=nil {
    return cli.NewExitError(err.Error(), 1)
  }

  return nil
}

func generateWixCommands(c *cli.Context) error {
  path := c.String("path")
  src := c.String("src")
  out := c.String("out")
  msi := c.String("msi")
  arch := c.String("arch")

  templates, err := tpls.Find(src)
  if err!=nil {
    return cli.NewExitError(err.Error(), 1)
  }
  if len(templates)==0 {
    return cli.NewExitError("No templates *.wxs found in this directory", 1)
  }

  var wixFile *manifest.WixManifest
  err = wixFile.Load(path)
  if err!=nil {
    return cli.NewExitError(err.Error(), 1)
  }

  if wixFile.NeedGuid() {
    fmt.Println("The manifest needs Guid")
    fmt.Println("To update your file automatically run:")
    fmt.Println("     go-msi set-guid")
    return cli.NewExitError("Cannot proceed, manifest file is incomplete", 1)
  }

  wixFile.RewriteFilePaths(out)

  cmdStr := GenerateCmd(wixFile, out, msi, arch)

  targetFile := filepath.Join(out, "build.bat")
  err := ioutil.WriteFile(targetFile, []byte(cmdStr), 0644)
  if err != nil {
    return cli.NewExitError(err.Error(), 1)
  }

  return nil
}

func runWixCommands(c *cli.Context) error {
  out := c.String("out")

  bin, err := exec.LookPath("cmd.exe")
  if err != nil {
    return cli.NewExitError(err.Error(), 1)
  }
  args := []string{"/C", "build.bat"}
  oCmd := exec.Command(bin, args...)
  oCmd.Dir = out
  oCmd.Stdout = os.Stdout
  oCmd.Stderr = os.Stderr
  err = oCmd.Run();
  if  err != nil {
    return cli.NewExitError(err.Error(), 1)
  }

  return nil
}

func make(c *cli.Context) error {
  path := c.String("path")
  src := c.String("src")
  out := c.String("out")
  version := c.String("version")
  license := c.String("license")
  msi := c.String("msi")
  arch := c.String("arch")

  var wixFile *manifest.WixManifest
  err = wixFile.Load(path)
  if err!=nil {
    return cli.NewExitError(err.Error(), 1)
  }

  if wixFile.NeedGuid() {
    _, err := wixFile.SetGuids()
    if err!=nil {
      return cli.NewExitError(err.Error(), 1)
    }
  }

  wixFile.RewriteFilePaths(out)

  if version!="" {
    wixFile.Version = version
  }

  if license!="" {
    wixFile.License = license
  }

  if wixFile.License!="" {
    if !rtf.IsRtf(wixFile.License) {
      target := filepath.Join(out, filepath.Base(wixFile.License)+".rtf")
      err := rtf.WriteAsRtf(wixFile.License, target, true)
      if err!=nil {
        return cli.NewExitError(err.Error(), 1)
      }
      wixFile.License = target
    }
  }

  templates, err := tpls.Find(src)
  if err!=nil {
    return cli.NewExitError(err.Error(), 1)
  }
  if len(templates)==0 {
    return cli.NewExitError("No templates *.wxs found in this directory", 1)
  }

  err = os.MkdirAll(out, 0744)
  if err!=nil {
    return cli.NewExitError(err.Error(), 1)
  }

  for _, tpl := range templates {
    dst := filepath.Join(out, filepath.Base(tpl))
    err = tpls.GenerateTemplate(wixFile, tpl, dst)
    if err!=nil {
      return cli.NewExitError(err.Error(), 1)
    }
  }

  cmdStr := GenerateCmd(wixFile, out, msi, arch)

  targetFile := filepath.Join(out, "build.bat")
  err := ioutil.WriteFile(targetFile, []byte(cmdStr), 0644)
  if err != nil {
    return cli.NewExitError(err.Error(), 1)
  }

  bin, err := exec.LookPath("cmd.exe")
  if err != nil {
    return cli.NewExitError(err.Error(), 1)
  }
  args := []string{"/C", "build.bat"}
  oCmd := exec.Command(bin, args...)
  oCmd.Dir = out
  oCmd.Stdout = os.Stdout
  oCmd.Stderr = os.Stderr
  err = oCmd.Run();
  if  err != nil {
    return cli.NewExitError(err.Error(), 1)
  }

  return nil
}
