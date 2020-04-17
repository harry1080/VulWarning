package main

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Config -
type Config struct {
	MySQL struct {
		Host    string `yaml:"host"`
		User    string `yaml:"user"`
		Pass    string `yaml:"pass"`
		Name    string `yaml:"name"`
		Charset string `yaml:"charset"`
	} `yaml:"mysql"`
	Pusher string `yaml:"pusher"`
	Server struct {
		Debug bool `yaml:"debug"`
	} `yaml:"server"`
}

var (
	configFileName = "config.yaml"
	conf           Config
)

func templateConfig() (data []byte, err error) {
	conf := &Config{}
	conf.MySQL.Charset = "utf8mb4"
	conf.MySQL.Host = "127.0.0.1"
	conf.MySQL.User = "root"
	conf.MySQL.Pass = "123456"
	conf.MySQL.Name = "vulwarning"
	// yamlData
	data, err = yaml.Marshal(conf)
	if err != nil {
		return
	}
	return data, nil
}

func loadConfig() (err error) {
	var (
		yamlFile []byte
	)
	_, err = os.Stat(configFileName)
	if err != nil && os.IsNotExist(err) {
		if data, _ := templateConfig(); data != nil {
			if err = ioutil.WriteFile(configFileName, data, 0666); err != nil {
				return err
			}
		}
	}
	if yamlFile, err = ioutil.ReadFile(configFileName); err != nil {
		return err
	}
	if err = yaml.Unmarshal(yamlFile, &conf); err != nil {
		return err
	}
	return nil
}
