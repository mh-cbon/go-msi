package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/kardianos/service"
)

var logger service.Logger

func main() {
	svcConfig := &service.Config{
		Name:        "HelloSvc",
		DisplayName: "The Hello service",
		Description: "This is an example Go service.",
	}
	runAsService(svcConfig, func() {
		fmt.Println("Starting....")
		http.HandleFunc("/", handleIndex)
		log.Fatal(http.ListenAndServe(":8080", nil))
	})
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello, world\n")
}

func runAsService(svcConfig *service.Config, run func()) error {
	s, err := service.New(&program{exec: run}, svcConfig)
	if err != nil {
		return err
	}
	logger, err = s.Logger(nil)
	if err != nil {
		return err
	}
	return s.Run()
}

type program struct {
	exec func()
}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.exec()
	return nil
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}
