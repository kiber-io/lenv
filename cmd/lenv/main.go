package main

import (
	"kiber-io/lenv/cmd/languages/java"
	"kiber-io/lenv/cmd/languages/python"
	"kiber-io/lenv/common"

	"github.com/spf13/cobra"
)

var version = "0.2.0"

func main() {
	var rootCmd = &cobra.Command{
		Use: "lenv",
	}
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			println("lenv", version)
		},
	}
	var printRootCmd = &cobra.Command{
		Use: "root",
		Run: func(cmd *cobra.Command, args []string) {
			println(common.GetRoot())
		},
	}
	var javaCmd = &cobra.Command{
		Use:     "java",
		Aliases: []string{"j"},
		Short:   "Manage Java versions",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				_ = cmd.Help()
				return
			}
		},
	}
	var pythonCmd = &cobra.Command{
		Use:     "python",
		Aliases: []string{"p", "py"},
		Short:   "Manage Python versions",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				_ = cmd.Help()
				return
			}
		},
	}
	rootCmd.AddCommand(printRootCmd)
	rootCmd.AddCommand(javaCmd)
	rootCmd.AddCommand(pythonCmd)
	rootCmd.AddCommand(versionCmd)
	java.Init(javaCmd)
	python.Init(pythonCmd)
	rootCmd.Execute()
}
