package cli

import (
	"fmt"
	"log"

	"github.com/NhyiraAmofaSekyi/go-reverse-proxy/internal/configs"
	"github.com/NhyiraAmofaSekyi/go-reverse-proxy/internal/server"
	"github.com/spf13/cobra"
)

var configPath string

var rootCmd = &cobra.Command{
	Use:   "rp",
	Short: "rp is a reverse proxy CLI tool",
	Long:  "rp is a CLI tool for managing reverse proxy configurations",
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommands are provided, display the help message
		cmd.Help()
	},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the server with specified configuration",
	Long:  "Run the reverse proxy server based on the specified configuration file.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running the reverse proxy server with configuration...")

		// Load the configuration using the specified config file path
		cfg, err := configs.LoadConfig(configPath)
		if err != nil {
			log.Fatalf("Failed to load configuration: %v", err)
		}

		// Now use 'cfg' to configure and run your server as needed
		if err := server.ConfigRun(cfg); err != nil { // Adjust 'server.Run' to accept config as parameter
			log.Fatalf("Could not start the server: %v", err)
		}
	},
}

func init() {

	runCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to the configuration YAML file")
	runCmd.MarkFlagRequired("config") // Makes the -config flag required

	rootCmd.AddCommand(runCmd)

}

// Execute the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
