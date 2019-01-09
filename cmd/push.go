package cmd

import (
    // "fmt"
    "github.com/spf13/cobra"
    "github.com/seveirbian/gear/pkg/image"
    "github.com/seveirbian/gear/core/push"
)

var pushUsage = `Usage:  gear push --push-url http://xxx.xxx.xxx.xxx:xxx/... NAME:TAG
Options:
  --push-url              server url address
`

var pushURL string

func init() {
    rootCmd.AddCommand(pushCmd)
    pushCmd.SetUsageTemplate(pushUsage)
    pushCmd.Flags().StringVarP(&pushURL, "push-url", "", "", "server url address")
    pushCmd.MarkFlagRequired("push-url")
}

var pushCmd = &cobra.Command{
    Use:   "push",
    Short: "push a gear image to server",
    Long:  `push a gear image to server`,
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        image := image.ParseGearImage(args[0])

        gearImagePusher := push.InitPusher(image, pushURL)

        gearImagePusher.Push()
    },
}