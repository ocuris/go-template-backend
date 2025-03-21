package config

import (
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type (
	ImmutableConfigs interface {
		GetPort() int
		GetDBConf() DB
		GetMongoConf() MongoDB
		GetRedisConf() Redis
	}

	config struct {
		AppPort int     `mapstructure:"APP_PORT"`
		DB      DB      `mapstructure:"DB"`
		MongoDB MongoDB `mapstructure:"MONGO_DB"`
		Redis   Redis   `mapstructure:"REDIS"`
	}

	DB struct {
		Host             string `mapstructure:"HOST"`
		Port             string `mapstructure:"PORT"`
		Name             string `mapstructure:"NAME"`
		User             string `mapstructure:"USER"`
		Password         string `mapstructure:"PASSWORD"`
		MaxIdleConns     int    `mapstructure:"MAX_IDLE_CONNS"`
		MaxOpenConns     int    `mapstructure:"MAX_OPEN_CONNS"`
		MaxLifetimeConns int    `mapstructure:"MAX_LIFETIME_CONNS"`
		SSLMode          string `mapstructure:"SSL_MODE"`
	}

	MongoDB struct {
		ConnectionURI     string `mapstructure:"CONNECTION_URI"`
		Name              string `mapstructure:"NAME"`
		MaxWriteBatchSize int    `mapstructure:"MAX_WRITE_BATCH_SIZE"`
	}

	Redis struct {
		Host     string `mapstructure:"HOST"`
		Name     string `mapstructure:"NAME"`
		Password string `mapstructure:"PASSWORD"`
	}
)

func (im *config) GetPort() int {
	return im.AppPort
}

func (im *config) GetDBConf() DB {
	return im.DB
}

func (im *config) GetMongoConf() MongoDB {
	return im.MongoDB
}

func (im *config) GetRedisConf() Redis {
	return im.Redis
}

var (
	once sync.Once
	conf *config
)

func NewImmutableConfigs() ImmutableConfigs {
	once.Do(func() {
		v := viper.New()
		appEnv, exists := os.LookupEnv("APP_ENV")
		configName := "app.config.local"
		if exists {
			switch appEnv {
			case "development":
				configName = "app.config.dev"
			case "production":
				configName = "app.config.prod"
			}
		}

		slog.Debug("Config loaded", slog.String("ConfigName", configName), slog.String("Level", "warn"))

		v.SetConfigName("configs/" + configName)
		v.AddConfigPath(".")

		v.SetEnvPrefix("GO_TEMPLATE")
		v.AutomaticEnv()

		if err := v.ReadInConfig(); err != nil {
			panic(err.Error())
		}

		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		err := v.Unmarshal(&conf)
		if err != nil {
			panic(err.Error())
		}
	})
	return conf
}
