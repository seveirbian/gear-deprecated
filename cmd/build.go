package cmd

import (
    "github.com/spf13/cobra"
    "github.com/seveirbian/gear/build"
    "github.com/seveirbian/gear/pkg/image"
)

var dockerImage string

func init() {
    rootCmd.AddCommand(buildCmd)
    buildCmd.Flags().StringVarP(&dockerImage, "docker-image", "", "", "build a gear image from a standard docker image")
}   

var buildCmd = &cobra.Command{
    Use:   "build",
    Short: "build a gear image from a standard docker image",
    Long:  `build a gear image from a standard docker image`,
    Run: func(cmd *cobra.Command, args []string) {
        image := image.Parse(dockerImage)
        gearImageBuilder := build.InitBuilder(image)

        gearImageBuilder.Build()
    },
}