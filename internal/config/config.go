package config

import (
	"os"

	"github.com/palantir/go-baseapp/baseapp"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server  baseapp.HTTPConfig    `yaml:"server"`
	Logging baseapp.LoggingConfig `yaml:"logging"`
}

func ReadConfig(path string) (Config, error) {
	var c Config

	bytes, err := os.ReadFile(path)
	if err != nil {
		return c, errors.Wrapf(err, "failed reading server config file: %s", path)
	}

	if err := yaml.UnmarshalStrict(bytes, &c); err != nil {
		return c, errors.Wrap(err, "failed parsing configuration file")
	}

	return c, nil
}
