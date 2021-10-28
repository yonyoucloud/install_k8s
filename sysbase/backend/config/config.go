package config

import (
	"io/ioutil"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type (
	Config struct {
		Static     string     `yaml:"static"`
		Mysql      Mysql      `yaml:"mysql"`
		InstallK8s InstallK8s `yaml:"install-k8s"`
	}
	Mysql struct {
		MasterDsn       string        `yaml:"master-dsn"`
		SourcesDsn      []string      `yaml:"sources-dsn"`
		ReplicasDsn     []string      `yaml:"replicas-dsn"`
		MaxIdleConns    int           `yam:"set-max-idle-conns"`
		MaxOpenConns    int           `yam:"set-max-open-conns"`
		ConnMaxLifetime time.Duration `yaml:"set-conn-max-lifetime"`
	}
	InstallK8s struct {
		SourceDir string `yaml:"source-dir"`
	}
)

func (c *Config) ReadConfigFile(fileName string) error {
	data, err := ioutil.ReadFile(fileName)
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
	err = ioutil.WriteFile(configPath, data, 0755)
	if err != nil {
		return err
	}

	return nil
}
