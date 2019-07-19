
package main

import (
	"fmt"
	"os"
)

func main() {
	id := os.Args[1]
	client := ChatPlugClient{}
	client.Connect(id, "http://localhost:2137/query", "ws://localhost:2137/query")
	msgChan := client.SubscribeToNewMessages()

	for msg := range msgChan {
		fmt.Printf(msg.Message.Attachments[0].SourceURL)
	}
	defer client.Close()
}
