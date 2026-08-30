package main

import (
	"log"

	"taski_backend/internal/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatal(err)
	}
	defer application.Close()

	log.Printf("listening on %s", application.Addr())
	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
