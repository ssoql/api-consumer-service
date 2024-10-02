package config

import (
	"log"

	"github.com/spf13/viper"
)

const (
	LogLevel   = "info"
	production = "prod"
	develop    = "dev"
)

type Env struct {
	AppEnv            string `mapstructure:"APP_ENV"`
	ApiSeedUrl        string `mapstructure:"API_SEED_URL"`
	ContextTimeout    int    `mapstructure:"CONTEXT_TIMEOUT"`
	RabbitMqUrl       string `mapstructure:"RABBITMQ_URL"`
	RabbitMqQueueName string `mapstructure:"RABBITMQ_QUEUE_NAME"`
}

func (e *Env) IsProduction() bool {
	return e.AppEnv == production
}

func (e *Env) IsDevelop() bool {
	return e.AppEnv == develop
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	if env.AppEnv == "development" {
		log.Println("The App is running in development env")
	}

	return &env
}
