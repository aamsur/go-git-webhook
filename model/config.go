package model

import (
	"encoding/json"
	"webhook/helper"
	"io/ioutil"
)

type ConfigFilter struct {
	Branch string
	Ref    string
	Path   string
}

//ConfigRepository represents a repository from the config file
type ConfigRepository struct {
	Name, Url, Event string
	Commands         []string
	Filters          []ConfigFilter
}

type ConfigDeploy struct {
	Before, After string
}

//Config represents the config file
type Config struct {
	Logfile      string
	Host         string
	Port         int64
	Deploy       ConfigDeploy
	Repositories []ConfigRepository
}

func LoadConfig(configFile string) (config Config) {
	file, e := ioutil.ReadFile(configFile)
	if e != nil {
		helper.PanicIf(e)
	}

	e = json.Unmarshal(file, &config)
	helper.PanicIf(e)

	return
}