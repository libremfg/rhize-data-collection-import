package main

import (
	"fmt"
	"log"
	"os"
	"rhize-data-collection-import/commands"

	"github.com/joho/godotenv"
)

const versionString = "rhize-data-collection-import v1.6.0"

func main() {
	log.Println(versionString)
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file, running without")
	}
	if err := commands.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
