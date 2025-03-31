package config

import (
	"github.com/spf13/viper"
	"strings"
	"time"
)

type Config struct {
	Log      Log
	Server   Server
	Database Database
}

func LoadConfig(filePath string) (Config, error) {
	viper.SetEnvPrefix("FLAGON")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()
	if filePath != "" {
		viper.SetConfigFile(filePath)
		if err := viper.ReadInConfig(); err != nil {
			return Config{}, err
		}
	}
	for _, key := range viper.AllKeys() {
		viper.SetDefault(key, viper.Get(key))
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, err
	}
	return config, nil
}

func init() {
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "text")
	viper.SetDefault("log.file", "stderr")
	viper.SetDefault("log.addSource", true)

	viper.SetDefault("http.host", "0.0.0.0")
	viper.SetDefault("http.port", 8080)
	viper.SetDefault("http.enableTLS", false)
	viper.SetDefault("http.certFile", "")
	viper.SetDefault("http.keyFile", "")

	viper.SetDefault("db.driver", "sqlite")
	viper.SetDefault("db.host", "")
	viper.SetDefault("db.port", "")
	viper.SetDefault("db.username", "")
	viper.SetDefault("db.password", "")
	viper.SetDefault("db.database", "flagon.sqlite")
	viper.SetDefault("db.sslmode", "")
	viper.SetDefault("db.maxConnLifetime", 0)
	viper.SetDefault("db.maxConnIdleTime", 0)
	viper.SetDefault("db.maxOpenConns", 0)
	viper.SetDefault("db.maxIdleConns", 0)
}

type Log struct {
	Level     string
	Format    string
	File      string
	AddSource bool
}
type Server struct {
	Host      string
	Port      int
	EnableTLS bool
	CertFile  string
	KeyFile   string
}

type Database struct {
	Driver          string
	Host            string
	Port            uint16
	Username        string
	Password        string
	Database        string
	SSLMode         string
	MaxConns        int
	MaxIdleConns    int
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}
