package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/safing/portbase/updater"
	"github.com/safing/portbase/utils"
	"github.com/spf13/cobra"
)

var (
	registry *updater.ResourceRegistry
	distDir  string
)

var rootCmd = &cobra.Command{
	Use:   "updatemgr",
	Short: "A simple tool to assist in the update and release process",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Check if the distribution directory exists.
		absDistPath, err := filepath.Abs(distDir)
		if err != nil {
			return fmt.Errorf("failed to get absolute path of distribution directory: %w", err)
		}
		_, err = os.Stat(absDistPath)
		if err != nil {
			return fmt.Errorf("failed to access distribution directory: %w", err)
		}

		registry = &updater.ResourceRegistry{}
		err = registry.Initialize(utils.NewDirStructure(absDistPath, 0o755))
		if err != nil {
			return err
		}

		registry.AddIndex(updater.Index{
			Path:   "stable.json",
			Stable: true,
			Beta:   false,
		})

		registry.AddIndex(updater.Index{
			Path:   "beta.json",
			Stable: false,
			Beta:   true,
		})

		err = registry.LoadIndexes(context.TODO())
		if err != nil {
			return err
		}

		err = registry.ScanStorage("")
		if err != nil {
			return err
		}

		return nil
	},
	SilenceUsage: true,
}

func init() {
	flags := rootCmd.PersistentFlags()
	flags.StringVar(&distDir, "dist-dir", "dist", "Set the distribution directory. Falls back to ./dist if available.")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
