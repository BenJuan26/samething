package config

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
)

type config struct {
	User   string `json:"user"`
	DBName string `json:"db_name"`
}

var configData *config

func init() {
	loadConfig("config.json")
}

func loadConfig(path string) {
	_, err := os.Stat(path)
	if err != nil {
		panic("Config file doesn't exist at the specified path")
	}

	if configData == nil {
		configData = new(config)
		buff, err := ioutil.ReadFile(path)
		if err != nil {
			panic(errors.Wrap(err, "Problem reading config file"))
		}
		err = json.Unmarshal(buff, configData)
		if err != nil {
			panic(errors.Wrap(err, "Problem unmarshaling config structure"))
		}
	}
}

func GetUser() string {
	return configData.User
}

func GetDBName() string {
	return configData.DBName
}
