package cmd

import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
    "github.com/sirupsen/logrus"
)

var loglevel string

var rootCmd = &cobra.Command{
    Use:   "gear",
    Short: "Gear is a fast docker container deployment system for edge computing",
    Long: `A fast docker container deployment system for edge computing.
Complete documentation is available at https://github.com/seveirbian/gear
    #######  ######      #      ######        
   ##       ##         ###     ##  ##    
  ## ####  #####     ## ##    ######     
 ##   ##  ##       ## # ##   ## ##       
#######  ######  ##     ##  ##  ###`,
}

func init() {
    rootCmd.Flags().StringVarP(&loglevel, "log-level", "l", "info", 
        "Set the logging level (\"debug\"|\"info\"|\"warn\"|\"error\"|\"fatal\") (default \"info\")")
    switch loglevel {
        case "debug": logrus.SetLevel(logrus.DebugLevel)
        case "info": logrus.SetLevel(logrus.InfoLevel)
        case "warn": logrus.SetLevel(logrus.WarnLevel)
        case "error": logrus.SetLevel(logrus.ErrorLevel)
        case "fatal": logrus.SetLevel(logrus.FatalLevel)
    } 
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}