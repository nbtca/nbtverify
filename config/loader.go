package config

import (
	"encoding/json"
	"os"

	"github.com/nbtca/zportal-web-verify/nbtverify/utils"
)

func LoadConfig(path string) (*Config, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	bytes = utils.RemoveComments(bytes)
	err = json.Unmarshal(bytes, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
