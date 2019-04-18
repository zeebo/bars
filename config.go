package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/zeebo/errs"
)

type Config struct {
	Ownerr    string   `json:"owner"`
	Repo      string   `json:"repo"`
	Reviewers []string `json:"reviewers"`
	TestRef   string   `json:"test_ref"`
	MasterRef string   `json:"master_ref"`
	Token     string   `json:"token"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	cfg := new(Config)
	return cfg, errs.Wrap(json.Unmarshal(data, cfg))
}
