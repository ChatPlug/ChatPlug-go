
package main

import (
	"fmt"
	"os"
)

func main() {
	id := os.Args[1]
	client := ChatPlugClient{}
	client.Connect(id, "http://localhost:2137/query", "ws://localhost:2137/query")
	conf := make([]ConfigurationField, 0)
	ques1 := ConfigurationField{
		Type: "STRING",
		Hint: "dupa",
		DefaultValue: "ay",
		Optional: false,
		Mask: false,
	}
	conf = append(conf, ques1)

	config := client.AwaitConfiguration(conf)
	for _, a := range config.FieldValues {
		fmt.Println(a)
	}
	msgChan := client.SubscribeToNewMessages()

	for msg := range msgChan {
		fmt.Printf(msg.Message.Attachments[0].SourceURL)
	}
	defer client.Close()
}
