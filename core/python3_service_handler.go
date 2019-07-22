package core

import (
	"fmt"
	"log"
	"os/exec"
)

type Python3ServiceHandler struct {
	*BaseServiceHandler
}

func (psh *Python3ServiceHandler) Init() error {
	psh.ServiceProcesses = make(map[string]*exec.Cmd)

	if !commandExists("python3") {
		// TODO: add downloading Python
		return fmt.Errorf("Service %s needs Python 3, but it was not found in the system", psh.Service.Name)
	}

	if !Exists(psh.GetPath() + "requirements.txt") {
		return nil
	}

	serviceName := psh.Service.Name
	log.Printf("Preparing dependencies for service %s...", serviceName)

	if !commandExists("pip3") {
		return fmt.Errorf("Service %s needs installing dependencies, but Pip3 was not found", serviceName)
	}

	cmd := exec.Command("pip3", "install", "-r", "requirements.txt")
	cmd.Dir = psh.GetPath()
	if err := RunCommand(cmd, psh.Service.Name+": pip3"); err != nil {
		return fmt.Errorf("Installing dependencies for service %s failed: %s", serviceName, err)
	}

	log.Printf("Installed dependencies for service %s", serviceName)
	return nil
}

func (psh *Python3ServiceHandler) LoadService(instanceID string) {
	psh.RunInstance(exec.Command("python3", psh.GetEntrypointPath()), instanceID)
}

func (psh *Python3ServiceHandler) ShutdownService(instanceID string) {
	// Kill service's process
	psh.ServiceProcesses[instanceID].Process.Kill()
}
