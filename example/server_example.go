package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lunargon/bolt-tui/src/bolt"
	"github.com/lunargon/bolt-tui/src/gin_server"
)

// RunServer starts the HTMX web server
func RunServer(port string) error {
	// Create a new BoltDB factory
	factory, err := bolt.NewFactory("default", "./bolt.db")
	if err != nil {
		return err
	}
	defer func() {
		// Close all databases when server shuts down
		factory.Close("default")
	}()

	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	// Create router
	r := gin.Default()

	// Setup routes
	gin_server.SetupRoutes(r, factory)

	// Start server
	log.Printf("Starting HTMX server on port %s", port)
	log.Printf("Visit http://localhost:%s to access the web interface", port)

	return r.Run(":" + port)
}

// RunServerWithCustomFactory starts the server with a custom factory
func RunServerWithCustomFactory(factory *bolt.Factory, port string) error {
	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	// Create router
	r := gin.Default()

	// Setup routes
	gin_server.SetupRoutes(r, factory)

	// Start server
	log.Printf("Starting HTMX server on port %s", port)
	log.Printf("Visit http://localhost:%s to access the web interface", port)

	return r.Run(":" + port)
}


func main() {
    err := RunServer("8080")
    if err != nil {
        panic(err)
    }
}