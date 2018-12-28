package image

import (
    "strings"
    "github.com/sirupsen/logrus"
    "github.com/seveirbian/gear/types"
)

func Parse(image string) types.Image{
    // is image valid
    imageInfo := strings.Split(image, ":")
    if len(imageInfo) <= 1 || len(imageInfo) > 3{
        logrus.WithFields(logrus.Fields{
            "image": image,
            }).Fatal("Invalid imagename...")
    }

    var parsedImage = types.Image{}

    // image:tag
    if len(imageInfo) == 2 {
        parsedImage.Name = imageInfo[0]
        parsedImage.Tag = imageInfo[1]
        parsedImage.Digest = ""
    } else {
        // image:digest
        parsedImage.Name = imageInfo[0]
        parsedImage.Tag = ""
        parsedImage.Digest = imageInfo[1]+":"+imageInfo[2]
    }

    return parsedImage
}