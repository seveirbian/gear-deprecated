package pull

import (
    "os"
    // "io"
    // "crypto/sha256"
    // "fmt"
    // "strings"
    "path/filepath"

    // "golang.org/x/net/context"

    // "github.com/seveirbian/gear/types"
    // "github.com/docker/docker/client"
    "github.com/sirupsen/logrus"
    "github.com/seveirbian/gear/pkg/gear"
    "github.com/seveirbian/gear/pkg/http_utils"
    // dtypes "github.com/docker/docker/api/types"
)

type Puller struct {
    FileName string

    GearRootPath string       // $HOME/.gear/
    ImageDir string             // $HOME/.gear/image/

    PullURL string
}

func InitPuller(fileName string, url string) *Puller {
    // 1. create image dir
    imageDir := filepath.Join(gear.GearRootPath, "image")
    err := os.MkdirAll(imageDir, os.ModePerm)
    if err != nil {
        logrus.WithFields(logrus.Fields{
                "err": err,
                }).Fatal("Fail to create tmpDir:/home/.gears/tmp.")
    }

    return &Puller {
        FileName: fileName, 
        GearRootPath: gear.GearRootPath, 
        ImageDir: imageDir, 
        PullURL: url, 
    }
}

func (b *Puller) Pull() {
    downloadURL := b.PullURL+"/"+b.FileName

    http_utils.Download(filepath.Join(b.ImageDir, b.FileName), downloadURL)
}


















