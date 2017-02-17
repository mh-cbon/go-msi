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
)

func main() {

	wd := "c:\\gopath\\src\\github.com\\mh-cbon\\go-msi\\testing\\hello"
	os.Chdir(wd)

	execCmd("go", "build", "-o", "build\\amd64\\hello.exe", "hello.go")
	os.MkdirAll("build\\amd64", os.ModePerm)
	execCmd("C:\\go-msi\\go-msi.exe", "make", "--msi", "hello.msi", "--version", "0.0.1", "--arch", "amd64")

	fileExists("hello.msi")
	execCmd("c:\\windows\\system32\\msiexec.exe", "/i", "hello.msi", "/q")

	execCmd("PowerShell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", "$env:path")
	tryCmd("PowerShell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", "ls env:some")

	readDir("C:\\Program Files")
	readDir("C:\\Program Files\\hello")
	readDir("C:\\Program Files\\hello\\assets")
	hello := execBgCmd("C:\\Program Files\\hello\\hello.exe")
	b := fecthURL("http://localhost:8080/")
	hello.Process.Kill()
	expected := "hello, world\n"
	if b != expected {
		panic(fmt.Errorf("Invalid http response %q, expected %q", b, expected))
	}
	execCmd("c:\\windows\\system32\\msiexec.exe", "/x", "hello.msi", "/quiet")

	execCmd("C:\\go-msi\\go-msi.exe", "choco", "--input", "hello.msi", "--version", "0.0.1", "-c", "\"C:\\Program Files\\changelog\\changelog.exe\" ghrelease --version 0.0.1")
	execCmd("choco", "install", "hello.0.0.1.nupkg", "-y")
	execCmd("PowerShell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", "$env:path")
	hello2 := execBgCmd("C:\\Program Files\\hello\\hello.exe")
	b2 := fecthURL("http://localhost:8080/")
	hello2.Process.Kill()
	if b2 != expected {
		panic(fmt.Errorf("Invalid http response %q, expected %q", b2, expected))
	}

	execCmd("choco", "uninstall", "hello", "-y")
}

func getCmd(w string, a ...string) *exec.Cmd {
	fmt.Printf("getCmd: %v %v\n", w, a)
	cmd := exec.Command(w, a...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func tryCmd(w string, a ...string) {
	cmd := getCmd(w, a...)
	if err := cmd.Start(); err != nil {
		fmt.Println(err)
	}
	if err := cmd.Wait(); err != nil {
		fmt.Println(err)
	}
}

func execCmd(w string, a ...string) {
	cmd := getCmd(w, a...)
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}

func execBgCmd(w string, a ...string) *exec.Cmd {
	cmd := getCmd(w, a...)
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	return cmd
}

func fecthURL(url string) string {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	var b bytes.Buffer
	defer response.Body.Close()
	_, err2 := io.Copy(&b, response.Body)
	if err2 != nil {
		log.Fatal(err2)
	}
	return b.String()
}

func fileExists(f string) {
	if _, err := os.Stat(f); os.IsNotExist(err) {
		log.Fatalf("File not found : %q\n", f)
	}
}
func readDir(s string) {
	files, err := ioutil.ReadDir(s)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Content of directory %q\n", s)
	for _, f := range files {
		fmt.Printf("    %v\n", f.Name())
	}
}
