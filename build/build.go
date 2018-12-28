package build

import (
    "fmt"
    "golang.org/x/net/context"

    "github.com/seveirbian/gear/types"
    "github.com/sirupsen/logrus"

    "github.com/docker/docker/client"
    dockerTypes "github.com/docker/docker/api/types"
    // "github.com/docker/docker/api/types/container"
)

type Builder struct {
    OldImage types.Image
}

func InitBuilder(image types.Image) *Builder {
    return &Builder{ OldImage: image}
}

func (b *Builder) Build() {
    ctx := context.Background()
    cli, err := client.NewEnvClient()
    if err != nil {
        logrus.WithFields(logrus.Fields{
            "err": err,
            }).Fatal("Fail to create client...")
    }

    fmt.Println(cli.ContainerList(ctx, dockerTypes.ContainerListOptions{}))
}