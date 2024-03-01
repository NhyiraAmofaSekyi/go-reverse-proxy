package main

import (
	"fmt"
	"log"

	"github.com/NhyiraAmofaSekyi/go-reverse-proxy/cli"
	"github.com/NhyiraAmofaSekyi/go-reverse-proxy/internal/configs"
)

func main() {
	cli.Execute()
}

func testConfig() {
	config, err := configs.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("testing config file")
	// Use the config struct as needed
	fmt.Printf("Server Host: %s\n", config.Server.Host)
	fmt.Printf("Server Listen Port: %s\n", config.Server.Listen_port)

	fmt.Println("Resources:")
	for _, res := range config.Resources {
		fmt.Printf("- Name: %s\n", res.Name)
		fmt.Printf("  Endpoint: %s\n", res.Endpoint)
		fmt.Printf("  Destination URL: %s\n", res.Destination_URL)
	}
}
