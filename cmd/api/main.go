package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sub-rat/social_network_api/internals/server"
)

func main() {
	fmt.Println("Starting Service")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println(os.Getenv("DB_NAME"))
	srv := server.GetServer()
	srv.Run()
}
