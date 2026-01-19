package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for wt.

To load completions:

Bash:
  $ source <(wt completion bash)
  # To load completions for each session, execute once:
  # Linux:
  $ wt completion bash > /etc/bash_completion.d/wt
  # macOS:
  $ wt completion bash > $(brew --prefix)/etc/bash_completion.d/wt

Zsh:
  $ source <(wt completion zsh)
  # To load completions for each session, execute once:
  $ wt completion zsh > "${fpath[1]}/_wt"

Fish:
  $ wt completion fish | source
  # To load completions for each session, execute once:
  $ wt completion fish > ~/.config/fish/completions/wt.fish

PowerShell:
  PS> wt completion powershell | Out-String | Invoke-Expression
  # To load completions for each session, add to profile:
  PS> wt completion powershell >> $PROFILE
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		}
		return nil
	},
}
