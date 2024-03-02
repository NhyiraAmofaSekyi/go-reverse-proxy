package configs

import (
	"fmt"
	"os"
	"path/filepath"

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

func LoadConfig(filePath string) (*Configuration, error) {
	// Convert relative path to absolute path (optional step based on use case)
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("error getting absolute path: %v", err)
	}

	data, err := os.ReadFile(absPath) // ioutil.ReadFile is used for simplicity; os.ReadFile can also be used
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file: %v", err)
	}

	var config Configuration
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error unmarshalling YAML: %v", err)
	}

	return &config, nil
}
