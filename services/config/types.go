package config

type Config struct {
	PublicHost         string
	Host               string
	Port               string
	SSL                string
	DB_URL             string
	CookiesAuthSecret  string
	GoogleClientID     string
	GoogleClientSecret string
	GithubClientID     string
	GithubClientSecret string
	ConfigPath         string
}

type Route struct {
	Name       string `yaml:"route"`
	Method     string
	Auth       bool
	Permission string
}
