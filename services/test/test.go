package main

import (
	"context"
	"log"
	"os"

	"github.com/machinebox/graphql"
)

func main() {
	id := os.Args[1]
	client := graphql.NewClient("http://localhost:2137/query")

	// make a request
	req := graphql.NewRequest(`
	mutation ($id: ID!) {
		setInstanceStatus(instanceId:$id, status:INITIALIZED) {
		  status
		  name
		}
	  }`)
	ctx := context.Background()
	var respData map[string]interface{}
	req.Var("id", id)

	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
	}
}
