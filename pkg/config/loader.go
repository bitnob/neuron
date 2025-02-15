package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Loader handles configuration loading from different sources
type Loader struct {
	sources []Source
}

// Source represents a configuration source
type Source interface {
	Load() (*Config, error)
	Priority() int
}

// FileSource represents a file-based configuration source
type FileSource struct {
	path     string
	format   string
	priority int
}

// EnvSource represents environment variable based configuration
type EnvSource struct {
	prefix   string
	priority int
}

// NewLoader creates a new configuration loader
func NewLoader(sources ...Source) *Loader {
	return &Loader{
		sources: sources,
	}
}

// LoadConfig loads configuration from all sources
func (l *Loader) LoadConfig() (*Config, error) {
	var config Config

	// Sort sources by priority
	sortSources(l.sources)

	// Load from each source
	for _, source := range l.sources {
		sourceConfig, err := source.Load()
		if err != nil {
			return nil, fmt.Errorf("error loading config from source: %w", err)
		}

		// Merge configurations
		if err := mergeConfig(&config, sourceConfig); err != nil {
			return nil, fmt.Errorf("error merging config: %w", err)
		}
	}

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation error: %w", err)
	}

	return &config, nil
}

// NewFileSource creates a new file-based configuration source
func NewFileSource(path string, priority int) *FileSource {
	return &FileSource{
		path:     path,
		format:   strings.ToLower(filepath.Ext(path)),
		priority: priority,
	}
}

func (fs *FileSource) Load() (*Config, error) {
	data, err := os.ReadFile(fs.path)
	if err != nil {
		return nil, err
	}

	var config Config
	switch fs.format {
	case ".json":
		err = json.Unmarshal(data, &config)
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, &config)
	default:
		return nil, fmt.Errorf("unsupported file format: %s", fs.format)
	}

	return &config, err
}

func (fs *FileSource) Priority() int {
	return fs.priority
}
