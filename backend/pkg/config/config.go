package config

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

// Config :
type Config struct {
	Matsuno AppConfig `mapstructure:"matsuno"`
	Prd     AppConfig `mapstructure:"prd"`
	Stg     AppConfig `mapstructure:"stg"`
	Test    AppConfig `mapstructure:"test"`
}

// HTTP :
type HTTP struct {
	Cors string `mapstructure:"cors" validate:"required"`
	Port int    `mapstructure:"port" validate:"required"`
}

// Logger :
type Logger struct {
	Debug   bool `mapstructure:"debug"`
	LogJSON bool `mapstructure:"log_json"`
}

// Postgres :
type Postgres struct {
	DBName  string `mapstructure:"dbname" validate:"required"`
	Host    string `mapstructure:"host" validate:"required"`
	Pass    string `mapstructure:"pass" validate:"required"`
	Port    string `mapstructure:"port" validate:"required"`
	Sslmode string `mapstructure:"sslmode" validate:"required"`
	User    string `mapstructure:"user" validate:"required"`
	Pseudo  bool
}

// Firebase :
type Firebase struct {
	CredentialKey string `mapstructure:"credential_key" validate:"required"`
	Pseudo        bool
}

// AppConfig :
type AppConfig struct {
	HTTP     HTTP     `mapstructure:"http"`
	Logger   Logger   `mapstructure:"logger"`
	Postgres Postgres `mapstructure:"postgres"`
	Firebase Firebase `mapstructure:"firebase"`
}

// Prepare :
func Prepare() AppConfig {
	viper.SetConfigName("config")
	viper.SetEnvPrefix("WDC")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigType("yaml")

	_, b, _, _ := runtime.Caller(0)
	configDir := filepath.Dir(b)
	pkgDir := filepath.Dir(configDir)
	backendDir := filepath.Dir(pkgDir)

	viper.AddConfigPath(backendDir)
	viper.AddConfigPath("./")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	viper.AutomaticEnv()

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		panic(err)
	}

	var appConfig AppConfig
	env := viper.GetString("env.name")
	if env == "test" {
		appConfig = c.Test
	} else {
		panic(fmt.Sprintf("Unknown env: %s", env))
	}
	return appConfig
}
