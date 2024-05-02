package config

import (
	"encoding/json"
	"os"

	"github.com/nbtca/zportal-web-verify/nbtverify/utils"
)

func LoadConfig(path string, cfg *Config) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	bytes = utils.RemoveComments(bytes)
	err = json.Unmarshal(bytes, cfg)
	if err != nil {
		return err
	}
	return nil
}
