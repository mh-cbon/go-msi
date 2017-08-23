// +build windows

package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

func main() {

	svcName := "HelloSvc"

	confirm(rmFile("log-install.txt"), "install log removal")
	confirm(rmFile("log-uninstall.txt"), "uninstall log removal")

	wd := makeDir("c:/gopath/src/github.com/mh-cbon/go-msi/testing/hello")
	mustContains(wd, "hello.go")
	mustChdir(wd)

	mustEnvEq("$env:some", "")

	mustMkdirAll("build/amd64")
	helloBuild := makeCmd("go", "build", "-o", _p("build/amd64/hello.exe"), "hello.go")
	mustExec(helloBuild, "hello build failed %v")

	helloPkgSetup := makeCmd("C:/go-msi/go-msi.exe", "set-guid")
	mustExec(helloPkgSetup, "hello packaging setup failed %v")

	helloPkg := makeCmd("C:/go-msi/go-msi.exe", "make",
		"--msi", "hello.msi",
		"--version", "0.0.1",
		"--arch", "amd64",
		"--keep",
	)
	mustExec(helloPkg, "hello packaging failed %v")

	resultPackage := makeFile("hello.msi")
	mustExists(resultPackage, "Package file is missing %v")

	mustNotHaveWindowsService("HelloSvc")

	helloPackageInstall := makeCmd("msiexec", "/i", "hello.msi", "/q", "/log", "log-install.txt")
	mustExec(helloPackageInstall, "hello package install failed %v")
	readFile("log-install.txt")
	mustSucceed(rmFile("log-install.txt"), "rmfile failed %v")

	// mustShowEnv("$env:path")
	// mustEnvEq("$env:some", "value")

	readDir("C:/Program Files/hello")
	readDir("C:/Program Files/hello/assets")

	mgr, helloSvc := mustHaveWindowsService(svcName)
	mustHaveStartedWindowsService(svcName, helloSvc)

	helloEpURL := "http://localhost:8080/"
	// helloExecPath := "C:/Program Files/hello/hello.exe"
	// mustExecHello(helloExecPath, helloEpURL)
	mustQueryHello(helloEpURL)
	// mustStopWindowsService(svcName, helloSvc)
	mustSucceed(helloSvc.Close(), "Failed to close the service %v")
	mgr.Disconnect()

	helloPackageUninstall := makeCmd("msiexec", "/x", "hello.msi", "/q", "/log", "log-uninstall.txt")
	mustExec(helloPackageUninstall, "hello package uninstall failed %v")
	readFile("log-uninstall.txt")
	mustSucceed(rmFile("log-uninstall.txt"), "rmfile failed %v")

	mustNotHaveWindowsService("HelloSvc")

	// mustShowEnv("$env:path")
	// mustEnvEq("$env:some", "")

	helloChocoPkg := makeCmd("C:/go-msi/go-msi.exe", "choco",
		"--input", "hello.msi",
		"--version", "0.0.1",
		"-c", _qp("C:/Program Files/changelog/changelog.exe")+" ghrelease --version 0.0.1",
		"--keep",
	)
	mustExec(helloChocoPkg, "hello choco package make failed %v")

	helloNuPkg := makeFile("hello.0.0.1.nupkg")
	mustExists(helloNuPkg, "Chocolatey nupkg file is missing %v")

	mustNotHaveWindowsService("HelloSvc")

	helloChocoInstall := makeCmd("choco", "install", "hello.0.0.1.nupkg", "-y")
	mustExec(helloChocoInstall, "hello choco package install failed %v")

	readDir("C:/Program Files/hello")
	readDir("C:/Program Files/hello/assets")

	mgr, helloSvc = mustHaveWindowsService(svcName)
	mustHaveStartedWindowsService(svcName, helloSvc)
	mustSucceed(helloSvc.Close(), "Failed to close the service %v")
	mgr.Disconnect()

	// mustShowEnv("$env:path")
	// mustEnvEq("$env:some", "value")

	// mustExecHello(helloExecPath, helloEpURL)
	mustQueryHello(helloEpURL)
	// mustStopWindowsService(svcName, helloSvc)

	helloChocoUninstall := makeCmd("choco", "uninstall", "hello", "-v", "-d", "-y", "--force")
	mustExec(helloChocoUninstall, "hello choco package uninstall failed %v")
	readFile("C:\\ProgramData\\chocolatey\\logs\\chocolatey.log")

	mustNotHaveWindowsService("HelloSvc")

	// mustShowEnv("$env:path")
	// mustEnvEq("$env:some", "")

}

