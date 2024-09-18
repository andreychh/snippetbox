package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log/slog"
	"net/url"
	"os"
)

type Config struct {
	DB  database `yaml:"database"`
	App app      `yaml:"app"`
	Log logging  `yaml:"logging"`
}

// Load #me загружает конфигурацию из YAML-файла
func Load(configPath string) (*Config, error) {
	var data, err = os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling: %w", err)
	}

	return &config, nil
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

// DSN #me формирует dataSourceName
// Формат строки подключения (dsn): user:password@protocol(host:port)/dbname?param=value
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
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	StaticDir string `yaml:"staticDir"`
}

func (a app) Addr() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}

type logging struct {
	Level slog.Level `yaml:"level"`
}
