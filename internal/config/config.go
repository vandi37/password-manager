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
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password"  json:"password"`
	Name     string `yaml:"name"  json:"name"`
}

type Config struct {
	Token     string   `yaml:"token"  json:"token"`
	DB        DBConfig `yaml:"db" json:"db"`
	HashSalt  string   `yaml:"hash_salt" json:"hash_salt"`
	ArgonSalt string   `yaml:"argon_salt" json:"argon_salt"`
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
