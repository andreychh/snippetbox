package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log/slog"
	"net/url"
	"os"
)

type Config struct {
	EnvName  string   `yaml:"envName"`
	Database database `yaml:"database"`
	App      app      `yaml:"app"`
	Logging  logging  `yaml:"logging"`
}

func MustLoad(configPath string) *Config {
	config := &Config{}

	data, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("reading file: %v", err))
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		panic(fmt.Sprintf("unmarshaling: %v", err))
	}

	return config
}

type database struct {
	User     string            `yaml:"user"`
	Password string            `yaml:"password"`
	Protocol string            `yaml:"protocol"`
	Host     string            `yaml:"host"`
	Port     int               `yaml:"port"`
	Name     string            `yaml:"name"`
	Params   map[string]string `yaml:"params"`
}

func (db database) DSN() string {
	var dsn = fmt.Sprintf("%s:%s@%s(%s:%d)/%s",
		db.User, db.Password, db.Protocol, db.Host, db.Port, db.Name,
	)

	if len(db.Params) == 0 {
		return dsn
	}

	var query = url.Values{}
	for k, v := range db.Params {
		query.Add(k, v)
	}

	return fmt.Sprintf("%s?%s", dsn, query.Encode())
}

type app struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (a app) Addr() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}

type logging struct {
	Level slog.Level `yaml:"level"`
}
