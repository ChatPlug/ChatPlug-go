
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"context"

	"github.com/gorilla/websocket"
	"github.com/machinebox/graphql"
)

type session struct {
	ws      *websocket.Conn
	errChan chan error
}

// ChatPlugClient holds connection with chatplug core server
type ChatPlugClient struct {
	session *session
	client *graphql.Client
	instanceID string
	wsEndpoint string
}

const (
	connectionInitMsg      = "connection_init"      // Client -> Server
	connectionTerminateMsg = "connection_terminate" // Client -> Server
	startMsg               = "start"                // Client -> Server
	stopMsg                = "stop"                 // Client -> Server
	connectionAckMsg       = "connection_ack"       // Server -> Client
	connectionErrorMsg     = "connection_error"     // Server -> Client
	dataMsg                = "data"                 // Server -> Client
	errorMsg               = "error"                // Server -> Client
	completeMsg            = "complete"             // Server -> Client
	connectionKeepAliveMsg = "ka"                 // Server -> Client
	sendMessageMutation = `
	mutation sendMessage($instanceId: ID!, $body: String!, $originId: String!, $originThreadId: String!, $username: String!, $authorOriginId: String!, $authorAvatarUrl: String!, $attachments: [AttachmentInput!]!) {
		sendMessage(
		  instanceId: $instanceId,
		  input: {
			body: $body,
			originId: $originId,
			originThreadId: $originThreadId,
			author: {
			  username: $username, 
			  originId: $authorOriginId,
			  avatarUrl: $authorAvatarUrl
			},
			attachments: $attachments
		  }
		) {
		  id
		}
	  }`
	messageReceivedSubscription = `
	  subscription ($id: ID!) {
		  messageReceived(instanceId:$id) {
			message {
			  body
			  id
			  originId
			  attachments {
				  type
				  sourceUrl
				  originId
				  id
			  }
			  thread {
				  id
				  originId
				  name
			  }
			  threadGroupId
			  author {
				  username
				  originId
				  avatarUrl
			  }
			}
			targetThreadId
		  }
		}`
	requestConfigurationRequest = `
	subscription confRequest($fields: [ConfigurationField!]!){
		configurationReceived(configuration:{fields: $fields}) {
		  fieldValues
		}
	  }`

	setInstanceStatusMutation = `
	mutation ($id: ID!) {
		setInstanceStatus(instanceId:$id, status:INITIALIZED) {
		  status
		  name
		}
	  }`
)

// MessageAuthor holds information about single message's atuhor
type MessageAuthor struct {
	ID string `json: "id"`
	Username string    `json:"username"`
	OriginID string    `json:"originId"`
	AvatarURL string `json:"avatarUrl"`
}

type AttachmentInput struct {
	OriginID string `json:"originId"`
	Type string    `json:"type"`
	SourceURL string    `json:"sourceUrl"`
}

type Attachment struct {
	ID string `json:"id"`
	OriginID string `json:"originId"`
	Type string    `json:"type"`
	SourceURL string    `json:"sourceUrl"`
}

// Thread holds information about single thread
type Thread struct {
	ID string `json:"id"`
	Name              string    `json:"name"`
	OriginID          string    `json:"originId"`
	ThreadGroupID     string    `json:"threadGroupId"`
	ServiceInstanceID string    `json:"serviceInstanceId"`
}

// Message holds information about single message
type Message struct {
	ID string `json:"string"`
	OriginID        string `json:"originId"`
	Author MessageAuthor `json:"author"`
	Thread Thread `json:"thread"`
	Body            string `json:"body"`
	ThreadGroupID   string `json:"threadGroupId"`
	Attachments []Attachment `json:"attachments"`
} 

// ErrorLocation holds data about location of gql error
type ErrorLocation struct {
	Line int`json:"line,omitempty"`
	Column int`json:"column,omitempty"`
}
// ErrorMessage holds data about graphql error
type ErrorMessage struct {
	Message string`json:"message,omitempty"`
	Locations []*ErrorLocation `json:"locations,omitempty"`
}

// MessageReceived holds data about incoming message
type MessageReceived struct {
	Message Message `json:"message"`
	TargetThreadID string `json:"targetThreadId"`
}

type messageReceivedPayload struct {
	Data struct {
		MessageReceived MessageReceived `json:"messageReceived"`
	}`json:"data"`
}

type configurationReceivedPayload struct {
	Data struct {
		ConfigurationReceived ConfigurationResponse `json:"configurationReceived"`
	}`json:"data"`
}

type operationMessage struct {
	Payload payloadMessage `json:"payload,omitempty"`
	ID      string          `json:"id,omitempty"`
	Type    string          `json:"type"`
}

// IncomingPayload is a struct holding graphql operation data
type IncomingPayload struct {
	Payload *json.RawMessage `json:"payload,omitempty"`
	Type    string          `json:"type"`
}


