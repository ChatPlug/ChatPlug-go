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

// Message stores all information of single message handled by chatplug
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

// Attachment stores data about a attachments contained in single message
type Attachment struct {
	Base
	OriginID  string         `json:"originId"`
	SourceURL string         `json:"sourceUrl"`
	Type      AttachmentType `json:"type"`
	MessageID string
}

// ThreadGroup holds data about a group containing many Threads
type ThreadGroup struct {
	Base
	Name     string    `json:"name"`
	Messages []Message `json:"messages"`
	Threads  []Thread  `json:"threads"`
}

// Thread holds a data about single thread from a single ServiceInstance
type Thread struct {
	Base
	Readonly          *bool     `json:"readonly"`
	Name              string    `json:"name"`
	OriginID          string    `json:"originId"`
	Messages          []Message `json:"messages"`
	ThreadGroupID     string    `json:"threadGroupId"`
	ServiceInstanceID string    `json:"serviceInstanceId"`
}

// ConfigurationField holds data about one field in config
type ConfigurationField struct {
	Type         ConfigurationFieldType `json:"type"`
	DefaultValue string                 `json:"defaultValue"`
	Optional     bool                   `json:"optional"`
	Hint         string                 `json:"hint"`
	Mask         bool                   `json:"mask"`
}

// ConfigurationRequest holds all the requested config fields and keeps a
// result chan that returns a configured response
type ConfigurationRequest struct {
	Fields  []ConfigurationField `json:"fields"`
	resChan chan *ConfigurationResponse
}

// ConfigurationResponse holds a completed config info
type ConfigurationResponse struct {
	FieldValues []string `json:"fieldValues"`
}

// Service holds info about a single service. Can be initialized with many instances via ServiceInstance
type Service struct {
	Base
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Type        string `json:"type"`
	EntryPoint  string `json:"entrypoint"`
}

// ServiceInstance holds info about a single running instance of a Service
type ServiceInstance struct {
	Base
	Name        string         `json:"name"`
	AccessToken string         `json:"-"`
	ModuleName  string         `json:"moduleName"`
	Threads     []Thread       `json:"threads"`
	Status      InstanceStatus `json:"status"`
}

// MessageAuthor holds data about a single user in single ServiceInstance
type MessageAuthor struct {
	Base
	Username  string    `json:"username"`
	OriginID  string    `json:"originId"`
	AvatarURL string    `json:"avatarUrl" gorm:"default:'https://i.imgur.com/3yPh9fE.png'"`
	Messages  []Message `json:"messages"`
}
