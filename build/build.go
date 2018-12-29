package build

import (
    "fmt"
    "os"
    "io"
    "path/filepath"
    "crypto/sha256"
    "encoding/json"
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
    TmpDir string

    RegularFiles map[string]string
    IrregularFiles []string

    Dockerfile types.Dockerfile
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

    // 3. create tmp dir and others
    tmpDir := filepath.Join(gear.GearRootPath, "tmp")
    err = os.MkdirAll(tmpDir, os.ModePerm)
        if err != nil {
            logrus.WithFields(logrus.Fields{
                    "err": err,
                    }).Fatal("Fail to create tmpDir:/home/.gears/tmp.")
        }

    return &Builder { 
        DockerImage: image, 
        Ctx: ctx, 
        Client: cli, 
        DockerImageInfo: imageInspect, 
        GearRootPath: gear.GearRootPath, 
        TmpDir: tmpDir, 
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

    // 3. create gear.json
    b.InitGearJSON()

    // 4. create Dockerfile
    b.InitDockerfile()

    // 5. create the gear image

    // 6. destroy tmp files
    b.Destroy()

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
    
    // 3. assign regularFiles and irregularFiles to builder
    b.RegularFiles = regularFiles
    b.IrregularFiles = irregularFiles
}

func (b *Builder) InitGearJSON() {
    // 1. encode regularfiles map[string]string
    json, err := json.Marshal(b.RegularFiles)
    if err != nil {
        logrus.WithFields(logrus.Fields{
                "err": err,
                }).Fatal("Fail to encode regularfiles struct...")
    }

    // 2. create the gear.json file
    f, err := os.Create(filepath.Join(b.TmpDir, "gear.json"))
    if err != nil {
        logrus.WithFields(logrus.Fields{
                "err": err,
                }).Fatal("Fail to create gear.json...")
    }

    // 3. write to gear.json
    _, err = f.Write(json)
    if err != nil {
        logrus.WithFields(logrus.Fields{
                "err": err,
                }).Fatal("Fail to write gear.json...")
    }
}

func (b *Builder) Destroy() {
    // 1. unmount overlay

    // 2. remove tmp dir
    err := os.RemoveAll(b.TmpDir)
    if err != nil {
        logrus.WithFields(logrus.Fields{
                "err": err,
                }).Fatal("Fail to remove tmp dir...")
    }
}

func (b * Builder) InitDockerfile() {
    b.Dockerfile.FROM = "scratch"
    b.Dockerfile.ENV = b.DockerImageInfo.Config.Env
    b.Dockerfile.LABELS = b.DockerImageInfo.Config.Labels

    exposedPorts := b.DockerImageInfo.Config.ExposedPorts
    for key, value := range exposedPorts {
        b.Dockerfile.EXPOSE[key] = value
    }

    entryPoints := b.DockerImageInfo.Config.Entrypoint
    for _, value := range entryPoints {
        b.Dockerfile.ENTRYPOINT = append(b.Dockerfile.ENTRYPOINT, string(value))
    }

    b.Dockerfile.VOLUME = b.DockerImageInfo.Config.Volumes
    b.Dockerfile.WORKDIR = b.DockerImageInfo.Config.WorkingDir

    cmds := b.DockerImageInfo.Config.Cmd
    for _, value = range cmds {
        b.Dockerfile.CMD = append(b.Dockerfile.CMD, string(value))
    }

    fmt.Println(b.Dockerfile)
}

















