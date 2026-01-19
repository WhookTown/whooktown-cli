package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/gofrs/uuid"
	"github.com/spf13/cobra"

	whooktown "github.com/fredericalix/whooktown-golang-sdk"
	"github.com/fredericalix/whooktown-cli/internal/client"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Manage authentication tokens",
}

func init() {
	tokenCmd.AddCommand(tokenListCmd)
	tokenCmd.AddCommand(tokenCreateCmd)
	tokenCmd.AddCommand(tokenRevokeCmd)
}

var tokenListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tokens for the current account",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		c, err := client.New(client.Options{
			Token:  tokenFlag,
			Config: cfg,
		})
		if err != nil {
			return err
		}

		tokens, err := c.Auth.ListTokens(ctx)
		if err != nil {
			return err
		}

		if jsonOutput {
			return formatter.Format(os.Stdout, tokens)
		}

		if len(tokens) == 0 {
			fmt.Println("No tokens found")
			return nil
		}

		headers := []string{"TOKEN", "TYPE", "NAME", "VALIDATED"}
		rows := make([][]string, len(tokens))
		for i, t := range tokens {
			tokenShort := t.Token
			if len(tokenShort) > 12 {
				tokenShort = tokenShort[:8] + "..." + tokenShort[len(tokenShort)-4:]
			}
			validated := "pending"
			if t.ValidationLink == "" && t.AccountID != uuid.Nil {
				validated = "yes"
			}
			name := t.Name
			if name == "" {
				name = "-"
			}
			rows[i] = []string{tokenShort, t.Type, name, validated}
		}

		return formatter.FormatTable(os.Stdout, headers, rows)
	},
}

var (
	tokenName string
	tokenType string
)

var tokenCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new token",
	Long: `Create a new token for the current account.

Token types:
  user    - Full user access
  sensor  - Only sensor write access
  viewer  - Read-only access`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		c, err := client.New(client.Options{
			Token:  tokenFlag,
			Config: cfg,
		})
		if err != nil {
			return err
		}

		token, err := c.Auth.CreateToken(ctx, &whooktown.CreateTokenRequest{
			Name: tokenName,
			Type: tokenType,
		})
		if err != nil {
			return err
		}

		if jsonOutput {
			return formatter.Format(os.Stdout, token)
		}

		fmt.Printf("Token created: %s\n", token.Token)
		fmt.Printf("Type: %s\n", token.Type)
		if token.Name != "" {
			fmt.Printf("Name: %s\n", token.Name)
		}
		fmt.Println("\nSave this token - it will not be shown again!")

		return nil
	},
}

func init() {
	tokenCreateCmd.Flags().StringVar(&tokenName, "name", "", "token name")
	tokenCreateCmd.Flags().StringVar(&tokenType, "type", "user", "token type (user, sensor, viewer)")
}

var tokenRevokeCmd = &cobra.Command{
	Use:   "revoke <token>",
	Short: "Revoke a token",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		c, err := client.New(client.Options{
			Token:  tokenFlag,
			Config: cfg,
		})
		if err != nil {
			return err
		}

		if err := c.Auth.RevokeToken(ctx, args[0]); err != nil {
			return err
		}

		formatter.Success(fmt.Sprintf("Token revoked: %s", args[0]))
		return nil
	},
}
