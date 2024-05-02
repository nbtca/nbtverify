package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nbtca/zportal-web-verify/config"
	"github.com/nbtca/zportal-web-verify/nbtverify"
	"github.com/nbtca/zportal-web-verify/nbtverify/utils"
)

var mobile = true

type OnlineStatus int

const (
	AlreadyOnline OnlineStatus = iota
	Offline
	LoginSuccess
	LoginFailed
)

func login(cfg *config.Config, force bool) (OnlineStatus, error) {
	find, address, err := nbtverify.GetBaseLoginUrl(mobile)
	if err != nil {
		return LoginFailed, err
	}
	if !find {
		if !force {
			return AlreadyOnline, nil
		}
		if cfg.CacheFile == "" {
			return AlreadyOnline, nil
		}
		if err := utils.FileNotExists(cfg.CacheFile); err != nil {
			return AlreadyOnline, nil
		}
		// get base url from file
		bytes, err := os.ReadFile(cfg.CacheFile)
		if err != nil {
			return AlreadyOnline, nil
		}
		address = string(bytes)
	}
	if cfg.CacheFile != "" {
		if err := utils.FileNotExists(cfg.CacheFile); err != nil {
			return LoginFailed, err
		}
		// save base url to file
		if err := os.WriteFile(cfg.CacheFile, []byte(address), 0644); err != nil {
			return LoginFailed, err
		}
	}
	fmt.Println("found login base url:", address)
	info := nbtverify.LoginInfo{
		Username: cfg.Username,
		Password: cfg.Password,
		AsMobile: mobile,
	}
	v, err := nbtverify.Login(address, info)
	if err != nil {
		return LoginFailed, err
	}
	fmt.Println("login result:", v)
	detail, err := v.GetDetail()
	if err != nil {
		return LoginFailed, err
	}
	fmt.Println("detail", *detail)
	result, err := detail.Logout()
	if err != nil {
		return LoginFailed, err
	}
	fmt.Println("result", result.Message)
	return LoginSuccess, nil
}

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "c", "config.json", "path of config file")
}
func main() {
	flag.Parse()
	if configPath == "config.json" {
		fmt.Println("use default config path: config.json")
	}
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch status, err := login(cfg, false); status {
	case AlreadyOnline:
		fmt.Println("already online")
	case Offline:
		fmt.Println("offline")
	case LoginSuccess:
		fmt.Println("login success")
	case LoginFailed:
		fmt.Println("login failed")
		fmt.Println(err)
	}
}
