package config

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

//TestConfig struct for configuring operator test
type TestConfig struct {
	//Csv is a clusterServiceVersion which contains the packaging details of the operator
	Csv struct {
		//Name of csv  operator package with version name
		Name string `yaml:"name"`
		//Namespace where the operator will be running
		Namespace string `yaml:"namespace"`
		//Expected status of the Csv
		Status string `yaml:"status"`
	} `yaml:"csv"`
}

// NewConfig returns a new decoded Config struct
func NewConfig(configPath string) (*TestConfig, error) {
	// Create config structure
	config := &TestConfig{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

// ValidateConfigPath just makes sure, that the path provided is a file,
// that can be read
func validateConfigPath(path string) error {

	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

// ParseFlags will create and parse the CLI flags
// and return the path to be used elsewhere
func parseFlags() (string, error) {
	// String that contains the configured configuration path
	var configPath string
	// Set up a CLI flag called "-config" to allow users
	// to supply the configuration file
	flag.StringVar(&configPath, "config", "config/config.yml", "path to config file")
	// Actually parse the flags
	// flag.Parse() - No need you want ginkgo to call it eventually

	// Validate the path first
	if err := validateConfigPath(configPath); err != nil {
		return "", err
	}

	// Return the configuration path
	return configPath, nil
}

//GetConfig ...Get Operator test config
func GetConfig() (*TestConfig, error) {
	// Generate our config based on the config supplied
	// by the user in the flags
	cfgPath, err := parseFlags()
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := NewConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}
	return cfg, err
}
