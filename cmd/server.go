package cmd

import (
    // "fmt"
    "github.com/spf13/cobra"
    "github.com/seveirbian/gear/core/server"
)

var serverIp string
var serverPort string

var serverUsage = `Usage:  gear server --server-ip IPADDRESS --server-port PORT
`

func init() {
    rootCmd.AddCommand(serverCmd)
    serverCmd.SetUsageTemplate(serverUsage)
    pushCmd.Flags().StringVarP(&serverIp, "server-ip", "", "", "gear file server ip address")
    pushCmd.Flags().StringVarP(&serverPort, "server-port", "", "9333", "gear file server port")
}   

var serverCmd = &cobra.Command{
    Use:   "server",
    Short: "create a central file server for gear",
    Long:  `create a central file server for gear`,
    Args:  cobra.NoArgs,
    Run: func(cmd *cobra.Command, args []string) {
        s := server.InitServer(serverIp, serverPort)

        s.InitRoute()

        s.Start()
    },
}