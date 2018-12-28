package main

import (
    "github.com/seveirbian/gear/cmd"
    "github.com/seveirbian/gear/pkg/gear"
)

func main() {
    gear.Init()

    cmd.Execute()
}