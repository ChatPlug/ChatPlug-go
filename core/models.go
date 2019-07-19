package core

import (
	"crypto/rand"
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

// Base contains common columns for all tables.
type Base struct {
	ID        string     `gorm:"type:uuid;primary_key;" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"update_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *Base) BeforeCreate(scope *gorm.Scope) error {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	// guid := ksuid.New()
	return scope.SetColumn("ID", uuid)
}

type Message struct {
	Base
	OriginID        string       `json:"originId"`
	MessageAuthorID string       `json:"messageAuthorId"`
	ThreadID        string       `json:"threadId"`
	Body            string       `json:"body"`
	ThreadGroupID   string       `json:"threadGroupId"`
	Author          string       `gorm:"-" json:"author"`
	Attachments     []Attachment `json:"attachments"`
}

type Attachment struct {
	Base
	OriginID  string         `json:"originId"`
	SourceURL string         `json:"sourceUrl"`
	Type      AttachmentType `json:"type"`
	MessageID string
}

type ThreadGroup struct {
	Base
	Name     string    `json:"name"`
	Messages []Message `json:"messages"`
	Threads  []Thread  `json:"threads"`
}

type Thread struct {
	Base
	Name              string    `json:"name"`
	OriginID          string    `json:"originId"`
	Messages          []Message `json:"messages"`
	ThreadGroupID     string    `json:"threadGroupId"`
	ServiceInstanceID string    `json:"serviceInstanceId"`
}

type Service struct {
	Base
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Type        string `json:"type"`
	EntryPoint  string `json:"entrypoint"`
}

type ServiceInstance struct {
	Base
	Name       string         `json:"name"`
	ModuleName string         `json:"moduleName"`
	Threads    []Thread       `json:"threads"`
	Status     InstanceStatus `json:"status"`
}

type MessageAuthor struct {
	Base
	Username string    `json:"username"`
	OriginID string    `json:"originId"`
	Messages []Message `json:"messages"`
}
