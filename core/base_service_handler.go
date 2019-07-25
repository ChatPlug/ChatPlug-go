package core

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	return bsh.Service.EntryPoint
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func (bsh *BaseServiceHandler) RunInstance(command *exec.Cmd, instanceID string) {
	command.Args = append(command.Args, instanceID)
	command.Dir = filepath.FromSlash("services/" + bsh.Service.Name)
	command.Env = os.Environ()
	command.Env = append(command.Env, "WS_ENDPOINT=ws://localhost:2137/query")
	command.Env = append(command.Env, "HTTP_ENDPOINT=http://localhost:2137/query")
	command.Env = append(command.Env, "INSTANCE_ID="+instanceID)

	bsh.ServiceProcesses[instanceID] = command
	fmt.Println(command.Path)
	fmt.Println(command.Dir)

	go func() {
		RunCommand(command, bsh.Service.Name)
	}()
}
