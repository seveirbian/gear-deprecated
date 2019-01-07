package push

import (
    "os"
    "io"
    "crypto/sha256"
    "fmt"
    "strings"
    "path/filepath"

    "golang.org/x/net/context"

    "github.com/seveirbian/gear/types"
    "github.com/docker/docker/client"
    "github.com/sirupsen/logrus"
    "github.com/seveirbian/gear/pkg/gear"
    "github.com/seveirbian/gear/pkg/http_upload"
    dtypes "github.com/docker/docker/api/types"
)

type Pusher struct {
    DockerImage types.GearImage
    DockerImageInfo dtypes.ImageInspect

    Ctx context.Context
    Client *client.Client

    GearRootPath string       // $HOME/.gear/
    TmpDir string             // $HOME/.gear/tmp/

    RegularFiles map[string]string

    Ip string
}

func InitPusher(image types.GearImage, ip string) *Pusher {
    // 1. create client to interact with docker daemon
    ctx := context.Background()
    cli, err := client.NewClientWithOpts(client.WithVersion("1.38"))
    if err != nil {
        logrus.WithFields(logrus.Fields{
            "err": err,
            }).Fatal("Fail to create client...")
    }

    // 2. inspect the image, get the info of this image
    imageInspect, _, err := cli.ImageInspectWithRaw(ctx, image.RawID)
    if err != nil {
        logrus.WithFields(logrus.Fields{
            "err": err,
            }).Warn("Fail to inspect the image")
    }

    // 3. create tmp dir and others
    tmpDir := filepath.Join(gear.GearRootPath, "tmp")
    err = os.MkdirAll(tmpDir, os.ModePerm)
    if err != nil {
        logrus.WithFields(logrus.Fields{
                "err": err,
                }).Fatal("Fail to create tmpDir:/home/.gears/tmp.")
    }

    return &Pusher { 
        DockerImage: image, 
        Ctx: ctx, 
        Client: cli, 
        DockerImageInfo: imageInspect, 
        GearRootPath: gear.GearRootPath, 
        TmpDir: tmpDir, 
        Ip: ip, 
    }
}

func (b *Pusher) Push() {
    fmt.Println("Start pushing...")

    // 1. get all layers path of this image
    logrus.Info("Inspecting this image...")
    var layers_path []string
    if b.DockerImageInfo.GraphDriver.Data["LowerDir"] == ""{
        layers_path = append(layers_path, b.DockerImageInfo.GraphDriver.Data["UpperDir"])
    }else {
        layers_path = append(layers_path, strings.Split(b.DockerImageInfo.GraphDriver.Data["LowerDir"], ":")...)
        layers_path = append(layers_path, b.DockerImageInfo.GraphDriver.Data["UpperDir"])
    }

    // 2. walk through these lowerdirs, hash regular files and record irregular files
    logrus.Info("Collecting file information...")
    b.WalkThroughLayers(layers_path)

    // 3. create hard links to regulars
    logrus.Info("Creating sym links to regular files...")
    b.CreateHardlinks()

    // 4. push regular files to seaweedfs
    logrus.Info("Pushing...")
    b.PushFiles()

    // 5. destroy tmp files
    logrus.Info("Cleaning...")
    b.Destroy()

}

// This Func is used to walk through the lowerdirs, calculate regular files' hash 
// and copy irregular file to parsedImages path
func (b *Pusher) WalkThroughLayers(LayerDirs []string) {
    var regularFiles = map[string]string{}
    // var irregularFiles = map[string]os.FileInfo{}

    // 1. get all files of this image
    for _, path := range LayerDirs {
        err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
                if (f == nil) {return err}
                if f.IsDir() {return nil}
                
                // 2. hash each regular file and record other files
                // if this file is a regular file, hash it
                if f.Mode().IsRegular() {
                    f, err := os.Open(path)
                    if err != nil {
                        logrus.WithFields(logrus.Fields{
                                "err": err,
                                }).Fatal("Fail to open file: "+path)
                    }
                    defer f.Close()
                    h := sha256.New()
                    if _, err := io.Copy(h, f); err != nil {
                        logrus.WithFields(logrus.Fields{
                                "err": err,
                                }).Fatal("Fail to copy file: "+path)
                    }

                    regularFiles[path] = fmt.Sprintf("%x", h.Sum(nil))
                }else {
                    // record the irregular files
                    // irregularFiles[path] = f
                }
                return nil
        })
        if err != nil {
            logrus.WithFields(logrus.Fields{
                    "err": err,
                    }).Fatal("Fail to walk image's dirs...")
        }
    }
    
    // 3. assign regularFiles and irregularFiles to builder
    b.RegularFiles = regularFiles
}

// This Func push files to seaweed
func (b *Pusher) PushFiles() {
    for _, hash := range b.RegularFiles {
        http_upload.Upload(filepath.Join(b.TmpDir, hash), b.Ip)
    }
}

// This Func create hardlinks of regular files for pushing
func (b *Pusher) CreateHardlinks() {
    for path, hash := range b.RegularFiles {
        err := os.Symlink(path, filepath.Join(b.TmpDir, hash))
        if err != nil {
            logrus.WithFields(logrus.Fields{
                "err": err,
                }).Fatal("Fail to creating hardlinks...")
        }
    }
}

// This Func remove all files under tmpDir
func (b *Pusher) Destroy() {
    // 1. remove tmp dir
    err := os.RemoveAll(b.TmpDir)
    if err != nil {
        logrus.WithFields(logrus.Fields{
                "err": err,
                }).Fatal("Fail to remove tmp dir...")
    }
}









