package config

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type (
	Config struct {
		Port    int           `yaml:"port"`
		Storage StorageConfig `yaml:"storage"`
	}

	StorageConfig struct {
		S3 S3Config `yaml:"s3"`
	}

	S3Config struct {
		Bucket   string `yaml:"bucket"`
		Region   string `yaml:"region"`
		Endpoint string `yaml:"endpoint"`
	}
)

var defaultConfig = &Config{
	Port: 8080,
	Storage: StorageConfig{
		S3: S3Config{
			Bucket:   "functions",
			Region:   "eu-west-1",
			Endpoint: "http://localhost:9000",
		},
	},
}

const (
	configFile = "registry"
	configType = "yaml"
)

// Parse load the configuration from the different sources (default, file and environment).
func Parse() (*Config, error) {
	v := viper.New()

	v.SetConfigName(configFile)
	v.SetConfigType(configType)

	// Set the defaults values in Viper.
	// Viper needs to know if a key exists in order to override it.
	b, err := yaml.Marshal(defaultConfig)
	if err != nil {
		return nil, err
	}

	re := bytes.NewReader(b)
	if err := v.MergeConfig(re); err != nil {
		return nil, err
	}

	// Overwrite values from configuration files
	// This will check for a application.yaml file in the directories listed below
	for _, path := range []string{"/etc/morty-registry", "$HOME/.morty-registry", "."} {
		v.AddConfigPath(path)
	}

	v.SetEnvPrefix("registry")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Parse configuration from all sources
	if err := v.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	// Finally unmarshal the Viper loaded configuration into our config struct
	config := defaultConfig
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return config, nil
}
