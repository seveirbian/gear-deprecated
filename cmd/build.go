package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/seveirbian/gear/build"
    "github.com/seveirbian/gear/pkg/image"
    // "github.com/seveirbian/gear/pkg/gear"
)

var buildImage string

func init() {
    rootCmd.AddCommand(buildCmd)
    buildCmd.Flags().StringVarP(&buildImage, "docker-image", "", "", "build a gear image from a standard docker image")
    rootCmd.MarkFlagRequired("docker-image")
}   

var buildCmd = &cobra.Command{
    Use:   "build",
    Short: "build a gear image from a standard docker image",
    Long:  `build a gear image from a standard docker image`,
    Run: func(cmd *cobra.Command, args []string) {
        image := image.Parse(buildImage)

        gearImageBuilder := build.InitBuilder(image)

        if gearImageBuilder.HasParsedThisImage() {
            fmt.Println("This image has been built.")
        }else {
            gearImageBuilder.Build()
        }
    },
}