func mustHaveWindowsService(n string) (*mgr.Mgr, *mgr.Service) {
	mgr, err := mgr.Connect()
	mustSucceed(err, "Failed to connect to the service manager %v")
	s, err := mgr.OpenService(n)
	mustSucceed(err, "Failed to open the service %v")
	if s == nil {
		mustSucceed(err, "Failed to find the service %v")
	}
	log.Printf("SUCCESS: Service %q exists\n", n)
	return mgr, s
}

func mustNotHaveWindowsService(n string) bool {
	mgr, err := mgr.Connect()
	mustSucceed(err, "Failed to connect to the service manager %v")
	defer mgr.Disconnect()
	s, err := mgr.OpenService(n)
	mustNotSucceed(err, "Must fail to open the service %v")
	if s == nil {
		mustNotSucceed(err, "Must fail to find the service %v")
	} else {
		defer s.Close()
	}
	log.Printf("SUCCESS: Service %q does not exist\n", n)
	return s == nil
}

func mustHaveStartedWindowsService(n string, s *mgr.Service) {
	status, err := s.Query()
	mustSucceed(err, "Failed to query the service status %v")
	if status.State != svc.Running {
		mustSucceed(fmt.Errorf("Service not started %v", n))
	}
	log.Printf("SUCCESS: Service %q was started\n", n)
}

func mustStopWindowsService(n string, s *mgr.Service) {
	status, err := s.Control(svc.Stop)
	mustSucceed(err, "Failed to control the service status %v")
	if status.State != svc.Stopped {
		mustSucceed(fmt.Errorf("Service not stopped %v", n))
	}
	log.Printf("SUCCESS: Service %q was stopped\n", n)
}

func mustNotSucceed(err error, format ...string) {
	if err == nil {
		if len(format) > 0 {
			err = fmt.Errorf(format[0], err)
		}
		log.Fatal(err)
	}
}

func mustQueryHello(u string) {
	res := getURL(u)
	mustExec(res, "HTTP request failed %v")
	mustEqStdout(res, "hello, world\n", "Invalid HTTP response got=%q, want=%q")
	log.Printf("SUCCESS: Hello service query %q succeed\n", u)
}

func mustExecHello(p string, u string) {
	helloPackageExec := makeCmd(p)
	mustStart(helloPackageExec, "hello command failed %v")
	mustQueryHello(u)
	mustKill(helloPackageExec, "hello was not killed properly %v")
	log.Printf("SUCCESS: Hello program exec %q and query %q succeed\n", p, u)
}

func _qp(s string) string {
	return _q(_p(s))
}

func _q(s string) string {
	return "\"" + s + "\""
}

func _p(s string) string {
	return filepath.Clean(s)
}

