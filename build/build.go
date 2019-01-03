package build

import (
    "fmt"
    "os"
    "io"
    "bytes"
    "archive/tar"
    "path/filepath"
    "crypto/sha256"
    "encoding/json"
    "strings"
    // "bytes"
    "golang.org/x/net/context"

    "github.com/seveirbian/gear/types"
    "github.com/seveirbian/gear/pkg/gear"
    "github.com/seveirbian/gear/pkg/archive"
    "github.com/sirupsen/logrus"

    "github.com/docker/docker/client"
    dtypes "github.com/docker/docker/api/types"
    // darchive "github.com/docker/docker/pkg/archive"
    // "github.com/docker/docker/pkg/progress"
    // "github.com/docker/docker/pkg/streamformatter"
    // "github.com/docker/docker/api/types/container"
)

type Builder struct {
    DockerImage types.Image
    DockerImageInfo dtypes.ImageInspect

    Ctx context.Context
    Client *client.Client

    GearRootPath string       // $HOME/.gear/
    TmpDir string             // $HOME/.gear/tmp/
    // TmpTarPath string         // $HOME/.gear/tmp/tmp.tar

    RegularFiles map[string]string
    IrregularFiles map[string]os.FileInfo

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

    // 3. create gear.json
    logrus.Info("Creating gear.json...")
    b.InitGearJSON()

    // 4. create Dockerfile
    logrus.Info("Generating dockerfile...")
    b.InitDockerfile()

    // 5. create the gear image
    logrus.Info("Creating new gear image...")
    b.BuildGearImage()

    // 6. destroy tmp files
    logrus.Info("Cleaning...")
    // b.Destroy()

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
    var irregularFiles = map[string]os.FileInfo{}

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

                    regularFiles[fmt.Sprintf("%x", h.Sum(nil))] = path
                }else {
                    // record the irregular files
                    irregularFiles[path] = f
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
    // 0. cut long path
    for hash, path := range b.RegularFiles {
        var finalPath string
        pathSlice := strings.SplitAfter(path, "diff")
        for i := 1; i < len(pathSlice); i++ {
            finalPath += pathSlice[i]
        }
        b.RegularFiles[hash] = finalPath
    }

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
    defer f.Close()

    // 3. write to gear.json
    _, err = f.Write(json)
    if err != nil {
        logrus.WithFields(logrus.Fields{
                "err": err,
                }).Fatal("Fail to write gear.json...")
    }
}

func (b *Builder) Destroy() {
    // 1. remove tmp dir
    err := os.RemoveAll(b.TmpDir)
    if err != nil {
        logrus.WithFields(logrus.Fields{
                "err": err,
                }).Fatal("Fail to remove tmp dir...")
    }
}

func (b *Builder) InitDockerfile() {
    // 1. fill b.Dockerfile struct
    b.Dockerfile.FROM = "scratch"
    b.Dockerfile.ENV = b.DockerImageInfo.Config.Env
    b.Dockerfile.LABEL = b.DockerImageInfo.Config.Labels
    b.Dockerfile.VOLUME = b.DockerImageInfo.Config.Volumes
    b.Dockerfile.WORKDIR = b.DockerImageInfo.Config.WorkingDir

    b.Dockerfile.EXPOSE = map[string]struct{}{}
    exposedPorts := b.DockerImageInfo.Config.ExposedPorts
    for key, value := range exposedPorts {
        b.Dockerfile.EXPOSE[string(key)] = value
    }

    entryPoints := b.DockerImageInfo.Config.Entrypoint
    for _, value := range entryPoints {
        b.Dockerfile.ENTRYPOINT = append(b.Dockerfile.ENTRYPOINT, string(value))
    }

    cmds := b.DockerImageInfo.Config.Cmd
    for _, value := range cmds {
        b.Dockerfile.CMD = append(b.Dockerfile.CMD, string(value))
    }

    // 2. transform b.Dockerfile to []byte
    var dockerfile string

    dockerfile = dockerfile + "FROM "
    dockerfile = dockerfile + b.Dockerfile.FROM
    dockerfile = dockerfile + "\n"

    for _, env := range b.Dockerfile.ENV {
        dockerfile = dockerfile + "ENV "
        dockerfile = dockerfile + env
        dockerfile = dockerfile + "\n"
    }

    // Label is a tag to record some thing about this image, useless
    // for key, value := range b.Dockerfile.LABEL {
    //     dockerfile = dockerfile + "LABEL "
    //     dockerfile = dockerfile + env
    //     dockerfile = dockerfile + "\n"
    // }

    for key, _ := range b.Dockerfile.VOLUME {
        dockerfile = dockerfile + "VOLUME "
        dockerfile = dockerfile + key
        dockerfile = dockerfile + "\n"
    }

    if b.Dockerfile.WORKDIR != "" {
        dockerfile = dockerfile + "WORKDIR "
        dockerfile = dockerfile + b.Dockerfile.WORKDIR
        dockerfile = dockerfile + "\n"
    }

    for key, _ := range b.Dockerfile.EXPOSE {
        dockerfile = dockerfile + "EXPOSE "
        dockerfile = dockerfile + key
        dockerfile = dockerfile + "\n"
    }

    for _, entry := range b.Dockerfile.ENTRYPOINT {
        dockerfile = dockerfile + "ENTRYPOINT "
        dockerfile = dockerfile + entry
        dockerfile = dockerfile + "\n"
    }

    // tar irregular file to a .tar file and move it to TmpDir
    if b.needTarIrregularFiles() {
        dockerfile = dockerfile + "ADD "
        dockerfile = dockerfile + "tmp.tar "
        dockerfile = dockerfile + "gear.json /"
        dockerfile = dockerfile + "\n"
        // dockerfile = dockerfile + "COPY "
        // dockerfile = dockerfile + "gear.json /"
        // dockerfile = dockerfile + "\n"
    }
    

    for _, cmd := range b.Dockerfile.CMD {
        dockerfile = dockerfile + "CMD "
        dockerfile = dockerfile + cmd
        dockerfile = dockerfile + "\n"
    }

    // 3. create Dockerfile
    f, err := os.Create(filepath.Join(b.TmpDir, "Dockerfile"))
    if err != nil {
        logrus.WithFields(logrus.Fields{
                "err": err,
                }).Fatal("Fail to create Dockerfile...")
    }
    defer f.Close()

    // 4. write to Dockerfile
    _, err = f.Write([]byte(dockerfile))
    if err != nil {
        logrus.WithFields(logrus.Fields{
                "err": err,
                }).Fatal("Fail to write Dockerfile...")
    }
}

