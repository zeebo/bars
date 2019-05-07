package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/zeebo/errs"
)

type Config struct {
	Owner   string   `json:"owner"`
	Repo    string   `json:"repo"`
	Mergers []string `json:"mergers"`
	Token   string   `json:"token"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	cfg := new(Config)
	return cfg, errs.Wrap(json.Unmarshal(data, cfg))
}
