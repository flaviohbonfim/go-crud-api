package config

import (
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	AppEnv            string `mapstructure:"APP_ENV"`
	HTTPPort          string `mapstructure:"HTTP_PORT"`
	JWTSecret         string `mapstructure:"JWT_SECRET"`
	AccessTokenTTL    string `mapstructure:"ACCESS_TOKEN_TTL"`
	RefreshTokenTTL   string `mapstructure:"REFRESH_TOKEN_TTL"`
	DBHost            string `mapstructure:"DB_HOST"`
	DBPort            string `mapstructure:"DB_PORT"`
	DBUser            string `mapstructure:"DB_USER"`
	DBPassword        string `mapstructure:"DB_PASSWORD"`
	DBName            string `mapstructure:"DB_NAME"`
	DBSslMode         string `mapstructure:"DB_SSLMODE"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigFile(path + "/.env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