func (b *Builder) needTarIrregularFiles() bool {
    // test whether need a tarFile
    if len(b.IrregularFiles) == 0 {
        return false
    }
    
    b.TarIrregularFiles()

    return true
}

func (b *Builder) TarIrregularFiles() {
    // 1. create a file, which will store the tar data
    tarFile, err := os.Create(filepath.Join(b.TmpDir, "tmp.tar"))
    if err != nil {
        logrus.WithFields(logrus.Fields{
                "err": err,
                }).Fatal("Fail to create tmp.tar...")
    }
    defer tarFile.Close()

    // 2. create a tar writer on tarFile
    tw := tar.NewWriter(tarFile)
    defer tw.Close()

    // 3. write each file to tar
    for irFilePath, irFileInfo := range b.IrregularFiles {
        mode := irFileInfo.Mode()
        var target string
        if mode & os.ModeSymlink != 0 {
            target, err = os.Readlink(irFilePath)
            if err != nil {
                logrus.WithFields(logrus.Fields{
                    "err": err,
                    }).Fatal("Fail to get symlink's target...")
            }
        }
        // get file header info
        hd, err := tar.FileInfoHeader(irFileInfo, target)
        if err != nil {
            logrus.WithFields(logrus.Fields{
                "err": err,
                }).Fatal("Fail to get file header...")
        }
        // modify hd's name
        var finalPath string
        pathSlice := strings.SplitAfter(irFilePath, "diff")
        for i := 1; i < len(pathSlice); i++ {
            finalPath += pathSlice[i]
        }
        hd.Name = finalPath
        // write file header info
        err = tw.WriteHeader(hd)
        if err != nil {
            logrus.WithFields(logrus.Fields{
                "err": err,
                }).Fatal("Fail to write file header...")
        }
    }
}

func (b *Builder) BuildGearImage() {
    // 1. create a tarball which contains everthing including Dockerfile
    var files = []string {
        filepath.Join(b.TmpDir, "tmp.tar"), 
        filepath.Join(b.TmpDir, "gear.json"), 
        filepath.Join(b.TmpDir, "Dockerfile"), 
    }
    var DockerfileTar = filepath.Join(b.TmpDir, "Dockerfile.tar")
    archive.Archive(files, DockerfileTar)

    // 2. open tar
    buildTar, err := os.Open(DockerfileTar)
    if err != nil {
        logrus.WithFields(logrus.Fields{
                "err": err,
                }).Fatal("Fail to open Dockerfile.tar...")
    }

    // 3. init image build options
    opts := dtypes.ImageBuildOptions{
        // this field must be filled
        Context: buildTar, 
        Tags: []string{b.DockerImage.Name+"-gear:"+b.DockerImage.Tag, }, 
    }

    // 4. start to build
    buildResp, err := b.Client.ImageBuild(b.Ctx, buildTar, opts)
    if err != nil {
        logrus.WithFields(logrus.Fields{
                "err": err,
                }).Fatal("Fail to build gear image...")
    }
    defer buildResp.Body.Close()

    // this cannot be removed
    buf := new(bytes.Buffer)
    io.Copy(buf, buildResp.Body)
}














