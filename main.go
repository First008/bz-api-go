package main

import (
	"bulutzincir/api"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	//inits.Init()
	api.Gin()

}
func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}
