package config

import (
	"fmt"
	"log/slog"
	"path"

	"github.com/8thgencore/passfort/pkg/logger/slogpretty"
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Env      string   `yaml:"env" env-defaul:"local" env-required:"true"` // local, dev or prod
		App      App      `yaml:"app"`
		HTTP     HTTP     `yaml:"http"`
		Database Database `yaml:"database"`
		Cache    Cache    `yaml:"cache"`
		Log      Log      `yaml:"log"`
	}

	// App contains all the environment variables for the application
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	// HTTP contains all the environment variables for the http server
	HTTP struct {
		Port           string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
		AllowedOrigins string `env-required:"true" yaml:"allowed_origins"`
	}

	// Database contains all the environment variables for the database
	Database struct {
		Connection string
		Host       string `yaml:"host"     env:"DB_HOST" env-default:"localhost"`
		Port       string `yaml:"port"     env:"DB_PORT" env-default:"5432"`
		User       string `yaml:"user"     env:"DB_USER" env-default:"user"`
		Password   string `yaml:"password" env:"DB_PASSWORD"`
		Name       string `yaml:"name"     env:"DB_NAME" env-default:"postgres"`
	}

	// Cache contains all the environment variables for the cache service
	Cache struct {
		Addr     string `yaml:"address"  env-default:"localhost:6379"`
		Password string `yaml:"password" env-default:"password"`
	}

	// Logger settings
	Log struct {
		Slog Slog `yaml:"slog"`
	}
	Slog struct {
		Level     slog.Level              `yaml:"level"`
		AddSource bool                    `yaml:"add_source"`
		Format    slogpretty.FieldsFormat `yaml:"format"` // json, text or pretty
		Pretty    PrettyLog               `yaml:"pretty"`
	}
	PrettyLog struct {
		FieldsFormat slogpretty.FieldsFormat `yaml:"fields_format"` // json, json-indent or yaml
		Emoji        bool                    `yaml:"emoji"`
		TimeLayout   string                  `yaml:"time_layout"`
	}
)

func NewConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(path.Join("./", configPath), cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	err = cleanenv.UpdateEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("error updating env: %w", err)
	}

	return cfg, nil
}
