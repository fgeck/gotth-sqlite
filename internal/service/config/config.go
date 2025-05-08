package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Host          string `mapstructure:"host"`
	Port          string `mapstructure:"port"`
	JwtSecret     string `mapstructure:"jwtSecret"`
	AdminUser     string `mapstructure:"adminUser"`
	AdminPassword string `mapstructure:"adminPassword"`
	AdminEmail    string `mapstructure:"adminEmail"`
}

type DbConfig struct {
	DataBasePath   string `mapstructure:"databasePath"`
	MigrationsPath string `mapstructure:"migrationsPath"`
}

type Config struct {
	App AppConfig `mapstructure:"app"`
	Db  DbConfig  `mapstructure:"db"`
}

type ConfigLoaderInterface interface {
	LoadConfig() (*Config, error)
}

type ConfigLoader struct {
	viper *viper.Viper
}

func NewLoader() *ConfigLoader {
	return &ConfigLoader{
		viper: viper.New(),
	}
}

func (c *ConfigLoader) LoadConfig(cfgDirPath string) (*Config, error) {
	c.viper.SetConfigName("config.yaml")
	c.viper.SetConfigType("yaml")
	c.viper.AddConfigPath(".")
	c.viper.AddConfigPath("./cmd/web")
	c.viper.AddConfigPath(cfgDirPath)

	// Enable automatic environment variable binding
	c.viper.AutomaticEnv()

	// Replace `.` in environment variable keys with `_` to match YAML structure
	c.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read the config file
	if err := c.viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %v", err)

		return nil, err
	}

	var config Config
	if err := c.viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
