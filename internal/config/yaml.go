package config

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
	Session  Session  `yaml:"session"`
	Logger   Logger   `yaml:"logging"`
}

func New(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("reading file %q: %v", configPath, err)
	}

	config := new(Config)
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling %q: %v", configPath, err)
	}

	return config, nil
}

type Server struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	LogLevel     slog.Level    `yaml:"log-level"`
	IdleTimeout  time.Duration `yaml:"idle-timeout"`
	ReadTimeout  time.Duration `yaml:"read-timeout"`
	WriteTimeout time.Duration `yaml:"write-timeout"`
	CertPath     string        `yaml:"cert-path"`
	KeyPath      string        `yaml:"key-path"`
}

func (s Server) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

type Database struct {
	User     string            `yaml:"user"`
	Password string            `yaml:"password"`
	Name     string            `yaml:"name"`
	Protocol string            `yaml:"protocol"`
	Host     string            `yaml:"host"`
	Port     int               `yaml:"port"`
	Params   map[string]string `yaml:"params"`
}

func (db Database) DSN() string {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s",
		db.User, db.Password, db.Protocol, db.Host, db.Port, db.Name,
	)

	if len(db.Params) == 0 {
		return dsn
	}

	query := url.Values{}
	for k, v := range db.Params {
		query.Add(k, v)
	}

	return fmt.Sprintf("%s?%s", dsn, query.Encode())
}

type Session struct {
	Lifetime time.Duration `yaml:"lifetime"`
}

type Logger struct {
	Writer    string     `yaml:"writer"`
	Handler   string     `yaml:"handler"`
	Level     slog.Level `yaml:"level"`
	AddSource bool       `yaml:"add-source"`
}
