package models

type EnvsConfig struct {
	Domain             string
	RedirectTo         string
	Host               string
	Port               string
	SSL                string
	DBURL              string
	CookiesAuthSecret  string
	GoogleClientID     string
	GoogleClientSecret string
	GithubClientID     string
	GithubClientSecret string
	GithubApiToken     string
}

type Configuration struct {
	MaxDocumentationCount int64
}
