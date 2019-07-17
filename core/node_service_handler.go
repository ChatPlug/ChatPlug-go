package core

import (
	"fmt"
	"log"
	"os/exec"
)

type NodeServiceHandler struct {
	*BaseServiceHandler
}

func (nsh *NodeServiceHandler) Init() error {
	nsh.ServiceProcesses = make(map[string]*exec.Cmd)

	if !commandExists("node") {
		// TODO: add downloading Node.js
		return fmt.Errorf("Service %s needs Node.js, but it was not found in the system", nsh.Service.Name)
	}

	if !Exists(nsh.GetPath() + "package.json") {
		return nil
	}
	if Exists(nsh.GetPath() + "node_modules") {
		return nil
	}

	serviceName := nsh.Service.Name
	log.Printf("Preparing dependencies for service %s...", serviceName)

	if !commandExists("npm") {
		return fmt.Errorf("Service %s needs installing dependencies, but NPM was not found", serviceName)
	}

	err := exec.Command("npm", "install")

	if err != nil {
		return fmt.Errorf("Installing dependencies for service %s failed: %s", serviceName, err)
	}

	log.Printf("Installed dependencies for service %s", serviceName)
	return nil
}

func (nsh *NodeServiceHandler) LoadService(instanceID string) {
	nsh.RunInstance(exec.Command("node", nsh.GetEntrypointPath()), instanceID)
}

func (nsh *NodeServiceHandler) ShutdownService(instanceID string) {
	// Kill service's process
	nsh.ServiceProcesses[instanceID].Process.Kill()
}
