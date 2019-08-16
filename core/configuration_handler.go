package core

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/manifoldco/promptui"
)

type ConfigurationHandler struct {
	configurationQueue chan *ConfigurationRequest
}

func (ch *ConfigurationHandler) WatchForConfiguration() {
	ch.configurationQueue = make(chan *ConfigurationRequest, 10)
	go func() {
		for configRequest := range ch.configurationQueue {
			clearConsole()
			log.Println("Configuration requested!")
			response, err := ch.PromptForConfiguration(configRequest)
			if err == nil {
				configRequest.resChan <- response
			}
		}
	}()
}

func (ch *ConfigurationHandler) PromptForConfiguration(request *ConfigurationRequest) (*ConfigurationResponse, error) {
	response := &ConfigurationResponse{
		FieldValues: make([]ConfigurationResult, 0),
	}

	for _, field := range request.Fields {
		prompt := promptui.Prompt{
			Label: field.Hint,
		}

		if field.Mask {
			prompt.Mask = '*'
		}

		if field.Optional {
			prompt.Default = field.DefaultValue
		}

		switch fieldType := field.Type; fieldType {
		case "BOOLEAN":
			prompt.Validate = func(input string) error {
				if input != "y" && input != "n" && input != "Y" && input != "N" {
					return errors.New("Boolean should be \"y\" or \"n\"")
				}
				return nil
			}
			break
		case "STRING":
			break
		case "NUMBER":
			prompt.Validate = func(input string) error {
				_, err := strconv.ParseFloat(input, 64)
				if err != nil {
					return errors.New("Invalid number")
				}
				return nil
			}
			break
		}

		result, err := prompt.Run()
		if err != nil {
			return nil, errors.New("Configuration failed")
		}
		response.FieldValues = append(response.FieldValues, ConfigurationResult{Name: field.Name, Value: result})
	}

	return response, nil
}

func clearConsole() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}
