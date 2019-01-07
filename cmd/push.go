package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    // "github.com/seveirbian/gear/pkg/image"
    // "github.com/seveirbian/gear/pkg/gear"
)

var pushUsage = `Usage:  gear push -i xxx.xxx.xxx.xxx NAME:TAG
Options:
  -i,  --ip              seaweedfs ip address
`

var ipAddress string

func init() {
    rootCmd.AddCommand(pushCmd)
    pushCmd.SetUsageTemplate(pushUsage)
    pushCmd.Flags().StringVarP(&ipAddress, "ip", "i", "", "seaweedfs ip address")
    pushCmd.MarkFlagRequired("ip")
}

var pushCmd = &cobra.Command{
    Use:   "push",
    Short: "push a gear image to seaweedfs",
    Long:  `push a gear image to seaweedfs`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("push")
    },
}