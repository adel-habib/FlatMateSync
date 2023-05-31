package main

import (
	"FlatMateSync/config"
	"fmt"
	"log"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	fmt.Println("Database host:", cfg.Database.Host)
}
