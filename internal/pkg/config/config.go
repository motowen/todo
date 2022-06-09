package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/caarlos0/env/v6"
)

func Setup() error {
	err := env.Parse(&Env)
	if err != nil {
		fmt.Printf("config parse fail, %+v\n", err)
		return err
	}

	err = Env.Validate()
	if err != nil {
		fmt.Printf("config validate fail, %+v\n", err)
		return err
	}

	return nil
}

const (
	logLevelError       = "error"
	logLevelDebug       = "debug"
	logLevelWarning     = "warn"
	logLevelInfo        = "info"
	deployEnvDevelop    = "develop"
	deployEnvStage      = "stage"
	deployEnvProduction = "production"
)

var Env EnvVariable

type EnvVariable struct {
	Version                         string   `env:"VERSION" envDefault:"0.4.4"`
	Port                            string   `env:"GO_HTTP_PORT,required"`
	LogLevel                        string   `env:"LOG_LEVEL" envDefault:"INFO"`
	DeployEnvironment               string   `env:"DEPLOY_ENVIRONMENT" envDefault:"DEVELOP"`
	AuthServiceHost                 string   `env:"AUTH_SERVICE_HOST,required"`
	AuthClientID                    string   `env:"AUTH_CLIENT_ID,required"`
	AuthClientSecret                string   `env:"AUTH_CLIENT_SECRET,required"`
	AuthServiceCacheTTL             int64    `env:"AUTH_SERVICE_CACHE_TTL,required"`
	RedisType                       string   `env:"REDIS_TYPE,required"`
	RedisEndpointList               []string `env:"REDIS_ENDPOINT_LIST,required"`
	RedisPassword                   string   `env:"REDIS_PASSWORD,required"`
	MongoURI                        string   `env:"MONGO_URI,required"`
	PostgresHost                    string   `env:"POSTGRES_HOST,required"`
	PostgresPort                    string   `env:"POSTGRES_PORT,required"`
	PostgresUsername                string   `env:"POSTGRES_USERNAME,required"`
	PostgresPassword                string   `env:"POSTGRES_PASSWORD,required"`
	PostgresName                    string   `env:"POSTGRES_NAME,required"`
	PostgresMinConnSize             int32    `env:"POSTGRES_MIN_CONN_SIZE" envDefault:"0"`
	PostgresMaxConnSize             int32    `env:"POSTGRES_MAX_CONN_SIZE" envDefault:"64"`
	PostgresMaxConnIdleTimeBySecond int64    `env:"POSTGRES_CONN_IDLE_TIME_BY_SECOND" envDefault:"1"`
	PostgresMaxConnLifeTimeBySecond int64    `env:"POSTGRES_CONN_LIFE_TIME_BY_SECOND" envDefault:"60"`
	NatsUrl                         string   `env:"NATS_URL,required"`
	ElasticsearchIndexPrefix        string   `env:"ELASTICSEARCH_INDEX_PREFIX,required"`
	ElasticsearchUrl                string   `env:"ELASTICSEARCH_URL,required"`
}

func (env EnvVariable) Validate() (err error) {
	port, err := strconv.ParseUint(env.Port, 10, 16)
	if err != nil || port <= 0 || port > uint64(65535) {
		err = errors.New("required environment variable \"GO_HTTP_PORT\" should be 0~65535")
		return
	}
	if d := strings.ToLower(env.LogLevel); d != logLevelError && d != logLevelDebug && d != logLevelWarning && d != logLevelInfo {
		err = errors.New("required environment variable \"LOG_LEVEL\" should be \"ERROR|DEBUG|WARN|INFO\"")
		return
	}
	if d := strings.ToLower(env.DeployEnvironment); d != deployEnvDevelop && d != deployEnvStage && d != deployEnvProduction {
		err = errors.New("required environment variable \"DEPLOY_ENVIRONMENT\" should be \"DEVELOP|STAGE|PRODUCTION\"")
		return
	}

	return
}

func IsProduction() bool {
	return strings.ToLower(Env.DeployEnvironment) == deployEnvProduction
}
