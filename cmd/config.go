package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/typemd/typemd/core"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage vault configuration",
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a config value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		vault, err := openVault(vaultPath, reindex)
		if err != nil {
			return err
		}
		defer vault.Close()

		return vault.SetConfigValue(args[0], args[1])
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a config value",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vault, err := openVault(vaultPath, reindex)
		if err != nil {
			return err
		}
		defer vault.Close()

		val, known := vault.GetConfigValue(args[0])
		if !known {
			return fmt.Errorf("unknown config key %q. Known keys: %s", args[0], joinKeys())
		}
		if val != "" {
			fmt.Println(val)
		}
		return nil
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all config values",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		vault, err := openVault(vaultPath, reindex)
		if err != nil {
			return err
		}
		defer vault.Close()

		for _, key := range core.ConfigKeys() {
			val, _ := vault.GetConfigValue(key)
			if val != "" {
				fmt.Printf("%s: %s\n", key, val)
			}
		}
		return nil
	},
}

func joinKeys() string {
	keys := core.ConfigKeys()
	result := ""
	for i, k := range keys {
		if i > 0 {
			result += ", "
		}
		result += k
	}
	return result
}

func init() {
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configListCmd)
	rootCmd.AddCommand(configCmd)
}
