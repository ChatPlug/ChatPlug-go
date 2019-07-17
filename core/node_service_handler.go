package core

import (
	"os/exec"
)

type NodeServiceHandler struct {
	*BaseServiceHandler
}

func (nsh *NodeServiceHandler) VerifyDepedencies() bool {
	nsh.ServiceProcesses = make(map[string]*exec.Cmd)

	return commandExists("node")
	// TODO: add downloading Node.js
}

func (nsh *NodeServiceHandler) LoadService(instanceID string) {
	nsh.RunInstance(exec.Command("node", nsh.GetEntrypointPath()), instanceID)
}

func (nsh *NodeServiceHandler) ShutdownService(instanceID string) {
	// Kill service's process
	nsh.ServiceProcesses[instanceID].Process.Kill()
}
