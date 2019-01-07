package types

import (
    
)

type Image struct {
    // raw imageID
    RawID string
    // image name
    Name string
    // image tag
    Tag string
}

type GearImage struct {
    // gear imageID
    GearID string
    // raw imageID
    RawID string
    // image name
    Name string
    // image tag
    Tag string
}