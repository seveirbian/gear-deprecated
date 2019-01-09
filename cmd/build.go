package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/seveirbian/gear/core/build"
    "github.com/seveirbian/gear/pkg/image"
    // "github.com/seveirbian/gear/pkg/gear"
)

var buildUsage = `Usage:  gear build NAME:TAG`

func init() {
    rootCmd.AddCommand(buildCmd)
    buildCmd.SetUsageTemplate(buildUsage)
}   

var buildCmd = &cobra.Command{
    Use:   "build",
    Short: "build a gear image from a standard docker image",
    Long:  `build a gear image from a standard docker image`,
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        image := image.Parse(args[0])

        gearImageBuilder := build.InitBuilder(image)

        if gearImageBuilder.HasParsedThisImage() {
            fmt.Println("This image has been built.")
        }else {
            gearImageBuilder.Build()
        }
    },
}