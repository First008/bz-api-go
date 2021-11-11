package main

import (
	"bulutzincir/api"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	fmt.Println(os.Environ())

	api.Gin()
}
