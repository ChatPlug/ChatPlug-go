package core

import (
	"os"
	"os/exec"
)

type ServiceHandler interface {
	Init() error
	LoadService(instanceID string)
	ShutdownService(instanceID string)
}

type BaseServiceHandler struct {
	App              *App
	Service          *Service
	ServiceProcesses map[string]*exec.Cmd
}

func (bsh *BaseServiceHandler) GetPath() string {
	return "services/" + bsh.Service.Name + "/"
}

func (bsh *BaseServiceHandler) GetEntrypointPath() string {
	return bsh.GetPath() + bsh.Service.EntryPoint
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func (bsh *BaseServiceHandler) RunInstance(command *exec.Cmd, instanceID string) {
	command.Args = append(command.Args, instanceID)
	bsh.ServiceProcesses[instanceID] = command

	command.Stdout = os.Stdout
	command.Stdin = os.Stdin

	command.Run()
}
