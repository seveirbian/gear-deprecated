package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/seveirbian/gear/pkg/image"
    // "github.com/seveirbian/gear/pkg/gear"
)

var pushImage string
var ipAddress string

func init() {
    rootCmd.AddCommand(pushCmd)
    buildCmd.Flags().StringVarP(&pushImage, "docker-image", "", "", "push a gear image from a standard docker image")
    rootCmd.MarkFlagRequired("docker-image")
    buildCmd.Flags().StringVarP(&ipAddress, "ip", "", "", "seaweedfs ip address")
    rootCmd.MarkFlagRequired("ip")
}

var pushCmd = &cobra.Command{
    Use:   "push",
    Short: "push a gear image to seaweedfs",
    Long:  `push a gear image to seaweedfs`,
    Run: func(cmd *cobra.Command, args []string) {
        
    },
}