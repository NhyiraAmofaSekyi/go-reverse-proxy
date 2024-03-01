package configs

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type Resource struct {
	Name            string
	Endpoint        string
	Destination_URL string
}
type Configuration struct {
	Server struct {
		Host        string
		Listen_port string
	}
	Resources []Resource
}

var Config *Configuration

const configFileName = "config.yaml"

func ReadConfig() (*Configuration, error) {
	// Set the relative path to your config.yaml file
	yamlPath, err := filepath.Abs("../data/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("error getting absolute path: %v", err)
	}

	data, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file: %v", err)
	}

	var config Configuration
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error unmarshalling YAML: %v", err)
	}

	return &config, nil
}

func NewConfiguration() (*Configuration, error) {
	viper.AddConfigPath("data")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading config file: %s", err)
	}
	err = viper.Unmarshal(&Config)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %s", err)
	}
	return Config, nil
}

// createConfigInteractively creates a new configuration interactively from user input
func createConfigInteractively() (*Configuration, error) {
	var config Configuration

	fmt.Print("Enter server host: ")
	fmt.Scanln(&config.Server.Host)

	fmt.Print("Enter server listen port: ")
	fmt.Scanln(&config.Server.Listen_port)

	fmt.Print("Enter the number of resources: ")
	var numResources int
	fmt.Scanln(&numResources)

	config.Resources = make([]Resource, numResources)
	scanner := bufio.NewScanner(os.Stdin)

	for i := 0; i < numResources; i++ {
		fmt.Printf("Enter details for resource %d:\n", i+1)
		fmt.Print("Name: ")
		scanner.Scan()
		config.Resources[i].Name = scanner.Text()

		fmt.Print("Endpoint: ")
		scanner.Scan()
		config.Resources[i].Endpoint = scanner.Text()

		fmt.Print("Destination URL: ")
		scanner.Scan()
		config.Resources[i].Destination_URL = scanner.Text()
	}

	return &config, nil
}

// SaveConfig saves the configuration to the specified file

func SaveConfig(config *Configuration) error {

	yamlPath, err := filepath.Abs("../data/config.yaml")
	if err != nil {
		return fmt.Errorf("error setting file: %v", err)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("error marshalling YAML: %v", err)
	}

	// dir, err := os.Getwd()
	// if err != nil {
	// 	return fmt.Errorf("error getting current working directory: %v", err)
	// }

	// // Set the relative path to the config file
	// configPath := filepath.Join(dir, configFileName)

	if err := os.WriteFile(yamlPath, data, 0644); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}
