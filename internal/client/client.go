package client

import (
	"fmt"

	whooktown "github.com/fredericalix/whooktown-golang-sdk"
	"github.com/fredericalix/whooktown-cli/internal/config"
)

// Options configures the SDK client
type Options struct {
	Token      string // Override token from flag
	ConfigPath string
	Config     *config.Config // Pre-loaded config (optional)
}

// New creates a configured SDK client
func New(opts Options) (*whooktown.Client, error) {
	var cfg *config.Config
	var err error

	if opts.Config != nil {
		cfg = opts.Config
	} else {
		cfg, err = config.Load(opts.ConfigPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load config: %w", err)
		}
	}

	ctx := cfg.CurrentCtx()

	// Token priority: flag > config > error
	token := opts.Token
	if token == "" {
		token = ctx.Token
	}
	if token == "" {
		return nil, fmt.Errorf("not logged in. Run 'wt login' or use --token flag")
	}

	// Build SDK options
	sdkOpts := []whooktown.Option{
		whooktown.WithToken(token),
	}

	// Set environment
	if ctx.Environment == "DEV" {
		sdkOpts = append(sdkOpts, whooktown.WithEnvironment(whooktown.EnvDevelopment))
	} else {
		sdkOpts = append(sdkOpts, whooktown.WithEnvironment(whooktown.EnvProduction))
	}

	// Apply URL overrides if set
	if ctx.AuthURL != "" {
		sdkOpts = append(sdkOpts, whooktown.WithAuthURL(ctx.AuthURL))
	}
	if ctx.UIURL != "" {
		sdkOpts = append(sdkOpts, whooktown.WithUIURL(ctx.UIURL))
	}
	if ctx.SensorURL != "" {
		sdkOpts = append(sdkOpts, whooktown.WithSensorURL(ctx.SensorURL))
	}
	if ctx.WorkflowURL != "" {
		sdkOpts = append(sdkOpts, whooktown.WithWorkflowURL(ctx.WorkflowURL))
	}

	return whooktown.New(sdkOpts...)
}

// NewUnauthenticated creates a client without token (for login flow)
func NewUnauthenticated(cfg *config.Config) (*whooktown.Client, error) {
	ctx := cfg.CurrentCtx()

	sdkOpts := []whooktown.Option{}
	if ctx.Environment == "DEV" {
		sdkOpts = append(sdkOpts, whooktown.WithEnvironment(whooktown.EnvDevelopment))
	} else {
		sdkOpts = append(sdkOpts, whooktown.WithEnvironment(whooktown.EnvProduction))
	}

	// Apply URL overrides if set
	if ctx.AuthURL != "" {
		sdkOpts = append(sdkOpts, whooktown.WithAuthURL(ctx.AuthURL))
	}

	return whooktown.New(sdkOpts...)
}
