package cmd

import (
	"fmt"
	"github.com/CC11001100/go-StringBuilder/pkg/string_builder"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(aboutCmd)
}

var aboutCmd = &cobra.Command{
	Use:   "about",
	Short: "About this tool github, about author, blabla.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		msg := string_builder.New().
			AppendString("\n\n\n").
			AppendString(fmt.Sprintf("%20s", "Repo")).AppendString(" : ").AppendString("https://github.com/CC11001100/big-ip-hacker\n\n").
			AppendString(fmt.Sprintf("%20s", "Problem feedback")).AppendString(" : ").AppendString("https://github.com/CC11001100/big-ip-hacker/issues\n\n").
			AppendString(fmt.Sprintf("%20s", "Author")).AppendString(" : ").AppendString("CC11001100\n").
			AppendString("\n\n\n").
			String()
		color.HiGreen(msg)
	},
}
