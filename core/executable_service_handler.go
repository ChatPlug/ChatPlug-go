package core

import (
	"os/exec"
)

type ExecutableServiceHandler struct {
	*BaseServiceHandler
}

func (esh *ExecutableServiceHandler) Init() error {
	esh.ServiceProcesses = make(map[string]*exec.Cmd)
	return nil
}

func (esh *ExecutableServiceHandler) LoadService(instanceID string) {
	esh.RunInstance(exec.Command(esh.GetEntrypointPath()), instanceID)
}

func (esh *ExecutableServiceHandler) ShutdownService(instanceID string) {
	// Kill service's process
	esh.ServiceProcesses[instanceID].Process.Kill()
}
