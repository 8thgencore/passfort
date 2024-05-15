package config

import (
	"fmt"
	"log/slog"
	"path"
	"time"

	"github.com/8thgencore/passfort/pkg/logger/slogpretty"
	"github.com/ilyakaznacheev/cleanenv"
)

type Env string

const (
	Local Env = "local"
	Dev   Env = "dev"
	Prod  Env = "prod"
)

type (
	Config struct {
		Env            Env            `yaml:"env" env-defaul:"local" env-required:"true"` // local, dev or prod
		App            App            `yaml:"app"`
		HTTP           HTTP           `yaml:"http"`
		Database       Database       `yaml:"database"`
		Cache          Cache          `yaml:"cache"`
		Token          Token          `yaml:"token"`
		MasterPassword MasterPassword `yaml:"master_password"`
		Clients        ClientConfig   `yaml:"clients"`
		Log            Log            `yaml:"log"`
	}

	// App contains all the environment variables for the application
	App struct {
		Name    string `yaml:"name"    env:"APP_NAME"    env-required:"true"`
		Version string `yaml:"version" env:"APP_VERSION" env-required:"true"`
	}

	// HTTP contains all the environment variables for the http server
	HTTP struct {
		Host         string `env-required:"true" yaml:"host" env:"HTTP_HOST"`
		Port         string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
		AllowOrigins string `env-required:"true" yaml:"allow_origins"`
	}

	// Database contains all the environment variables for the database
	Database struct {
		Connection string `yaml:"connection"  env:"DB_CONNECTION" env-default:"postgres"`
		Host       string `yaml:"host"        env:"DB_HOST"       env-default:"localhost"`
		Port       string `yaml:"port"        env:"DB_PORT"       env-default:"5432"`
		User       string `yaml:"user"        env:"DB_USER"       env-default:"user"`
		Password   string `yaml:"password"    env:"DB_PASSWORD"`
		Name       string `yaml:"name"        env:"DB_NAME"       env-default:"postgres"`
	}

	// Cache contains all the environment variables for the cache service
	Cache struct {
		Addr     string `yaml:"address"  env:"REDIS_ADDRESS"  env-default:"localhost:6379"`
		Password string `yaml:"password" env:"REDIS_PASSWORD" env-default:"password"`
	}

	// Token contains all the environment variables for the token service
	Token struct {
		SigningKey      string        `yaml:"signing_key"       env:"TOKEN_SIGNING_KEY"`
		AccessTokenTTL  time.Duration `yaml:"access_token_ttl"  env-default:"30m"`
		RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl" env-default:"720h"`
	}

	// MasterPassword contains all the environment variables for the master password service
	MasterPassword struct {
		MasterPasswordTTL time.Duration `yaml:"master_password_ttl" env-default:"MasterPassword"`
	}

	//  Clients
	Client struct {
		Address      string        `yaml:"address"`
		Timeout      time.Duration `yaml:"timeout"`
		RetriesCount int           `yaml:"retries_count"`
		Insecure     bool          `yaml:"insecure"`
	}
	ClientConfig struct {
		Mail Client `yaml:"mail"`
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
