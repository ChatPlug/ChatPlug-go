package core

import "os/exec"

type ServiceHandler interface {
	VerifyDepedencies() bool
	LoadService(instanceID string)
	ShutdownService(instanceID string)
}

type BaseServiceHandler struct {
	App              *App
	Service          *Service
	ServiceProcesses map[string]*exec.Cmd
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
