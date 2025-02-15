package config

import (
	"sync"
)

// Config represents the main configuration structure
type Config struct {
	App        AppConfig        `json:"app" yaml:"app"`
	Server     ServerConfig     `json:"server" yaml:"server"`
	Database   DatabaseConfig   `json:"database" yaml:"database"`
	Cache      CacheConfig      `json:"cache" yaml:"cache"`
	Security   SecurityConfig   `json:"security" yaml:"security"`
	Middleware MiddlewareConfig `json:"middleware" yaml:"middleware"`
	Logger     LoggerConfig     `json:"logger" yaml:"logger"`
}

type AppConfig struct {
	Name        string            `json:"name" yaml:"name"`
	Environment string            `json:"environment" yaml:"environment"`
	Debug       bool              `json:"debug" yaml:"debug"`
	TimeZone    string            `json:"timezone" yaml:"timezone"`
	Metadata    map[string]string `json:"metadata" yaml:"metadata"`
}

type ServerConfig struct {
	Host            string `json:"host" yaml:"host"`
	Port            int    `json:"port" yaml:"port"`
	ReadTimeout     int    `json:"readTimeout" yaml:"readTimeout"`
	WriteTimeout    int    `json:"writeTimeout" yaml:"writeTimeout"`
	MaxHeaderBytes  int    `json:"maxHeaderBytes" yaml:"maxHeaderBytes"`
	GracefulTimeout int    `json:"gracefulTimeout" yaml:"gracefulTimeout"`
}

type DatabaseConfig struct {
	Driver          string            `json:"driver" yaml:"driver"`
	Host            string            `json:"host" yaml:"host"`
	Port            int               `json:"port" yaml:"port"`
	Name            string            `json:"name" yaml:"name"`
	User            string            `json:"user" yaml:"user"`
	Password        string            `json:"password" yaml:"password"`
	MaxConnections  int               `json:"maxConnections" yaml:"maxConnections"`
	MaxIdleConns    int               `json:"maxIdleConns" yaml:"maxIdleConns"`
	ConnMaxLifetime int               `json:"connMaxLifetime" yaml:"connMaxLifetime"`
	Options         map[string]string `json:"options" yaml:"options"`
}

type CacheConfig struct {
	Driver     string `json:"driver" yaml:"driver"`
	Host       string `json:"host" yaml:"host"`
	Port       int    `json:"port" yaml:"port"`
	Password   string `json:"password" yaml:"password"`
	DB         int    `json:"db" yaml:"db"`
	MaxRetries int    `json:"maxRetries" yaml:"maxRetries"`
	PoolSize   int    `json:"poolSize" yaml:"poolSize"`
	DefaultTTL int    `json:"defaultTTL" yaml:"defaultTTL"`
}

// ConfigManager handles configuration loading and access
type ConfigManager struct {
	config *Config
	mu     sync.RWMutex
}

// Global configuration instance
var (
	globalConfig *ConfigManager
	once         sync.Once
)

// GetConfig returns the global configuration instance
func GetConfig() *ConfigManager {
	once.Do(func() {
		globalConfig = &ConfigManager{
			config: &Config{},
		}
	})
	return globalConfig
}

// Get returns the current configuration
func (cm *ConfigManager) Get() *Config {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config
}

// Set updates the current configuration
func (cm *ConfigManager) Set(config *Config) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.config = config
}
