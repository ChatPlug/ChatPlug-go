// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package core

import (
	"fmt"
	"io"
	"strconv"
)

type AttachmentInput struct {
	OriginID  string         `json:"originId"`
	Type      AttachmentType `json:"type"`
	SourceURL string         `json:"sourceUrl"`
}

type MessageAuthorInput struct {
	OriginID  string `json:"originId"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatarUrl"`
}

type MessageInput struct {
	Body           string              `json:"body"`
	Author         *MessageAuthorInput `json:"author"`
	Attachments    []*AttachmentInput  `json:"attachments"`
	OriginID       string              `json:"originId"`
	OriginThreadID string              `json:"originThreadId"`
	AvatarURL      *string             `json:"avatarUrl"`
}

type MessagePayload struct {
	TargetThreadID string   `json:"targetThreadId"`
	Message        *Message `json:"message"`
}

type NewServiceInstanceCreated struct {
	Instance    *ServiceInstance `json:"instance"`
	AccessToken string           `json:"accessToken"`
}

type SearchRequest struct {
	Query string `json:"query"`
}

type ThreadInput struct {
	InstanceID string `json:"instanceId"`
	OriginID   string `json:"originId"`
	GroupID    string `json:"groupId"`
	Readonly   *bool  `json:"readonly"`
	Name       string `json:"name"`
}

type AttachmentType string

const (
	AttachmentTypeFile  AttachmentType = "FILE"
	AttachmentTypeImage AttachmentType = "IMAGE"
	AttachmentTypeAudio AttachmentType = "AUDIO"
	AttachmentTypeVideo AttachmentType = "VIDEO"
)

var AllAttachmentType = []AttachmentType{
	AttachmentTypeFile,
	AttachmentTypeImage,
	AttachmentTypeAudio,
	AttachmentTypeVideo,
}

func (e AttachmentType) IsValid() bool {
	switch e {
	case AttachmentTypeFile, AttachmentTypeImage, AttachmentTypeAudio, AttachmentTypeVideo:
		return true
	}
	return false
}

func (e AttachmentType) String() string {
	return string(e)
}

func (e *AttachmentType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AttachmentType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AttachmentType", str)
	}
	return nil
}

func (e AttachmentType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ConfigurationFieldType string

const (
	ConfigurationFieldTypeBoolean ConfigurationFieldType = "BOOLEAN"
	ConfigurationFieldTypeString  ConfigurationFieldType = "STRING"
	ConfigurationFieldTypeNumber  ConfigurationFieldType = "NUMBER"
)

var AllConfigurationFieldType = []ConfigurationFieldType{
	ConfigurationFieldTypeBoolean,
	ConfigurationFieldTypeString,
	ConfigurationFieldTypeNumber,
}

func (e ConfigurationFieldType) IsValid() bool {
	switch e {
	case ConfigurationFieldTypeBoolean, ConfigurationFieldTypeString, ConfigurationFieldTypeNumber:
		return true
	}
	return false
}

func (e ConfigurationFieldType) String() string {
	return string(e)
}

func (e *ConfigurationFieldType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ConfigurationFieldType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ConfigurationFieldType", str)
	}
	return nil
}

func (e ConfigurationFieldType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type InstanceStatus string

const (
	InstanceStatusRunning      InstanceStatus = "RUNNING"
	InstanceStatusInitialized  InstanceStatus = "INITIALIZED"
	InstanceStatusConfigured   InstanceStatus = "CONFIGURED"
	InstanceStatusShuttingDown InstanceStatus = "SHUTTING_DOWN"
	InstanceStatusStopped      InstanceStatus = "STOPPED"
)

var AllInstanceStatus = []InstanceStatus{
	InstanceStatusRunning,
	InstanceStatusInitialized,
	InstanceStatusConfigured,
	InstanceStatusShuttingDown,
	InstanceStatusStopped,
}

func (e InstanceStatus) IsValid() bool {
	switch e {
	case InstanceStatusRunning, InstanceStatusInitialized, InstanceStatusConfigured, InstanceStatusShuttingDown, InstanceStatusStopped:
		return true
	}
	return false
}

func (e InstanceStatus) String() string {
	return string(e)
}

func (e *InstanceStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = InstanceStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid InstanceStatus", str)
	}
	return nil
}

func (e InstanceStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
