package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/nbtca/nbtverify/config"
	"github.com/nbtca/nbtverify/nbtverify"
	"github.com/nbtca/nbtverify/nbtverify/utils"
)

var mobile = true
var cfg *config.Config = new(config.Config)

const defaultPingUrl = "http://10.80.92.85/"

func login(address string) (*nbtverify.OnlineDetail, error) {
	v, err := nbtverify.Login(address, nbtverify.LoginInfo{
		Username: cfg.Username,
		Password: cfg.Password,
		AsMobile: mobile,
	})
	if err != nil {
		return nil, err
	}
	if v.Success {
		fmt.Println("login success")
	} else {
		fmt.Println("login failed:", v.Message)
	}
	detail, err := v.GetDetail()
	if err != nil {
		return nil, err
	}
	if cfg.StatusFile != "" {
		data := map[string]interface{}{
			"result": *v,
			"detail": *detail,
		}
		if err := utils.SaveJson(cfg.StatusFile, data); err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("-----------------Detail-----------------")
		fmt.Println("Welcome:\t", detail.Welcome)
		fmt.Println("Account:\t", detail.Account)
		fmt.Println("UserName:\t", detail.UserName)
		fmt.Println("MAC:\t", detail.UserMac)
		fmt.Println("UserIP:\t", detail.UserIP)
		fmt.Println("DeviceIP:\t", detail.DeviceIP)
		fmt.Println("IsMacFastAuth:\t", detail.IsMacFastAuth)
		fmt.Println("-----------------End-----------------")
	}
	return detail, nil

}
func loadBaseUrl(force bool) (bool, *string, error) {
	find, address, err := nbtverify.GetBaseLoginUrl(cfg.PingUrl, mobile)
	if err != nil {
		return find, nil, err
	}
	if !find {
		if !force {
			return find, nil, nil
		}
		if cfg.CacheFile == "" {
			return find, nil, nil
		}
		if err := utils.FileNotExists(cfg.CacheFile); err != nil {
			return find, nil, nil
		}
		// get base url from file
		bytes, err := os.ReadFile(cfg.CacheFile)
		if err != nil {
			return find, nil, nil
		}
		address = string(bytes)
	}
	if cfg.CacheFile != "" {
		if err := utils.FileNotExists(cfg.CacheFile); err != nil {
			return find, nil, err
		}
		// save base url to file
		if err := os.WriteFile(cfg.CacheFile, []byte(address), 0644); err != nil {
			return find, nil, err
		}
	}
	return find, &address, nil
}
func service() {
	for {
		find, baseUrl, err := loadBaseUrl(false)
		if err != nil {
			fmt.Println(err)
			fmt.Println("error occurs, sleep 30s.")
			time.Sleep(30 * time.Second)
			continue
		}
		if find {
			// offline and got login url
			if baseUrl != nil {
				detail, err := login(*baseUrl)
				if err != nil {
					fmt.Println(err)
				}
				if detail != nil {
					// success
					time.Sleep(10 * time.Millisecond)
					continue
				}
			}
		}
		time.Sleep(1 * time.Second)
	}
}

var (
	configPath string
	times      int
	force      bool
)

func init() {
	flag.StringVar(&configPath, "c", "config.json", "path of config file")
	flag.IntVar(&times, "t", 1, "retry times if login failed")
	flag.BoolVar(&force, "f", false, "force login use url from cache file")
	flag.StringVar(&cfg.Username, "u", "", "username")
	flag.StringVar(&cfg.Password, "p", "", "password")
	flag.StringVar(&cfg.CacheFile, "url", "url.txt", "cache file path")
	flag.StringVar(&cfg.StatusFile, "s", "", "status json file path")
	flag.BoolVar(&mobile, "mobile", true, "use mobile login")
	flag.StringVar(&cfg.PingUrl, "ping", defaultPingUrl, "ping url")
}

func main() {
	flag.Parse()
	action := flag.Arg(0)
	if cfg.Username == "" || cfg.Password == "" {
		if configPath == "config.json" {
			fmt.Println("use default config path: config.json")
		}
		err := config.LoadConfig(configPath, cfg)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if action == "service" {
		service()
		return
	}
	find, baseUrl, err := loadBaseUrl(false)
	if err != nil {
		fmt.Println(err)
		return
	}
	reloadBaseUrlForce := func() error {
		_, baseUrlFromCache, err := loadBaseUrl(true)
		if err != nil {
			fmt.Println(err)
			return err
		}
		baseUrl = baseUrlFromCache
		if baseUrl == nil {
			fmt.Println("base url not found, please create cache file: ", cfg.CacheFile)
			return err
		}
		return nil
	}
	if baseUrl != nil {
		fmt.Println("get login base url:", *baseUrl)
	}
	//match action
	if action == "login" || action == "" {
		if !find && !force {
			fmt.Println("Already online. Tips: use '-f' to get status or 'logout' to logout.")
			return
		}
		if baseUrl == nil {
			if !force {
				fmt.Println("base url not found, please use -f to force login use url from cache file")
				return
			}
			if reloadBaseUrlForce() != nil {
				return
			}
		}
		for i := 0; i < times; i++ {
			login(*baseUrl)
		}
	} else if action == "logout" {
		if baseUrl == nil {
			if reloadBaseUrlForce() != nil {
				return
			}
		}
		detail, err := login(*baseUrl)
		if err != nil {
			fmt.Println(err)
			return
		}
		for i := 0; i < times; i++ {
			result, err := detail.Logout()
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("logout success")
				fmt.Println(result.Message)
			}
		}
	} else if action == "relogin" {
		if baseUrl == nil {
			if reloadBaseUrlForce() != nil {
				return
			}
		}
		for i := 0; i < times; i++ {
			detail, err := login(*baseUrl)
			if err != nil {
				fmt.Println(err)
				return
			}
			result, err := detail.Logout()
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("logout success")
				fmt.Println(result.Message)
			}
			login(*baseUrl)
		}
	} else if action == "flogin" {
		if baseUrl == nil {
			if reloadBaseUrlForce() != nil {
				return
			}
		}
		detail, err := login(*baseUrl)
		if err != nil {
			fmt.Println(err)
			return
		}
		detail.Logout()
		for i := 0; i < times; i++ {
			login(*baseUrl)
		}
		detail.Logout()
	} else {
		fmt.Println("unknown action: " + action)
		fmt.Println("avaliable actions: login logout relogin")
	}
}
