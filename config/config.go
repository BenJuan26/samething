package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

type config struct {
	DBUser     string `json:"db_user"`
	DBName     string `json:"db_name"`
	DBHost     string `json:"db_host"`
	HTTPSchema string `json:"http_schema"`
	WSSchema   string `json:"ws_schema"`
	BaseURL    string `json:"base_url"`
	Title      string `json:"title"`
}

var configData *config

// TODO: Allow for a cmdline flag specifying the config file path
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

func GetDBUser() string {
	return configData.DBUser
}

func GetDBName() string {
	return configData.DBName
}

func GetHTTPSchema() string {
	return configData.HTTPSchema
}

func GetWSSchema() string {
	return configData.WSSchema
}

func GetBaseURL() string {
	return configData.BaseURL
}

func GetFullBaseURL() string {
	return configData.HTTPSchema + "://" + configData.BaseURL
}

func GetTitle() string {
	return configData.Title
}

func GetDBHost() string {
	return configData.DBHost
}
