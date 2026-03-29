package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var logOneline bool

var logCmd = &cobra.Command{
	Use:   "log <object-id>",
	Short: "Show git commit history for an object",
	Long: `Display the git commit history for a specific object file.

Wraps git log --follow to track the object's history across renames.
Supports prefix matching — you can omit the ULID suffix if the prefix
uniquely identifies an object.

Examples:
  tmd log book/clean-code
  tmd log --oneline book/clean-code`,
	Args: cobra.ExactArgs(1),
	ValidArgsFunction: func(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return completeObjectID(toComplete)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		vault, err := openVault(vaultPath)
		if err != nil {
			return err
		}
		defer vault.Close()

		obj, err := resolveObjectInteractive(vault, args[0])
		if err != nil {
			return err
		}

		filePath := vault.ObjectPath(obj.Type, obj.Filename)

		gitArgs := []string{"-C", vault.Root, "log", "--follow"}
		if logOneline {
			gitArgs = append(gitArgs, "--oneline")
		}
		gitArgs = append(gitArgs, "--", filePath)

		gitCmd := exec.Command("git", gitArgs...)
		var stdout, stderr bytes.Buffer
		gitCmd.Stdout = &stdout
		gitCmd.Stderr = &stderr

		if err := gitCmd.Run(); err != nil {
			errMsg := strings.TrimSpace(stderr.String())
			if strings.Contains(errMsg, "does not have any commits yet") {
				fmt.Println("no commits found for this object")
				return nil
			}
			if strings.Contains(errMsg, "not a git repository") {
				return fmt.Errorf("vault is not inside a git repository")
			}
			if errMsg != "" {
				return fmt.Errorf("git log: %s", errMsg)
			}
			return fmt.Errorf("git log: %w", err)
		}

		output := strings.TrimSpace(stdout.String())
		if output == "" {
			fmt.Println("no commits found for this object")
			return nil
		}

		fmt.Println(output)
		return nil
	},
}

func init() {
	logCmd.Flags().BoolVar(&logOneline, "oneline", false, "show compact one-line output")
	rootCmd.AddCommand(logCmd)
}
