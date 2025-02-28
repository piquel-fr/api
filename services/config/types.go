package config

type EnvsConfig struct {
	Domain             string
	Host               string
	Port               string
	SSL                string
	DBURL              string
	CookiesAuthSecret  string
	GoogleClientID     string
	GoogleClientSecret string
	GithubClientID     string
	GithubClientSecret string
}
