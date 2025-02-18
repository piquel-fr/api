package config

type EnvsConfig struct {
	PublicHost         string
	OrgDomain          string
	Host               string
	Port               string
	SSL                string
	DBURL              string
	CookiesAuthSecret  string
	GoogleClientID     string
	GoogleClientSecret string
	GithubClientID     string
	GithubClientSecret string
	ConfigPath         string
}

type ConfigType struct {
	CORS CORSConfig `yaml:"cors"`
}

type CORSConfig struct {
	AllowedOrigins map[string]bool `yaml:"allowed_origins"`
	MaxAge         int             `yaml:"max_age"`
}
