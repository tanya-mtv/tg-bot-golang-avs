package config

import (
	"flag"
	"os"
	"tg-bot-golang/internal/logger"

	"github.com/spf13/viper"
)

var configPath string
var configType string

type Config struct {
	Logger    *logger.Config `mapstructure:"logger"`
	BotToken  string
	Redis     *ConfigRedis  `mapstructure:"redis"`
	MSSQL     *ConfigMSSQL  `mapstructure:"mssql"`
	Axelot    *ConfigAxelot `mapstructure:"axelot"`
	RaecKey   string
	RaecValue string
}

func init() {
	flag.StringVar(&configType, "config-type", "", "Format of configuration file type. Supported formats is: yaml")
	flag.StringVar(&configPath, "config", "", "Path to configuration file")
}

func InitConfig() (*Config, error) {
	if configPath == "" {
		configPathFromEnv := os.Getenv(CONFIG_PATH)
		if configPathFromEnv != "" {
			configPath = configPathFromEnv
		} else {
			configPath = cfgPath
		}
	}
	cfg := &Config{}

	configType := os.Getenv(CONFIG_TYPE)

	if configType == "" {
		configType = yaml
	}
	viper.SetConfigType(configType)
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	telegramBotToken := os.Getenv(botToken)
	if telegramBotToken != "" {
		cfg.BotToken = telegramBotToken
	}

	raecKey := os.Getenv(raecKey)
	if raecKey != "" {
		cfg.RaecKey = raecKey
	}

	raecValue := os.Getenv(raecValue)
	if raecKey != "" {
		cfg.RaecValue = raecValue
	}

	return cfg, nil
}
