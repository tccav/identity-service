package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Configs struct {
	Environment string `envconfig:"ENVIRONMENT" default:"dev"`
	Auth        auth
	API         api
	DB          db
	MemoryDB    memoryDB
	Kafka       kafka
	Swagger     swagger
}

type api struct {
	Port         int           `envconfig:"API_PORT" default:"8000"`
	ReadTimeout  time.Duration `envconfig:"API_READ_TIMEOUT" default:"15s"`
	WriteTimeout time.Duration `envconfig:"API_WRITE_TIMEOUT" default:"15s"`
	IdleTimeout  time.Duration `envconfig:"API_IDLE_TIMEOUT" default:"1m"`
}

type auth struct {
	Secret   string        `envconfig:"TOKEN_SECRET" required:"true"`
	Issuer   string        `envconfig:"TOKEN_ISSUER" default:"uerj"`
	Duration time.Duration `envconfig:"TOKEN_DURATION" default:"3h"`
}

func (a auth) TokenSecret() string {
	return a.Secret
}

func (a auth) TokenIssuer() string {
	return a.Issuer
}

func (a auth) TokenDuration() time.Duration {
	return a.Duration
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

type memoryDB struct {
	Host     string `envconfig:"MEMORY_DB_HOST" required:"true"`
	Port     string `envconfig:"MEMORY_DB_PORT" required:"true"`
	User     string `envconfig:"MEMORY_DB_USER"`
	Password string `envconfig:"MEMORY_DB_PASSWORD"`
}

func (d memoryDB) URL() string {
	return fmt.Sprintf("%s:%s", d.Host, d.Port)
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
