package config

import "time"

type JWT struct {
	AccessToken  Token `yaml:"access_token"`
	RefreshToken Token `yaml:"refresh_token"`
}

type Token struct {
	CookieName string        `yaml:"cookie_name"`
	TTL        time.Duration `yaml:"ttl"`
}
