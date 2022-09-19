package cmd

import (
	"context"
	"fmt"
	"github.com/elimity-com/abnf/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"path/filepath"
)

const Version = "v0.1"

func Do(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	rootCmd := &cobra.Command{Use: "abnf", SilenceUsage: true}
	rootCmd.PersistentFlags().StringP("file", "f", "", "specify an alternate config file (default: abnf.yml)")

	rootCmd.AddCommand(genCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(versionCmd)

	rootCmd.SetArgs(args)
	rootCmd.SetIn(stdin)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)

	ctx := context.Background()

	if err := rootCmd.ExecuteContext(ctx); err == nil {
		return 0
	}
	return 1
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the abnf version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", Version)
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create an empty abnf.yml settings file",
	RunE: func(cmd *cobra.Command, args []string) error {
		file := "abnf.yml"
		if f := cmd.Flag("file"); f != nil && f.Changed {
			file = f.Value.String()
			if file == "" {
				return fmt.Errorf("file argument is empty")
			}
		}
		if _, err := os.Stat(file); !os.IsNotExist(err) {
			return nil
		}
		blob, err := yaml.Marshal(config.Config{Version: "1", VerboseFlag: false})
		if err != nil {
			return err
		}
		return os.WriteFile(file, blob, 0644)
	},
}

var genCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Go code from an ABNF specification file",
	Run: func(cmd *cobra.Command, args []string) {
		stderr := cmd.ErrOrStderr()
		dir, name := getConfigPath(stderr, cmd.Flag("file"))

		err := Generate(dir, name, stderr)
		if err != nil {
			os.Exit(1)
		}
	},
}

func getConfigPath(stderr io.Writer, f *pflag.Flag) (string, string) {
	if f != nil && f.Changed {
		file := f.Value.String()
		if file == "" {
			fmt.Fprintln(stderr, "error parsing config: file argument is empty")
			os.Exit(1)
		}
		abspath, err := filepath.Abs(file)
		if err != nil {
			fmt.Fprintf(stderr, "error parsing config: absolute file path lookup failed: %s\n", err)
			os.Exit(1)
		}
		return filepath.Dir(abspath), filepath.Base(abspath)
	} else {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(stderr, "error parsing abnf.json: file does not exist")
			os.Exit(1)
		}
		return wd, ""
	}
}
