package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type (
	Config struct {
		Static     string     `yaml:"static"`
		Db         Db         `yaml:"db"`
		InstallK8s InstallK8s `yaml:"install-k8s"`
	}
	Db []struct {
		Name         string `yaml:"name"`
		Type         string `yaml:"type"`
		Dsn          string `yaml:"dsn"`
		MaxOpenConns int    `yaml:"set_max_open_conns"`
		TablePrefix  string `yaml:"table_prefix"`
	}
	InstallK8s struct {
		SourceDir string `yaml:"source-dir"`
	}
)

func (c *Config) ReadConfigFile(fileName string) error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal([]byte(data), c); err != nil {
		return err
	}
	return nil
}

func (c *Config) WriteConfigFile() error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	execPath, err := os.Getwd()
	if err != nil {
		return err
	}

	configPath := execPath + "/etc/config.yaml"
	err = os.WriteFile(configPath, data, 0755)
	if err != nil {
		return err
	}

	return nil
}
