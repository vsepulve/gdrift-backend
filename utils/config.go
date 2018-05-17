package utils

import (
	"gopkg.in/yaml.v2"

	"io/ioutil"
)

var (
	Config Configuration
)

type Database struct {
	Rdbms     string `yaml:"rdbms"`
	User      string `yaml:"user"`
	Pass      string `yaml:"pass"`
	Protocol  string `yaml: protocol`
	Ip        string `yaml:"ip"`
	Port      string `yaml:"port"`
	Name      string `yaml:"name"`
	Charset   string `yaml:"charset"`
	ParseTime bool   `yaml:"parseTime"`
	Loc       string `yaml:"loc"`
}

type Server struct {
	Port string `yaml:port`
}

type Daemon struct {
	Ip   string `yaml:"ip"`
	Port string `yaml:"port"`
}

type Configuration struct {
	Database Database `yaml:"database"`
	Server   Server   `yaml:"server"`
	Daemon   Daemon   `yaml:"daemon"`
}

func LoadConfig(filename string) {
	bytes, err := ioutil.ReadFile(filename)
	Check(err)

	err = yaml.Unmarshal(bytes, &Config)
	Check(err)
}
