package core

import (
	"bufio"
	"io"
	"log"
	"os/exec"
)

// https://stackoverflow.com/a/38870609
func RunCommand(command *exec.Cmd, prefix string) error {
	stdout, err := command.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := command.StderrPipe()
	if err != nil {
		return err
	}

	cout := make(chan struct{})
	cerr := make(chan struct{})

	go readPipe(stdout, prefix, cout)
	go readPipe(stderr, prefix, cerr)

	cout <- struct{}{}
	cerr <- struct{}{}
	command.Start()

	<-cout
	<-cerr
	if err := command.Wait(); err != nil {
		return err
	}

	command.Start()
	return nil
}

func readPipe(pipe io.ReadCloser, prefix string, c chan struct{}) {
	defer func() { c <- struct{}{} }()
	<-c
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		m := scanner.Text()
		log.Println("[" + prefix + "] " + m)
	}
}
