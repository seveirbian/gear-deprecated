package build

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "golang.org/x/net/context"

    "github.com/seveirbian/gear/types"
    // "github.com/seveirbian/gear/pkg/gear"
    "github.com/sirupsen/logrus"

    "github.com/docker/docker/client"
    dtypes "github.com/docker/docker/api/types"
    // "github.com/docker/docker/api/types/container"
)

type Builder struct {
    DockerImage types.Image
    DockerImageInfo dtypes.ImageInspect

    ctx context.Context
    Client *client.Client
}

func InitBuilder(image types.Image) *Builder {
    // create client to interact with docker daemon
    ctx := context.Background()
    cli, err := client.NewClientWithOpts(client.WithVersion("1.39"))
    if err != nil {
        logrus.WithFields(logrus.Fields{
            "err": err,
            }).Fatal("Fail to create client...")
    }

    // inspect the image, get the info of this image
    imageInspect, _, err := cli.ImageInspectWithRaw(ctx, image.RawID)
    if err != nil {
        logrus.WithFields(logrus.Fields{
            "err": err,
            }).Warn("Fail to inspect the image")
    }

    return &Builder { 
        DockerImage: image, 
        ctx: ctx, 
        Client: cli, 
        DockerImageInfo: imageInspect, 
    }
}

func (b *Builder) Build() {
    fmt.Println("Start building...")

    // get all layers path of this image
    var layers_path []string
    if b.DockerImageInfo.GraphDriver.Data["LowerDir"] == ""{
        layers_path = append(layers_path, b.DockerImageInfo.GraphDriver.Data["UpperDir"])
    }else {
        layers_path = append(layers_path, strings.Split(b.DockerImageInfo.GraphDriver.Data["LowerDir"], ":")...)
        layers_path = append(layers_path, b.DockerImageInfo.GraphDriver.Data["UpperDir"])
    }

    // walk through these lowerdirs
    b.WalkThroughLayers(layers_path)

}

// This func detec whether this image has been built
func (b *Builder) HasParsedThisImage() bool{
    _, _, err := b.Client.ImageInspectWithRaw(b.ctx, b.DockerImage.Name+"-gear"+":"+b.DockerImage.Tag)

    if err != nil {
        return false
    }

    return true
}

// This Func is used to walk through the lowerdirs, calculate regular files' hash 
// and copy irregular file to parsedImages path
func (b *Builder) WalkThroughLayers(lowerDirs []string) {
    // var allFiles []string
    fmt.Println(lowerDirs)
    for _, path := range lowerDirs{
        fmt.Println(path)
        err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
                if ( f == nil ) {return err}
                if f.IsDir() {return nil}
                // allFiles = append(allFiles, path+)
                println(path)
                return nil
        })
        if err != nil {
            logrus.WithFields(logrus.Fields{
                    "err": err,
                    }).Fatal("Fail to walk image's dirs...")
        }
    }
}



















