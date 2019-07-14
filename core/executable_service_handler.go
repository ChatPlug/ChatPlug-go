package core

import (
	"log"
	"os/exec"
)

type ExecutableServiceHandler struct {
	App              *App
	ServiceProcesses map[string]*exec.Cmd
}

func (esh *ExecutableServiceHandler) VerifyDepedencies() bool {
	esh.ServiceProcesses = make(map[string]*exec.Cmd)
	return true
}

func (esh *ExecutableServiceHandler) LoadService(instanceID string) {
	var instance ServiceInstance
	esh.App.db.Where("id = ?", instanceID).First(&instance)

	service := esh.App.sm.FindServiceWithName(instance.ModuleName)

	log.Println("asd")
	// Startup service
	esh.ServiceProcesses[instanceID] = exec.Command("services/"+service.Name+"/"+service.EntryPoint, instance.ID)
	esh.ServiceProcesses[instanceID].Run()
}

func (esh *ExecutableServiceHandler) ShutdownService(instanceID string) {
	// Kill service's process
	esh.ServiceProcesses[instanceID].Process.Kill()
}
