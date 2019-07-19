package core

import (
	"errors"
	"strconv"

	"github.com/manifoldco/promptui"
)

type ConfigurationFieldType string

const (
	Boolean ConfigurationFieldType = "boolean"
	Number  ConfigurationFieldType = "number"
	String  ConfigurationFieldType = "string"
)

var configurationFieldTypes = []string{"boolean", "number", "string"}

type ConfigurationField struct {
	Type         ConfigurationFieldType `json:"type"`
	DefaultValue string                 `json:"defaultValue"`
	Optional     bool                   `json:"optional"`
	Hint         string                 `json:"hint"`
	Mask         bool                   `json:"mask"`
}

type ConfigurationRequest struct {
	fields []ConfigurationField `json:"fields"`
}

type ConfigurationResponse struct {
	fields []string `json:"fieldValues"`
}

type ConfigurationHandler struct {
}

func (ch *ConfigurationHandler) PromptForConfiguration(request *ConfigurationRequest) (*ConfigurationResponse, error) {
	response := &ConfigurationResponse{
		fields: make([]string, 0),
	}

	for _, field := range request.fields {
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
		case Boolean:
			prompt.Validate = func(input string) error {
				if input != "y" && input != "n" && input != "Y" && input != "N" {
					return errors.New("Boolean should be \"y\" or \"n\"")
				}
				return nil
			}
			break
		case String:
			break
		case Number:
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
		response.fields = append(response.fields, result)
	}

	return response, nil
}