func confirm(err error, message string) {
	if err == nil {
		log.Printf("DONE: %v\n", message)
	} else {
		log.Printf("NOT-DONE: (%v) %v", err, message)
	}
}
func mustSucceed(err error, format ...string) {
	if err != nil {
		if len(format) > 0 {
			err = fmt.Errorf(format[0], err)
		}
		log.Fatal(err)
	}
}
func mustSucceedDetailed(err error, e interface{}, format ...string) {
	if x, ok := e.(stdouter); ok {
		fmt.Printf("%T:%v\n", x, x.Stdout())
	}
	if x, ok := e.(stderrer); ok {
		fmt.Printf("%T:%v\n", x, x.Stderr())
	}
	if err != nil {
		if len(format) > 0 {
			err = fmt.Errorf(format[0], err)
		} else {
			err = fmt.Errorf("%v", err)
		}
		log.Fatal(err)
	} else {

	}
}
func mustFail(err error, format ...string) {
	if err == nil {
		msg := "Expected to fail"
		if len(format) > 0 {
			msg = format[0]
		}
		log.Fatal(msg)
	}
}
func mustShowEnv(e string) {
	psShowEnv := makeCmd("PowerShell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", e)
	mustExec(psShowEnv, "powershell command failed %v")
	log.Printf("showEnv ok %v %q", e, psShowEnv.Stdout())
}
func maybeShowEnv(e string) *cmdExec {
	psShowEnv := makeCmd("PowerShell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", e)
	warnExec(psShowEnv, "powershell command failed %v")
	log.Printf("showEnv ok %v %q", e, psShowEnv.Stdout())
	return psShowEnv
}
func mustLookPath(search string) string {
	path, err := exec.LookPath(search)
	mustSucceed(err, fmt.Sprintf("lookPath failed %q\n%%v", search))
	path = filepath.Clean(path)
	log.Printf("lookPath ok %v=%q", search, path)
	return path
}
func mustChdir(s fmt.Stringer) {
	path := filepath.Clean(s.String())
	mustSucceed(os.Chdir(path), fmt.Sprintf("chDir failed %q\n%%v", path))
	log.Printf("chdir ok %v", path)
}
func mustEnvEq(env string, expect string, format ...string) {
	c := maybeShowEnv(env)
	got := c.Stdout()
	if len(got) > 0 {
		got = got[0 : len(got)-2]
	}
	f := fmt.Sprintf("Env %q is not equal to want=%q, got=%q", env, expect, got)
	mustSucceed(isTrue(got == expect, f))
	log.Printf("mustEnvEq ok %v=%q", env, expect)
}
func mustContains(path fmt.Stringer, file string) {
	s := mustLs(path)
	_, ex := s[file]
	f := fmt.Sprintf("File %q not found in %q", file, path)
	mustSucceed(isTrue(ex, f))
	log.Printf("mustContains ok %v %v", path, file)
}
func mustMkdirAll(path string, perm ...os.FileMode) {
	if len(perm) == 0 {
		perm = []os.FileMode{os.ModePerm}
	}
	path = filepath.Clean(path)
	mustSucceed(os.MkdirAll(path, perm[0]), fmt.Sprintf("mkdirAll failed %q\n%%v", path))
	log.Printf("mkdirAll ok %v %v", path, perm)
}
func mustLs(s fmt.Stringer) map[string]os.FileInfo {
	ret := make(map[string]os.FileInfo)
	path := filepath.Clean(s.String())
	files, err := ioutil.ReadDir(path)
	mustSucceed(err, fmt.Sprintf("readdir failed %q, err=%%v", s))
	for _, f := range files {
		ret[f.Name()] = f
	}
	return ret
}

type starter interface {
	Start() error
}

func mustStart(e starter, format ...string) {
	if len(format) < 1 {
		format[0] = "Start err: %v"
	}
	mustSucceed(e.Start(), format[0])
}

type waiter interface {
	Wait() error
}

func mustWait(e waiter, format ...string) {
	if len(format) < 1 {
		format[0] = "Wait err: %v"
	}
	mustSucceed(e.Wait(), format[0])
}

type execer interface {
	Exec() error
}

func mustExec(e execer, format ...string) {
	if len(format) < 1 {
		format[0] = "Exec err: %v"
	}
	mustSucceedDetailed(e.Exec(), e, format[0])
	log.Printf("mustExec success %v", e)
}

func warnExec(e execer, format ...string) {
	if err := e.Exec(); err != nil {
		if len(format) < 1 {
			format[0] = "Exec err: %v"
		}
		log.Printf(format[0], err)
	}
}

type killer interface {
	Kill() error
}

func mustKill(e killer, format ...string) {
	if len(format) < 1 {
		format[0] = "Kill err: %v"
	}
	mustSucceed(e.Kill(), format[0])
}

type exister interface {
	exists() bool
}

func mustExists(e exister, format ...string) {
	if len(format) < 1 {
		format[0] = fmt.Sprintf("mustExists err: %T does not exist %q, got %%v", e, e)
	}
	mustSucceed(isTrue(e.exists(), format[0]))
}

func isTrue(b bool, format ...string) error {
	if b == false {
		if len(format) < 1 {
			format[0] = "mustTrue got %v"
		}
		return fmt.Errorf(format[0], b)
	}
	return nil
}

type stderrer interface {
	Stderr() string
}

type stdouter interface {
	Stdout() string
}

func mustEqStdout(e stdouter, expected string, format ...string) {
	got := e.Stdout()
	if len(format) < 1 {
		format[0] = fmt.Sprintf("mustEqStdout failed: output does not match, got=%q, want=%q", got, expected)
	}
	mustSucceed(isTrue(got == expected), format[0])
}

