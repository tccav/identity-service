package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Configs struct {
	Environment string `envconfig:"ENVIRONMENT" default:"dev"`
	API         api
	DB          db
	Kafka       kafka
	Swagger     swagger
}

type api struct {
	Port         int           `envconfig:"API_PORT" default:"8000"`
	ReadTimeout  time.Duration `envconfig:"API_READ_TIMEOUT" default:"15s"`
	WriteTimeout time.Duration `envconfig:"API_WRITE_TIMEOUT" default:"15s"`
	IdleTimeout  time.Duration `envconfig:"API_IDLE_TIMEOUT" default:"1m"`
}

type db struct {
	Host     string `envconfig:"DB_HOST" required:"true"`
	Port     string `envconfig:"DB_PORT" required:"true"`
	User     string `envconfig:"DB_USER" required:"true"`
	Password string `envconfig:"DB_PASSWORD" required:"true"`
	Name     string `envconfig:"DB_NAME" required:"true"`
	Options  string `envconfig:"DB_OPTIONS"`
}

func (d db) URL() string {
	u := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", d.User, d.Password, d.Host, d.Port, d.Name)
	if d.Options != "" {
		u += fmt.Sprintf("?%s", d.Options)
	}
	return u
}

type kafka struct {
	Host     string `envconfig:"KAFKA_HOST" required:"true"`
	Port     string `envconfig:"KAFKA_PORT" required:"true"`
	User     string `envconfig:"KAFKA_USER"`
	Password string `envconfig:"KAFKA_PASSWORD"`
}

func (k kafka) URL() string {
	return fmt.Sprintf("%s:%s", k.Host, k.Port)
}

type swagger struct {
	Enabled bool `envconfig:"SWAGGER_ENABLED" default:"false"`
}

func LoadConfigs() (Configs, error) {
	var config Configs
	err := envconfig.Process("", &config)
	if err != nil {
		return Configs{}, err
	}
	return config, nil
}
