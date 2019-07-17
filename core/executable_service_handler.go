package core

import (
	"os/exec"
)

type ExecutableServiceHandler struct {
	*BaseServiceHandler
}

func (esh *ExecutableServiceHandler) VerifyDepedencies() bool {
	esh.ServiceProcesses = make(map[string]*exec.Cmd)
	return true
}

func (esh *ExecutableServiceHandler) LoadService(instanceID string) {
	esh.RunInstance(exec.Command(esh.GetEntrypointPath()), instanceID)
}

func (esh *ExecutableServiceHandler) ShutdownService(instanceID string) {
	// Kill service's process
	esh.ServiceProcesses[instanceID].Process.Kill()
}
