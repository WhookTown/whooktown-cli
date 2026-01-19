package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the CLI configuration
type Config struct {
	CurrentContext string              `yaml:"current_context"`
	Contexts       map[string]*Context `yaml:"contexts"`
}

// Context represents a named environment configuration
type Context struct {
	Name          string `yaml:"name"`
	Token         string `yaml:"token,omitempty"`
	Environment   string `yaml:"environment"` // "PROD" or "DEV"
	DefaultLayout string `yaml:"default_layout,omitempty"`

	// Optional URL overrides (for custom deployments)
	AuthURL     string `yaml:"auth_url,omitempty"`
	UIURL       string `yaml:"ui_url,omitempty"`
	SensorURL   string `yaml:"sensor_url,omitempty"`
	WorkflowURL string `yaml:"workflow_url,omitempty"`
}

// DefaultConfigPath returns ~/.whooktown/config.yaml
func DefaultConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".whooktown", "config.yaml")
}

// Load reads the config file or returns defaults
func Load(path string) (*Config, error) {
	cfg := &Config{
		CurrentContext: "default",
		Contexts: map[string]*Context{
			"default": {
				Name:        "default",
				Environment: "PROD",
			},
		},
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil // Return defaults
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Save writes config to disk
func (c *Config) Save(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// CurrentCtx returns the active context
func (c *Config) CurrentCtx() *Context {
	if ctx, ok := c.Contexts[c.CurrentContext]; ok {
		return ctx
	}
	return c.Contexts["default"]
}

// SetToken updates the token for the current context
func (c *Config) SetToken(token string) {
	c.CurrentCtx().Token = token
}

// GetToken returns the token for the current context
func (c *Config) GetToken() string {
	return c.CurrentCtx().Token
}
