/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"strings"

	"github.com/codeready-toolchain/argocd-checker/pkg/validation"

	charmlog "github.com/charmbracelet/log"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := checkCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var apps, components []string
var baseDir string
var verbose bool

// checkCmd represents the base command when called without any subcommands
var checkCmd = &cobra.Command{
	Use:   "check-argocd",
	Short: "Checks the Argo CD configuration",

	Run: func(cmd *cobra.Command, args []string) {

		logger := charmlog.New(cmd.OutOrStderr())
		logger.SetLevel(charmlog.InfoLevel)
		if verbose {
			logger.SetLevel(charmlog.DebugLevel)
		}

		afs := afero.Afero{
			Fs: afero.NewOsFs(),
		}

		// verifies that the source path of the Applications and ApplicationSets exists
		if err := validation.CheckApplications(logger, afs, baseDir, apps...); err != nil {
			logger.Error(strings.ReplaceAll(err.Error(), ": ", ":\n"))
			os.Exit(1)
		}
		// verifies that `kustomize build` on each component completes successfully
		if err := validation.CheckComponents(logger, afs, baseDir, components...); err != nil {
			logger.Error(strings.ReplaceAll(err.Error(), ": ", ":\n"))
			os.Exit(1)
		}
	},
}

func init() {
	checkCmd.Flags().StringSliceVar(&apps, "apps", []string{}, "path(s) to the applications (comma-separated, relative to '--baseDir')")
	// if err := checkCmd.MarkFlagRequired("apps"); err != nil {
	// 	panic(fmt.Sprintf("failed to mark flag as required: %s", err))
	// }
	checkCmd.Flags().StringVar(&baseDir, "base-dir", ".", "base directory of the repository")
	checkCmd.Flags().StringSliceVar(&components, "components", []string{}, "path(s) to the components (comma-separated, relative to '--baseDir')")
	// if err := checkCmd.MarkFlagRequired("components"); err != nil {
	// 	panic(fmt.Sprintf("failed to mark flag as required: %s", err))
	// }
	checkCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}
