package types

import (
    "github.com/docker/docker/api/types/strslice"
    "github.com/docker/go-connections/nat"
)

// Port is a string containing port number and protocol in the format "80/tcp"
type Port string

type DockerFile struct {
    FROM string
    ENV []string
    RUN []string
    LABELS map[string]string
    EXPOSE nat.PortSet             // map[Port]struct{}
    ENTRYPOINT strslice.StrSlice   // []string
    VOLUME map[string]struct{}
    WORKDIR string
    CMD strslice.StrSlice          // []string
}