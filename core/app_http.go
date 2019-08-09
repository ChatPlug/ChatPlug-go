package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
var instanceCtxKey = &contextKey{"serviceInstance"}

type contextKey struct {
	name string
}

func (app *App) CreateServiceInstanceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.GetHeader("Authentication")
		fmt.Println(accessToken)

		// Allow unauthenticated users in
		if accessToken == "" || c.Request.Context() == nil {
			c.Next()
			return
		}

		var serviceInstance ServiceInstance

		app.db.First(&serviceInstance, "access_token = ?", accessToken)

		// put it in context
		ctx := context.WithValue(c.Request.Context(), instanceCtxKey, &serviceInstance)

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// ForContext finds the ServiceInstance from the context
func (app *App) InstanceForContext(ctx context.Context) *ServiceInstance {
	raw, ok := ctx.Value(instanceCtxKey).(*ServiceInstance)

	if !ok {
		payload := handler.GetInitPayload(ctx)
		if payload == nil {
			return nil
		}
		b, err := json.Marshal(payload)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(b))
		accessToken := payload.GetString("accessToken")

		var serviceInstance ServiceInstance
		fmt.Printf("dooppappapapapapapapapapapapapapapapapapapapapapapaapaapapapapapapapapa")
		fmt.Printf(accessToken)

		app.db.First(&serviceInstance, "access_token = ?", accessToken)

		return &serviceInstance
	}

	return raw
}
func (app *App) RunHTTPServer() {
	// Setting up Gin
	app.router = gin.New()
	app.router.RedirectTrailingSlash = false
	app.router.Use(cors.Default(), app.CreateServiceInstanceMiddleware())
	app.router.Any("/query", app.graphqlQueryHandler())
	app.router.GET("/playground", app.graphqlPlaygroundHandler())

	go func() {
		for {
			time.Sleep(time.Second)

			log.Println("Checking if started...")
			resp, err := http.Get("http://localhost:2137/playground")
			if err != nil {
				log.Println("Failed:", err)
				continue
			}
			resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				log.Println("Not OK:", resp.StatusCode)
				continue
			}

			// Reached this point: server is up and running!
			break
		}

		// Start all instances after the server is up and running
		log.Println("Startin up!")
		app.sl.StartupAllInstances()
	}()

	app.router.Run(":2137")
}

func (app *App) graphqlQueryHandler() gin.HandlerFunc {
	h := handler.GraphQL(NewExecutableSchema(Config{Resolvers: &Resolver{
		App: app,
	}}), handler.WebsocketUpgrader(websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}))

	return func(c *gin.Context) {

		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func (app *App) graphqlPlaygroundHandler() gin.HandlerFunc {
	h := handler.Playground("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
