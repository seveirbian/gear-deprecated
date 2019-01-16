package cmd

import (
    // "fmt"
    "strings"

    "github.com/spf13/cobra"
    "github.com/seveirbian/gear/core/nonCooperativeFs"
    // "github.com/seveirbian/gear/pkg/image"
    // // "github.com/seveirbian/gear/pkg/gear"
    // "github.com/docker/go-plugins-helpers/graphdriver"
)

var nonCooperativeFsOptions string

var noCoopFsUsage = `Usage:  gear non-cooperative-fs -o lowerdir=...,upperdir=...,workdir=... merged
Central-fs will create a fs that does not cooperate with other fs
  --non-cooperative-fs-options   -o              lowerdir, upperdir and workdir
`

func init() {
    rootCmd.AddCommand(noCoopFsCmd)
    noCoopFsCmd.SetUsageTemplate(noCoopFsUsage)
    noCoopFsCmd.Flags().StringVarP(&nonCooperativeFsOptions, "nonCooperativeFsOptions", "o", "", "lowerdir, upperdir and workdir")
    noCoopFsCmd.MarkFlagRequired("nonCooperativeFsOptions")
}   

var noCoopFsCmd = &cobra.Command{
    Use:   "non-cooperative-fs",
    Short: "create a gear fs that does not cooperate with other fs",
    Long:  `create a gear fs that does not cooperate with other fs`,
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        mergedDir := args[0]
        lowerDir, upperDir, workDir, publicDir := parseArgs(nonCooperativeFsOptions)
        
        nonCooperativeFs.Mount(lowerDir, upperDir, workDir, mergedDir, publicDir)
    },
}

func parseArgs(args string) (string, string, string, string) {
    argStrings := strings.Split(args, ",")

    lowerDir := ""
    upperDir := ""
    workDir := ""
    publicDir := ""

    for _, argString := range argStrings {
        s := strings.TrimPrefix(argString, "lowerdir=")
        if len(s) < len(argString) {
            lowerDir = s
            continue
        }

        s = strings.TrimPrefix(argString, "upperdir=")
        if len(s) < len(argString) {
            upperDir = s
            continue
        }

        s = strings.TrimPrefix(argString, "workdir=")
        if len(s) < len(argString) {
            workDir = s
            continue
        }

        s = strings.TrimPrefix(argString, "publicdir=")
        if len(s) < len(argString) {
            publicDir = s
            continue
        }
    }

    return lowerDir, upperDir, workDir, publicDir
}













