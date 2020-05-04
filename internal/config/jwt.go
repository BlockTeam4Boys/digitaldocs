package config

import "time"

type JWT struct {
	Token Token  `yaml:"token"`
	Key   string `yaml:"key"`
}

type Token struct {
	CookieName  string        `yaml:"cookie_name"`
	TTL         time.Duration `yaml:"ttl"`
	RefreshTime time.Duration `yaml:"refresh_time"`
}
