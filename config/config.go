package config

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
)

type config struct {
	User    string `json:"user"`
	DBName  string `json:"db_name"`
	Schema  string `json:"schema"`
	BaseURL string `json:"base_url"`
	Title   string `json:"title"`
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

func GetUser() string {
	return configData.User
}

func GetDBName() string {
	return configData.DBName
}

func GetSchema() string {
	return configData.Schema
}

func GetBaseURL() string {
	return configData.BaseURL
}

func GetFullBaseURL() string {
	return configData.Schema + "://" + configData.BaseURL
}

func GetTitle() string {
	return configData.Title
}
