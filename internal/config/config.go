package config

import (
	"fmt"
	"os"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	API      APIConfig      `yaml:"api"`
	Admin    AdminConfig    `yaml:"admin"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

type APIConfig struct {
	AuthToken string `yaml:"auth_token"`
}

type AdminConfig struct {
	Password string `yaml:"password"`
}

func LoadConfig(configPath string) (*Config, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 使用环境变量覆盖配置
	overrideWithEnvVars(&config)

	// 验证必需的配置
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

// overrideWithEnvVars 使用环境变量覆盖配置值
func overrideWithEnvVars(config *Config) {
	// 服务器配置
	if port := os.Getenv("MAILCAT_SERVER_PORT"); port != "" {
		config.Server.Port = port
	}
	if host := os.Getenv("MAILCAT_SERVER_HOST"); host != "" {
		config.Server.Host = host
	}

	// 数据库配置
	if dbPath := os.Getenv("MAILCAT_DATABASE_PATH"); dbPath != "" {
		config.Database.Path = dbPath
	}

	// API配置 - 必须通过环境变量设置
	if authToken := os.Getenv("MAILCAT_API_AUTH_TOKEN"); authToken != "" {
		config.API.AuthToken = authToken
	}

	// 管理员配置 - 必须通过环境变量设置
	if adminPassword := os.Getenv("MAILCAT_ADMIN_PASSWORD"); adminPassword != "" {
		config.Admin.Password = adminPassword
	}
}

// validateConfig 验证配置的必需字段
func validateConfig(config *Config) error {
	if config.API.AuthToken == "" {
		return fmt.Errorf("API auth token is required. Please set MAILCAT_API_AUTH_TOKEN environment variable")
	}
	if config.Admin.Password == "" {
		return fmt.Errorf("Admin password is required. Please set MAILCAT_ADMIN_PASSWORD environment variable")
	}
	return nil
}