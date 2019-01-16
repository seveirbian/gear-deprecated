package fs

import (
    "os"
    "fmt"
    "syscall"
    "io/ioutil"
    "path/filepath"
    "encoding/json"
    "os/signal"

    "bazil.org/fuse"
    "bazil.org/fuse/fs"

    "github.com/sirupsen/logrus"
    "github.com/seveirbian/gear/types"
)

func Mount(lowerDir, upperDir, workDir, mergedDir, publicDir string) {
    var gearJson = map[string]types.ExtendFileInfo{}

    // 1. read gear.json file
    data, err := ioutil.ReadFile(filepath.Join(lowerDir, "gear.json"))
    if err != nil {
        if err != nil {
            logrus.WithFields(logrus.Fields{
                "err": err,
                }).Fatal("Fail to read gear.json in lowerDir...")
        }
    }

    // 2. unmarshal gear.json file
    err = json.Unmarshal(data, &gearJson)
    if err != nil {
        if err != nil {
            logrus.WithFields(logrus.Fields{
                "err": err,
                }).Fatal("Fail to unmarshal json file...")
        }
    }

    fmt.Println(gearJson)

    // 3. mount to mergedDir
    c, err := fuse.Mount(mergedDir)
    if err != nil {
        if err != nil {
            logrus.WithFields(logrus.Fields{
                "err": err,
                }).Fatal("Fail to mount to mergedDir...")
        }
    }
    defer c.Close()

    fmt.Println("Mount OK")

    // 4. create signal channel
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigs
        c.Close()
        os.Exit(1)
    }()

    // 5. create fs struct
    fileSystem := &types.FS {
        Files: gearJson, 
        LowerDir: lowerDir, 
        UpperDir: upperDir, 
        WorkDir: workDir, 
        MergedDir: mergedDir, 
        PublicDir: publicDir, 
    }

    // 5. use fs to serving fs requests
    err = fs.Serve(c, fileSystem)
    if err != nil {
        if err != nil {
            logrus.WithFields(logrus.Fields{
                "err": err,
                }).Fatal("Fail to serve filesystem...")
        }
    }

    <- c.Ready
    err = c.MountError
    if err != nil {
        if err != nil {
            logrus.WithFields(logrus.Fields{
                "err": err,
                }).Fatal("Have something to report...")
        }
    }

    fmt.Println("fs started...")
}












