package types

import (
    // "github.com/docker/docker/api/types/strslice"
    // "github.com/docker/go-connections/nat"
)

type Dockerfile struct {
    FROM string
    ENV []string
    RUN []string
    LABEL map[string]string
    EXPOSE map[string]struct{}   // "80/tcp":{}
    ENTRYPOINT []string
    VOLUME map[string]struct{}
    WORKDIR string
    CMD []string
}