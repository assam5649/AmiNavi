package main

import (
	"server/internal/db"
)

func main() {
	db.Connect()
	db.Migrate()

	// if err := server.Start(); err != nil {
	// 	log.Fatal(err)
	// }
}
