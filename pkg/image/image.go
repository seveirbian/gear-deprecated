package image

import (
    "strings"
    "github.com/sirupsen/logrus"
    "github.com/seveirbian/gear/types"
)

func Parse(image string) types.Image{
    // is image valid
    imageInfo := strings.Split(image, ":")
    if len(imageInfo) <= 1 || len(imageInfo) > 2{
        logrus.WithFields(logrus.Fields{
            "image": image,
            }).Fatal("Invalid imagename...Valid image name should like image:tag...")
    }

    var parsedImage = types.Image{}

    // image:tag
    parsedImage.RawID = image
    parsedImage.Name = imageInfo[0]
    parsedImage.Tag = imageInfo[1]

    return parsedImage
}

func ParseGearImage(image string) types.GearImage{
    tmpImage := Parse(image)
    parsedImage := types.GearImage{
        RawID: tmpImage.RawID, 
        Name: tmpImage.Name, 
        Tag: tmpImage.Tag, 
    }

    if strings.HasSuffix(parsedImage.RawID, "-gear") {
        parsedImage.GearID = parsedImage.RawID
        parsedImage.RawID = strings.Split(parsedImage.RawID, "-gear")[0]
        return parsedImage
    } else {
        logrus.WithFields(logrus.Fields{
            "image": parsedImage,
            }).Fatal("Invalid image...You should use a gear image...")
    }

    return parsedImage
}