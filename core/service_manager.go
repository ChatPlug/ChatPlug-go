package core

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type ServiceManager struct {
	services []*Service
}

func (sm *ServiceManager) LoadAvailableServices() {
	files, err := ioutil.ReadDir("services")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		manifestPath := "services/" + f.Name() + "/manifest.json"
		if f.IsDir() && Exists(manifestPath) {
			jsonFile, err := os.Open(manifestPath)
			if err == nil {
				byteVal, _ := ioutil.ReadAll(jsonFile)
				var service Service
				json.Unmarshal(byteVal, &service)

				log.Printf("Service found! %s", service.DisplayName)
				sm.services = append(sm.services, &service)
			}
		}
	}
}

func (sm *ServiceManager) FindServiceWithName(moduleName string) *Service {
	for _, n := range sm.services {
		if n.Name == moduleName {
			return n
		}
	}
	return nil
}

func (sm *ServiceManager) StartupServiceInstance(instance *ServiceInstance) {
	service := sm.FindServiceWithName(instance.ModuleName)

	cmd := exec.Command("services/"+service.Name+"/"+service.EntryPoint, "--id="+instance.ID)
	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)
	cmd.Stderr = mw
	cmd.Stdout = mw
	err := cmd.Run() //blocks until sub process is complete
	if err != nil {
		panic(err)
	}
	log.Println(stdBuffer.String())
}

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
