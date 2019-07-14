package core

import (
	"os"
	"os/exec"
)

type NodeServiceHandler struct {
	App              *App
	ServiceProcesses map[string]*exec.Cmd
}

func (nsh *NodeServiceHandler) VerifyDepedencies() bool {
	nsh.ServiceProcesses = make(map[string]*exec.Cmd)
	return true
}

func (nsh *NodeServiceHandler) LoadService(instanceID string) {
	var instance ServiceInstance
	nsh.App.db.Where("id = ?", instanceID).First(&instance)

	service := nsh.App.sm.FindServiceWithName(instance.ModuleName)

	// Startup service
	command := exec.Command("node", "services/"+service.Name+"/"+service.EntryPoint, instance.ID)
	nsh.ServiceProcesses[instanceID] = command
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Run()
}

func (nsh *NodeServiceHandler) ShutdownService(instanceID string) {
	// Kill service's process
	nsh.ServiceProcesses[instanceID].Process.Kill()
}
