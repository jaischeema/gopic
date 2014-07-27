package main

import (
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"os"
)

type Config struct {
	data map[string]interface{}
}

func (c *Config) Get(key string) string {
	value := c.data[key].(string)
	return value
}

func setProperEnvironment() {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}
	os.Setenv("GO_ENV", env)
}

func ParseConfig() Config {
	setProperEnvironment()

	content, err := ioutil.ReadFile("./config.yaml")
	CheckError(err, "Invalid configuration file")

	configData := make(map[string]map[string]interface{})
	err = yaml.Unmarshal([]byte(content), &configData)
	CheckError(err, "Invalid configuration options specified")
	return Config{configData[os.Getenv("GO_ENV")]}
}
