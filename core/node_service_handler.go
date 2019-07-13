package core

import "os/exec"

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
	nsh.App.db.First(&instance, instanceID)

	service := nsh.App.sm.FindServiceWithName(instance.ModuleName)

	// Startup service
	nsh.ServiceProcesses[instanceID] = exec.Command("node", "services/"+service.Name+"/"+service.EntryPoint, instance.ID)
}

func (nsh *NodeServiceHandler) ShutdownService(instanceID string) {
	// Kill service's process
	nsh.ServiceProcesses[instanceID].Process.Kill()
}
