package core

import (
	"log"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func (app *App) RunHTTPServer() {
	// Setting up Gin
	app.router = gin.New()
	app.router.RedirectTrailingSlash = false
	app.router.Use(cors.Default())
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
