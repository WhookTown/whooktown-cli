package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	whooktown "github.com/fredericalix/whooktown-golang-sdk"
	"github.com/fredericalix/whooktown-cli/internal/client"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to whooktown",
	Long: `Login to whooktown using email-based authentication.

An email with a validation link will be sent to your address.
Click the link to complete authentication.`,
	RunE: runLogin,
}

func runLogin(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Get email from user
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	if email == "" {
		return fmt.Errorf("email is required")
	}

	// Create unauthenticated client
	c, err := client.NewUnauthenticated(cfg)
	if err != nil {
		return err
	}

	// Initiate login
	fmt.Printf("Sending verification email to %s...\n", email)
	token, err := c.Auth.Login(ctx, &whooktown.LoginRequest{
		Email: email,
		Type:  "user",
		Name:  "wt-cli",
		AppID: "cli",
	})
	if err != nil {
		// Check if account doesn't exist - try signup
		fmt.Println("Account not found. Creating new account...")
		token, err = c.Auth.Signup(ctx, &whooktown.SignupRequest{
			Email: email,
			Type:  "user",
			Name:  "wt-cli",
			AppID: "cli",
		})
		if err != nil {
			return fmt.Errorf("signup failed: %w", err)
		}
	}

	fmt.Println("\nCheck your email and click the validation link.")
	fmt.Println("Waiting for validation...")

	// Poll for validation (token becomes valid after email click)
	validated, err := pollForValidation(ctx, c, token.Token, 5*time.Minute)
	if err != nil {
		return err
	}

	// Save token to config
	cfg.SetToken(validated.Token)
	if err := cfg.Save(cfgFile); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	formatter.Success(fmt.Sprintf("Logged in as %s", email))
	return nil
}

// pollForValidation waits for the token to be validated
func pollForValidation(ctx context.Context, c *whooktown.Client, token string, timeout time.Duration) (*whooktown.Token, error) {
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			if time.Now().After(deadline) {
				return nil, fmt.Errorf("validation timeout - please try again")
			}

			// Check if token is now valid
			info, err := c.Auth.CheckToken(ctx, token)
			if err != nil {
				// Error means still pending validation
				fmt.Print(".")
				continue
			}

			// Token is validated when it has an AccountID and Account is validated
			if info.Account != nil && info.Account.Validated {
				fmt.Println() // New line after dots
				return info, nil
			}
			fmt.Print(".")
		}
	}
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout and remove stored credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		// Try to call logout endpoint if we have a token
		if cfg.CurrentCtx().Token != "" {
			c, err := client.New(client.Options{
				Token:  tokenFlag,
				Config: cfg,
			})
			if err == nil {
				_ = c.Auth.Logout(ctx, "cli") // Ignore errors
			}
		}

		// Clear token from config
		cfg.SetToken("")
		if err := cfg.Save(cfgFile); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		formatter.Success("Logged out successfully")
		return nil
	},
}

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current account information",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		c, err := client.New(client.Options{
			Token:  tokenFlag,
			Config: cfg,
		})
		if err != nil {
			return err
		}

		// Get token info
		token := tokenFlag
		if token == "" {
			token = cfg.CurrentCtx().Token
		}

		info, err := c.Auth.CheckToken(ctx, token)
		if err != nil {
			return fmt.Errorf("failed to get account info: %w", err)
		}

		if jsonOutput {
			return formatter.Format(os.Stdout, info)
		}

		email := ""
		if info.Account != nil {
			email = info.Account.Email
		}
		fmt.Printf("Email:       %s\n", email)
		fmt.Printf("Account ID:  %s\n", info.AccountID)
		fmt.Printf("Token Type:  %s\n", info.Type)
		fmt.Printf("Context:     %s\n", cfg.CurrentContext)
		fmt.Printf("Environment: %s\n", cfg.CurrentCtx().Environment)

		return nil
	},
}
