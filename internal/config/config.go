package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	configPathEnvKey   = "CONFIG_PATH"
	DefaultFileMaxSize = 10 << 20
	max                = 100 << 20
)

type Config struct {
	HTTPServer 					`yaml:"http_server"										env-required:"true"`
	Database   StorageConfig 	`yaml:"database" 										env-required:"true"`
	Redis      RedisConfig   	`yaml:redis												env-required:"true"`
}

type HTTPServer struct {
	Address     string        	`yaml:"address"			env-default:"localhost:8080`
	Timeout     time.Duration 	`yaml:"timeout" 											env-required:"true"`
	IdleTimeout time.Duration 	`yaml:"iddle_timeout" 										env-required:"true"`
	TLS         TLSConfig     	`yaml:"tls"													env-required:"true"`
	MaxFileSize int64         	`yaml:"max_file_size"`
	Cors        CORSConfig    	`yaml:"cors"`
}

type StorageConfig struct {
	Host         string 		`yaml:"host" 			env-default:"localhost"	`
	Port         string 		`yaml:"port" 			env-default:"5432" 		`
	DatabaseName string 		`yaml:"databaseName" 	env-default:"postgres" 	`
	User         string 		`yaml:"user" 			env-default:"postgres" 	`
	Password     string 		`yaml:"password" 		env-default:"1488" 		`
}

type RedisConfig struct {
	Host     string 			`yaml:"host" 			env-default:"localhost"	`
	Port     string 			`yaml:"port" 			env-default:"5432" 		`
	User     string 			`yaml:"user" 			env-default:"default" 	`
	Password string 			`yaml:"password" 		env-default:"" 			`
}

type TLSConfig struct {
	CertFile string 			`yaml:"certfile"`
	KeyFile  string 			`yaml:"keyfile"	`
}

type CORSConfig struct {
	AllowedOrigins []string 	`yaml:"allowed_origins"`
	AllowedMethods []string 	`yaml:"allowed_methods"`
	AllowedHeaders []string 	`yaml:"allowed_headers"`
}

func MustLoadConfig() *Config {
	configPath := os.Getenv(configPathEnvKey)
	if configPath == "" {
		log.Fatalf("%s is not set up", configPathEnvKey)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist: %s", configPath, err.Error())
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to read config file: %s", err.Error())
	}

	if cfg.MaxFileSize == 0 {
		cfg.MaxFileSize = DefaultFileMaxSize
	}

	return &cfg
}