type payloadMessage struct {
	Query string `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

func wsConnect(url string) *websocket.Conn {
	headers := make(http.Header)
	headers.Add("Sec-Websocket-Protocol", "graphql-ws")
	c, _, err := websocket.DefaultDialer.Dial(url, headers)

	if err != nil {
		panic(err)
	}
	return c
}

// ReadOp reads a single graphql operation from underlying websocket connection
func (s *session) ReadOp() (*IncomingPayload, error) {
	var msg IncomingPayload
	err := s.ws.ReadJSON(&msg)
	if err != nil {
		panic(err)
	}
	return &msg, err
}

// Subscribe starts a given graphql subscription and returns a chan with incoming data
func (s *session) Subscribe(query string, variables map[string]interface{}) (<-chan *IncomingPayload, <-chan error) {

	channel := make(chan *IncomingPayload)

	s.ws.WriteJSON(&operationMessage{
		Type: startMsg,
		ID: "1",
		Payload: payloadMessage{
			Query: query,
			Variables: variables,
		},
	})

	go func() {
		for {

			msg, err := s.ReadOp()
			if err != nil {
				s.errChan <- err
			}
			// kok, _ := json.MarshalIndent(&msg, "", "    ")
			// fmt.Printf("%s\n", kok)

			if (msg.Type == "error") {
				var errs []*ErrorMessage
				err = json.Unmarshal(*msg.Payload, &errs)
				log.Println(errs[0].Message)
				log.Println(string(errs[0].Locations[0].Line))
			}
			channel <- msg

		}
		close(channel)
		close(s.errChan)
	}()

	return channel, s.errChan
}

// Connect starts a websocket connection to the server and notifies it about it's initialization
func (gqc *ChatPlugClient) Connect(instanceID string, httpEndpoint string, wsEndpoint string) {
	gqc.client = graphql.NewClient(httpEndpoint)
	gqc.instanceID = instanceID
	gqc.wsEndpoint = wsEndpoint
	c := wsConnect(wsEndpoint)

	gqc.session = &session{
		ws: c,
	}
	gqc.session.ws.WriteJSON(&operationMessage{Type: connectionInitMsg})
	gqc.session.ReadOp()

	req := graphql.NewRequest(setInstanceStatusMutation)
	req.Var("id", instanceID)

	gqc.Request(req)
}

// Close closes a websocket connection to the server
func (gqc *ChatPlugClient) Close() {
	gqc.session.ws.Close()
}

// SendMessage sends a message with given data to core server via graphql
func (gqc *ChatPlugClient) SendMessage(body string, originId string, originThreadId string, username string, authorOriginId string, authorAvatarUrl string, attachments []*AttachmentInput) {
	req := graphql.NewRequest(sendMessageMutation)
	req.Var("instanceId", gqc.instanceID)
	req.Var("body", body)
	req.Var("originId", originId)
	req.Var("originThreadId", originThreadId)
	req.Var("username", username)
	req.Var("authorOriginId", authorOriginId)
	req.Var("authorAvatarUrl", authorAvatarUrl)
	req.Var("attachments", attachments)

	fmt.Println("Sending sendMessage mutation to the core")
	_, err := gqc.Request(req)
	if err != nil {
		fmt.Println("boop")
		fmt.Println(err)
	}
}

// SubscribeToNewMessages starts a subscription to core server's messages and returns a chan with parsed data
func (gqc *ChatPlugClient) SubscribeToNewMessages() <-chan *MessageReceived {
	// gqc.session.ws.Close()
	c := wsConnect(gqc.wsEndpoint)

	gqc.session = &session{
		ws: c,
	}

	gqc.session.ws.WriteJSON(&operationMessage{Type: connectionInitMsg})
	gqc.session.ReadOp()

	variables := make(map[string]interface{})
	variables["id"] = gqc.instanceID
	channel := make(chan *MessageReceived)

	subscriptionChan, _ := gqc.session.Subscribe(messageReceivedSubscription, variables)
	go func() {
		for subscription := range subscriptionChan {
			if (subscription.Type == "data") {
				var msg messageReceivedPayload
				json.Unmarshal(*subscription.Payload, &msg)
				channel <- &msg.Data.MessageReceived
			}
		}
	}()
	return channel
}

type ConfigurationField struct {
	Type         string `json:"type"`
	DefaultValue string                 `json:"defaultValue"`
	Optional     bool                   `json:"optional"`
	Hint         string                 `json:"hint"`
	Mask         bool                   `json:"mask"`
}

type ConfigurationRequest struct {
	Fields  []ConfigurationField `json:"fields"`
}

type ConfigurationResponse struct {
	FieldValues []string `json:"fieldValues"`
}

func (gqc *ChatPlugClient) AwaitConfiguration(configurationSchema []ConfigurationField) *ConfigurationResponse {
	variables := make(map[string]interface{})
	variables["fields"] = configurationSchema
	channel := make(chan *ConfigurationResponse)

	subscriptionChan, _ := gqc.session.Subscribe(requestConfigurationRequest, variables)
	go func() {
		Loop:
		for {
		for subscription := range subscriptionChan {
			if (subscription.Type == "data") {
				var cfg configurationReceivedPayload
				json.Unmarshal(*subscription.Payload, &cfg)
				channel <- &cfg.Data.ConfigurationReceived
				break Loop
			}
			}
		}
	}()
	res := <- channel
	return res
}

// Request sends a graphql requests to the core server and returns a pointer to map with result
func (gqc *ChatPlugClient) Request(req *graphql.Request) (*map[string]interface{}, error) {
	// make a request
	ctx := context.Background()
	var respData map[string]interface{}

	if err := gqc.client.Run(ctx, req, &respData); err != nil {
		return nil, err
	}
	return &respData, nil
}