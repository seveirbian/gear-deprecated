package types

import (
    
)

type Image struct {
    // image name
    Name string
    // image tag
    Tag string
    // if image has no tag, there must be a digest
    Digest string
}