func makeFile(f string) *file {
	return &file{f}
}

type file struct {
	path string
}

func (f *file) exists() bool {
	if _, err := os.Stat(f.path); os.IsNotExist(err) {
		return false
	}
	return true
}
func (f *file) String() string {
	return f.path
}

func makeDir(f string) *dir {
	return &dir{f}
}

type dir struct {
	path string
}

func (d *dir) exists() bool {
	if _, err := os.Stat(d.path); os.IsNotExist(err) {
		return false
	}
	return true
}
func (d *dir) String() string {
	return d.path
}

func getURL(url string) *httpRequest {
	return &httpRequest{url: url}
}

type httpRequest struct {
	url        string
	err        error
	body       string
	statusCode int
	headers    map[string][]string
}

func (f *httpRequest) String() string {
	return f.url
}

func (f *httpRequest) Stdout() string {
	return f.body
}
func (f *httpRequest) Header(name string) []string {
	return f.headers[name]
}
func (f *httpRequest) RespondeCode() int {
	return f.statusCode
}
func (f *httpRequest) ExitOk() bool {
	return f.statusCode == 200
}
func (f *httpRequest) Exec() error {
	response, err := http.Get(f.url)
	f.headers = response.Header
	f.statusCode = response.StatusCode
	f.err = err
	if f.err == nil {
		var b bytes.Buffer
		defer response.Body.Close()
		_, f.err = io.Copy(&b, response.Body)
		if f.err == nil {
			f.body = b.String()
		}
	}
	return f.err
}

type cmdExec struct {
	*exec.Cmd
	bin        string
	args       []string
	startErr   error
	hasStarted bool
	waitErr    error
	hasWaited  bool
	stdout     *bytes.Buffer
	stderr     *bytes.Buffer
}

func (e *cmdExec) String() string {
	return e.bin + " " + strings.Join(e.args, " ")
}

func (e *cmdExec) SetArgs(args []string) error {
	if e.hasStarted {
		return fmt.Errorf("Cannot set arguments on command already started")
	}
	if e.hasWaited {
		return fmt.Errorf("Cannot set arguments on command already waited")
	}
	e.args = args
	return nil
}

func (e *cmdExec) Start() error {
	if !e.hasStarted {
		e.startErr = e.Cmd.Start()
	}
	e.hasStarted = true
	return e.startErr
}
func (e *cmdExec) Wait() error {
	if !e.hasWaited {
		e.waitErr = e.Cmd.Wait()
	}
	e.hasWaited = true
	return e.waitErr
}
func (e *cmdExec) Exec() error {
	if err := e.Start(); err != nil {
		return err
	}
	return e.Wait()
}
func (e *cmdExec) Kill() error {
	return e.Process.Kill()
}
func (e *cmdExec) Stdout() string {
	if e.hasStarted == false {
		log.Fatal("Process must have run")
	}
	return e.stdout.String()
}
func (e *cmdExec) Stderr() string {
	if e.hasStarted == false {
		log.Fatal("Process must have run")
	}
	return e.stderr.String()
}
func (e *cmdExec) ExitOk() bool {
	if e.hasWaited == false {
		log.Fatal("Process must have run")
	}
	return e.ProcessState.Exited() && e.ProcessState.Success()
}

func makeCmd(w string, a ...string) *cmdExec {
	log.Printf("makeCmd: %v %v\n", w, a)
	cmd := exec.Command(mustLookPath(w), a...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	// cmd.Stdout = &stdout
	// cmd.Stderr = &stderr
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	return &cmdExec{Cmd: cmd, stdout: &stdout, stderr: &stderr}
}

func readDir(s string) {
	s = filepath.Clean(s)
	files, err := ioutil.ReadDir(s)
	mustSucceed(err, fmt.Sprintf("readdir failed %q, err=%%v", s))
	log.Printf("Content of directory %q\n", s)
	for _, f := range files {
		log.Printf("    %v\n", f.Name())
	}
}

func readFile(s string) {
	s = filepath.Clean(s)
	fd, err := os.Open(s)
	mustSucceed(err, fmt.Sprintf("readfile failed %q, err=%%v", s))
	defer fd.Close()
	if fd != nil {
		io.Copy(fd, os.Stdout)
	}
}

func rmFile(s string) error {
	s = filepath.Clean(s)
	return os.Remove(s)
}
