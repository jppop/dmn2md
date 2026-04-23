package main

import (
	"fmt"
	"os"

	"github.com/jppop/dmn2md/internal/parser"
	"github.com/jppop/dmn2md/internal/renderer"
	"github.com/spf13/cobra"
)

var mode string

var rootCmd = &cobra.Command{
	Use:   "dmn2md <file.dmn>",
	Short: "Convert a DMN decision table to Markdown",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		defs, err := parser.ParseFile(args[0])
		if err != nil {
			return err
		}
		fmt.Print(renderer.Render(defs, renderer.Mode(mode)))
		return nil
	},
}

func init() {
	rootCmd.Flags().StringVarP(&mode, "mode", "m", string(renderer.ModeDoc),
		`output format: "doc" (human-readable) or "skill" (Claude Code slash command)`)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
