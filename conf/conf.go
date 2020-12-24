package conf

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	DataSourceName string `json:"dataSourceName"`
}

func ReadConfig(filename string) (conf *Config, err error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = json.Unmarshal(content, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
