package core

import (
	"log"
	"os"
	"os/exec"
	"path"
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
	var serviceInstance ServiceInstance

	bsh.App.db.First(&serviceInstance, "id = ?", instanceID)
	command.Args = append(command.Args, instanceID)
	command.Dir = path.Join("services", bsh.Service.Name)
	command.Env = os.Environ()
	command.Env = append(command.Env,
		"WS_ENDPOINT=ws://localhost:2137/query",
		"HTTP_ENDPOINT=http://localhost:2137/query",
		"INSTANCE_ID="+instanceID,
		"ACCESS_TOKEN="+serviceInstance.AccessToken)

	bsh.ServiceProcesses[instanceID] = command

	log.Printf("Launching service %s...", bsh.Service.Name)
	go func() {
		err := RunCommand(command, bsh.Service.Name)
		if err != nil {
			log.Println(err)
		}
	}()
}
