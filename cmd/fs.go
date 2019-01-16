package cmd

import (
    // "fmt"
    "strings"

    "github.com/spf13/cobra"
    "github.com/seveirbian/gear/core/fs"
    // "github.com/seveirbian/gear/pkg/image"
    // // "github.com/seveirbian/gear/pkg/gear"
    // "github.com/docker/go-plugins-helpers/graphdriver"
)

var fsOptions string

var fsUsage = `Usage:  gear fs -o lowerdir=...,upperdir=...,workdir=... merged
Central-fs will create a fs that does not cooperate with other fs
  --fsOptions   -o              lowerdir, upperdir and workdir
`

func init() {
    rootCmd.AddCommand(fsCmd)
    fsCmd.SetUsageTemplate(fsUsage)
    fsCmd.Flags().StringVarP(&fsOptions, "fsOptions", "o", "", "lowerdir, upperdir and workdir")
    fsCmd.MarkFlagRequired("fsOptions")
}   

var fsCmd = &cobra.Command{
    Use:   "fs",
    Short: "create a gear fs",
    Long:  `create a gear fs`,
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        mergedDir := args[0]
        lowerDir, upperDir, workDir, publicDir := parseArgs(fsOptions)
        
        fs.Mount(lowerDir, upperDir, workDir, mergedDir, publicDir)
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













