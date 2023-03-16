package config

import (
	"github.com/thomasgouveia/go-config"
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

var loaderOptions = &config.Options[Config]{
	Format: config.YAML,

	// Configure the loader to lookup for environment
	// variables with the following pattern: APP_*
	EnvEnabled: true,
	EnvPrefix:  "app",

	// Configure the loader to search for an alpha.yaml file
	// inside one or more locations defined in `FileLocations`
	FileName:      "registry",
	FileLocations: []string{"/etc/morty-registry", "$HOME/.morty", "."},

	// Inject a default configuration in the loader
	Default: &Config{
		Port: 8080,
		Storage: StorageConfig{
			S3: S3Config{
				Bucket:   "functions",
				Region:   "eu-west-1",
				Endpoint: "http://localhost:9000",
			},
		},
	},
}

// Load the configuration from the different sources (default, file and environment).
func Load() (*Config, error) {
	cl, err := config.NewLoader(loaderOptions)
	if err != nil {
		return nil, err
	}
	return cl.Load()
}
