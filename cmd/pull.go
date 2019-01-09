package cmd

import (
    // "fmt"
    "github.com/spf13/cobra"
    // "github.com/seveirbian/gear/pkg/image"
    "github.com/seveirbian/gear/core/pull"
)

var pullUsage = `Usage:  gear pull --pull-url http://xxx.xxx.xxx.xxx:xxx/... FILENAME
Options:
  --pull-url              server url address
`

var pullURL string

func init() {
    rootCmd.AddCommand(pullCmd)
    pullCmd.SetUsageTemplate(pullUsage)
    pullCmd.Flags().StringVarP(&pullURL, "pull-url", "", "", "server url address")
    pullCmd.MarkFlagRequired("pull-url")
}

var pullCmd = &cobra.Command{
    Use:   "pull",
    Short: "pull a file from server",
    Long:  `pull a file from server`,
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {

        gearFilePuller := pull.InitPuller(args[0], pullURL)

        gearFilePuller.Pull()
    },
}