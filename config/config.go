package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	HTTP      HTTPConfig
	Storage   StorageConfig
	LogLevel  string `mapstructure:"log_level"`
	Token     TokenConfig
	Flashcall Flashcall
	Redis     Redis
	Cookie    CookieConfig
}

type HTTPConfig struct {
	Address            string        `mapstructure:"address"`
	ReadTimeout        time.Duration `mapstructure:"readTimeout"`
	WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
	MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
}

type TokenConfig struct {
	TokenSecretKey       string        `mapstructure:"token_secret_key"`
	AccessTokenDuration  time.Duration `mapstructure:"access_token_duration"`
	RefreshTokenDuration time.Duration `mapstructure:"refresh_token_duration"`
}
type CookieConfig struct {
	AuthcookieName   string `mapstructure:"authcookie_name"`
	AuthcookiePath   string `mapstructure:"authcookie_path"`
	AuthcookieDomain string `mapstructure:"authcookie_domain"`
	AccessName       string `mapstructure:"accessname"`
	AccessTtl        int    `mapstructure:"access_ttl"`
	RefreshName      string `mapstructure:"refreshname"`
	RefreshTtl       int    `mapstructure:"refresh_ttl"`
}
type StorageConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port" `
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

type Flashcall struct {
	PublicKey  string `mapstructure:"public_key"`
	CampaignID string `mapstructure:"campaign_id"`
}
type Redis struct {
	Host      string        `mapstructure:"host"`
	Port      string        `mapstructure:"port"`
	ExpiredAt time.Duration `mapstructure:"expired_at"`
	Password  string        `mapstructure:"password"`
}

func LoadConfig(path string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
