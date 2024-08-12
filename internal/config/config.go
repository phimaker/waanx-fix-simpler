package config

import (
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/phimaker/waanx-fix-simpler/internal/logger"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Fix   *Fix `mapstructure:"fix"`
		Db    *Db
		Redis *Redis
	}

	Fix struct {
		ConfigPath string `mapstructure:"config-path"`
		Username   string
		Password   string
	}

	Db struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
		SSLMode  string
		TimeZone string
	}

	Redis struct {
		Host     string
		Port     int
		Password string
		DB       int
	}
)

var (
	once           sync.Once
	configInstance *Config
)

func GetConfig() *Config {

	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			logger.Warn("No .env file found")
		}

		viper.SetConfigName("config") // name of config file (without extension)
		viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
		viper.AddConfigPath(".")      // path to look for the config file in

		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		if err := viper.ReadInConfig(); err != nil {
			logger.Warn("No config file found")
		}

		if err := viper.Unmarshal(&configInstance); err != nil {
			panic(err)
		}
		postInit()
	})

	return configInstance
}

func postInit() {
	// Do some post initialization stuff here
	if configInstance.Fix == nil {
		configInstance.Fix = &Fix{
			ConfigPath: "config.cfg",
		}
	}
}
