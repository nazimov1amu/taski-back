package main

import (
	"log"

	"taski_backend/internal/app"
)

func main() {
	mcpApp, err := app.NewMCP()
	if err != nil {
		log.Fatal(err)
	}
	defer mcpApp.Close()

	log.Printf("listening on %s", mcpApp.Addr())
	if err := mcpApp.Run(); err != nil {
		log.Fatal(err)
	}
}
