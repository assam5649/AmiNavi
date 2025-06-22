package main

import (
	"log"
	"server/internal/db"
)

func main() {

	log.Println("------------------")
	db.Connect()
	db.Migrate()

	// if err := server.Start(); err != nil {
	// 	log.Fatal(err)
	// }
}
