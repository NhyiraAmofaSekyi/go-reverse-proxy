package cli

import (
	"fmt"
	"log"
	"strconv"

	"github.com/NhyiraAmofaSekyi/go-reverse-proxy/internal/configs"
	"github.com/NhyiraAmofaSekyi/go-reverse-proxy/internal/server"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rp",
	Short: "rp is a reverse proxy CLI tool",
	Long:  "rp is a CLI tool for managing reverse proxy configurations",
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommands are provided, display the help message
		cmd.Help()
	},
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage reverse proxy servers",
	Long:  "Manage reverse proxy servers in the configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommands are provided for "server", display the help message
		cmd.Help()
	},
}

var listServersCmd = &cobra.Command{
	Use:   "list",
	Short: "List reverse proxy servers",
	Long:  "List reverse proxy servers defined in the configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := configs.ReadConfig()
		if err != nil {
			log.Fatalf("Error reading configuration: %v", err)
		}

		fmt.Println("Reverse Proxy Servers:")
		for i, server := range config.Resources {
			fmt.Printf("%d. Name: %s, Endpoint: %s, Destination URL: %s\n", i+1, server.Name, server.Endpoint, server.Destination_URL)
		}
	},
}

var updateServerCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a reverse proxy server",
	Long:  "Update a reverse proxy server in the configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 4 {
			log.Fatal("Usage: rp server update [index] [name] [endpoint] [destination_url]")
		}

		// Parse the index argument from the command-line arguments
		index, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalf("Invalid index: %v", err)
		}

		// Read the current configuration from the file
		config, err := configs.ReadConfig()
		if err != nil {
			log.Fatalf("Error reading configuration: %v", err)
		}

		// Check if the index is within the valid range
		if index < 1 || index > len(config.Resources) {
			log.Fatal("Invalid index")
		}

		// Update the specified reverse proxy server with the new values
		config.Resources[index-1] = configs.Resource{
			Name:            args[1],
			Endpoint:        args[2],
			Destination_URL: args[3],
		}

		// Save the updated configuration back to the file
		if err := configs.SaveConfig(config); err != nil {
			log.Fatalf("Error saving configuration: %v", err)
		}

		// Print a success message
		fmt.Println("Configuration updated successfully.")

	},
}

var deleteServerCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a reverse proxy server",
	Long:  "Delete a reverse proxy server from the configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		// Check if the correct number of arguments is provided
		if len(args) != 1 {
			log.Fatal("Usage: rp server delete [index]")
		}

		// Parse the index argument from the command-line arguments
		index, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalf("Invalid index: %v", err)
		}

		// Read the current configuration from the file
		config, err := configs.ReadConfig()
		if err != nil {
			log.Fatalf("Error reading configuration: %v", err)
		}

		// Check if the index is within the valid range
		if index < 1 || index > len(config.Resources) {
			log.Fatal("Invalid index")
		}

		// Remove the specified reverse proxy server from the configuration
		config.Resources = append(config.Resources[:index-1], config.Resources[index:]...)

		// Save the updated configuration back to the file
		if err := configs.SaveConfig(config); err != nil {
			log.Fatalf("Error saving configuration: %v", err)
		}

		// Print a success message
		fmt.Println("Configuration deleted successfully.")

	},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the reverse proxy server",
	Long:  "Run the reverse proxy server based on the configuration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running the reverse proxy server...")
		if err := server.Run(); err != nil {
			log.Fatalf("could not start the server: %v", err)
		}
	},
}

func init() {
	// Add the "server" subcommand
	rootCmd.AddCommand(serverCmd)

	// Add subcommands for the "server" command
	serverCmd.AddCommand(listServersCmd)
	serverCmd.AddCommand(updateServerCmd)
	serverCmd.AddCommand(deleteServerCmd)

	// Add the "run" subcommand
	rootCmd.AddCommand(runCmd)

}

// Execute the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
