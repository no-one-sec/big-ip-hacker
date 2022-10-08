package cmd

import (
	"github.com/CC11001100/big-ip-hacker/pkg/big_ip"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	targetUrl         string
	bigIPCookieString string
)

func init() {

	cookieCmd.Flags().StringVarP(&targetUrl, "url", "u", "", "The URL was used the BIG IP")
	cookieCmd.Flags().StringVarP(&bigIPCookieString, "cookie", "c", "", "Cookies of Big IP")

	rootCmd.AddCommand(cookieCmd)
}

var cookieCmd = &cobra.Command{
	Use:   "cookie",
	Short: "Big ip cookie safe",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if bigIPCookieString != "" {
			big_ip.FromCookie(bigIPCookieString)
		} else if targetUrl != "" {
			big_ip.FromUrl(targetUrl)
		} else {
			color.Red("please use --url or --cookie specify the parameters")
			_ = cmd.Help()
		}
	},
}
