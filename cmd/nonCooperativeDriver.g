package cmd

import (
    // "fmt"
    // "github.com/spf13/cobra"
    "github.com/seveirbian/gear/core/nonCooperativeDriver"
    // "github.com/seveirbian/gear/pkg/image"
    // "github.com/seveirbian/gear/pkg/gear"
    "github.com/docker/go-plugins-helpers/graphdriver"
)

var noCoopDriverUsage = `Usage:  gear non-cooperative-driver
Central-fs will create a graphdriver that does not cooperate with other drivers
`

func init() {
    // rootCmd.AddCommand(noCoopDriverCmd)
    noCoopDriverCmd.SetUsageTemplate(noCoopDriverUsage)
}   

var noCoopDriverCmd = &cobra.Command{
    Use:   "non-cooperative-driver",
    Short: "create a gear graphdriver that does not cooperate with other drivers",
    Long:  `create a gear graphdriver that does not cooperate with other drivers`,
    Args:  cobra.NoArgs,
    Run: func(cmd *cobra.Command, args []string) {
        noCoopDirver := &nonCooperativeDriver.NonCooperativeDriver{}

        h := graphdriver.NewHandler(noCoopDirver)

        h.ServeUnix("myGraphDriver", 0)
    },
}