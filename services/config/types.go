package config

type EnvsConfig struct {
	PublicHost         string
	OrgDomain          string
	Host               string
	Port               string
	SSL                string
	DBURL             string
	CookiesAuthSecret  string
	GoogleClientID     string
	GoogleClientSecret string
	GithubClientID     string
	GithubClientSecret string
	ConfigPath         string
}

type CORSConfig struct {
	AllowedOrigins   map[string]bool
	// Duration of preflight caching
	MaxAge           int
}
