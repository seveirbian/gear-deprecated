package build

import (
    "fmt"
    "os"
    "io"
    "path/filepath"
    "crypto/sha256"
    "strings"
    "golang.org/x/net/context"

    "github.com/seveirbian/gear/types"
    "github.com/seveirbian/gear/pkg/gear"
    "github.com/sirupsen/logrus"

    "github.com/docker/docker/client"
    dtypes "github.com/docker/docker/api/types"
    // "github.com/docker/docker/api/types/container"
)

type Builder struct {
    DockerImage types.Image
    DockerImageInfo dtypes.ImageInspect

    Ctx context.Context
    Client *client.Client

    GearRootPath string

    RegularFiles map[string]string
    IrregularFiles []string
}

func InitBuilder(image types.Image) *Builder {
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

    return &Builder { 
        DockerImage: image, 
        Ctx: ctx, 
        Client: cli, 
        DockerImageInfo: imageInspect, 
        GearRootPath: gear.GearRootPath,  
    }
}

func (b *Builder) Build() {
    fmt.Println("Start building...")

    // 1. get all layers path of this image
    var layers_path []string
    if b.DockerImageInfo.GraphDriver.Data["LowerDir"] == ""{
        layers_path = append(layers_path, b.DockerImageInfo.GraphDriver.Data["UpperDir"])
    }else {
        layers_path = append(layers_path, strings.Split(b.DockerImageInfo.GraphDriver.Data["LowerDir"], ":")...)
        layers_path = append(layers_path, b.DockerImageInfo.GraphDriver.Data["UpperDir"])
    }

    // 2. walk through these lowerdirs, hash regular files and record irregular files
    b.WalkThroughLayers(layers_path)

}

// This func detect whether this image has been built
func (b *Builder) HasParsedThisImage() bool{
    _, _, err := b.Client.ImageInspectWithRaw(b.Ctx, b.DockerImage.Name+"-gear"+":"+b.DockerImage.Tag)

    if err != nil {
        return false
    }

    return true
}

// This Func is used to walk through the lowerdirs, calculate regular files' hash 
// and copy irregular file to parsedImages path
func (b *Builder) WalkThroughLayers(LayerDirs []string) {
    var regularFiles = map[string]string{}
    var irregularFiles = []string{}

    // 1. get all files of this image
    for _, path := range LayerDirs {
        fmt.Println(path)
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
                    regularFiles[fmt.Sprintf("%x", h.Sum(nil))] = path
                }else {
                    // record the irregular files
                    irregularFiles = append(irregularFiles, path)
                }
                return nil
        })
        if err != nil {
            logrus.WithFields(logrus.Fields{
                    "err": err,
                    }).Fatal("Fail to walk image's dirs...")
        }
    }
    
    b.RegularFiles = regularFiles
    b.IrregularFiles = irregularFiles
    fmt.Println(regularFiles)
    fmt.Println(irregularFiles)
}



















