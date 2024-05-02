package config

type Config struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	AsMobile  bool   `json:"mobile"`
	CacheFile string `json:"cache"`
}
