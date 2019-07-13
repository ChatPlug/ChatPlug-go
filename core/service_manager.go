package core

import (
	"encoding/json"
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

	exec.Command("services/"+service.Name+"/"+service.EntryPoint, instance.ID)
}

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
