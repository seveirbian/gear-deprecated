package types

import (
    // "github.com/docker/docker/api/types/strslice"
    // "github.com/docker/go-connections/nat"
)

// Port is a string containing port number and protocol in the format "80/tcp"
type Port string

type Dockerfile struct {
    FROM string
    ENV []string
    RUN []string
    LABELS map[string]string
    EXPOSE map[Port]struct{}
    ENTRYPOINT []string
    VOLUME map[string]struct{}
    WORKDIR string
    CMD []string
}