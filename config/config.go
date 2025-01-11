package config

import (
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Kafka struct {
		Brokers           string `yaml:"brokers" validate:"required"`
		Topic             string `yaml:"topic" validate:"required"`
		GroupID           string `yaml:"group_id" validate:"required"`
		SchemaregistryUrl string `yaml:"schemaregistry_url" validate:"required"`
	} `yaml:"kafka"`
	Database struct {
		DatabaseSource string `yaml:"database_source"`
		DatabaseTarget string `yaml:"database_target"`
	} `yaml:"database"`
	Logging struct {
		Info  *bool `yaml:"info" validate:"boolean"`
		Debug *bool `yaml:"debug" validate:"boolean"`
		File  string `yaml:"file"`
	} `yaml:"logging"`
}

func (cfg *Config) loadFromYaml(path string) {
	rootDir, err := os.Getwd()

	if err != nil {
		panic(err)
	}
	yamlPath := filepath.Join(rootDir, path)
	file, err := os.Open(yamlPath)

	if err != nil {
		panic(err)
	}

	defer file.Close()
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		panic("Invalid Yaml File")
	}
}

func (cfg *Config) ApplyDefaults() {
	if cfg.Logging.Info == nil {
		defaultValue := false
		cfg.Logging.Info = &defaultValue
	}
	if cfg.Logging.Debug == nil {
		defaultValue := false
		cfg.Logging.Debug = &defaultValue
	}
}

func NewConfig() *Config {

	config := &Config{}

	config.loadFromYaml("config.yaml")

	config.ApplyDefaults()

	if err := config.Validate(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return config
}

func (cfg *Config) Validate() error {
	validate := validator.New()
	return validate.Struct(cfg)
}
