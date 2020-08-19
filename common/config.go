package common

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

//Configuration ...
type Configuration struct {
	DATABASE string `yaml:"database"`
	PORT     int    `yaml:"port"`
}

//Config ...
var Config Configuration

//MySQLConn ...
type MySQLConn struct {
	DB *sql.DB
}

//Conn ...
var Conn *MySQLConn

//WebLog ...
var WebLog *os.File

//ReadConfig ...
func ReadConfig(configFileName string) {
	filename, _ := filepath.Abs(configFileName)
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic("Can't read config file: " + configFileName)
	}

	err = yaml.Unmarshal(yamlFile, &Config)

	if err != nil {
		panic("Can't unmarshal config. " + configFileName + " -> " + err.Error())
	}
}
