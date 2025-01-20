package config

import (
	"os"

	"github.com/vandi37/vanerrors"
	"gopkg.in/yaml.v3"
)

const (
	ErrorToOpenConfig = "error to open config"
	ErrorDecodingData = "error decoding data"
)

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type LogConfig struct {
	Token string
	Chat  int64
}

type Config struct {
	Token     string    `yaml:"token"`
	DB        DBConfig  `yaml:"db"`
	Log       LogConfig `yaml:"log"`
	HashSalt  string    `yaml:"hash_salt"`
	ArgonSalt string    `yaml:"argon_salt"`
}

func Get(path string) (*Config, error) {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, vanerrors.Wrap(ErrorToOpenConfig, err)
	}
	defer file.Close()

	cfg := new(Config)

	err = yaml.NewDecoder(file).Decode(cfg)
	if err != nil {
		return nil, vanerrors.Wrap(ErrorDecodingData, err)
	}

	return cfg, nil
}
