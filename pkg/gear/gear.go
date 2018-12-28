package gear

import (
    "os"

    "github.com/sirupsen/logrus"
)

var GearRootPath = os.Getenv("HOME") + "/.gear/"
var GearParsedImagesPath = os.Getenv("HOME") + "/.gear/parsedImages/"


func Init() {
    // create GearRootPath
    _, err := os.Stat(GearRootPath)
    if err != nil {
        err = os.MkdirAll(GearRootPath, os.ModePerm)
        if err != nil {
            logrus.WithFields(logrus.Fields{
                    "err": err,
                    }).Fatal("Fail to create GearRootPath:/home/.gears/.")
        }
    }
    logrus.WithFields(logrus.Fields{
            }).Info("Created GearRootPath:/home/.gears/.")

    // create ParsedImagesPath
    _, err = os.Stat(GearParsedImagesPath)
    if err != nil {
        err = os.MkdirAll(GearParsedImagesPath, os.ModePerm)
        if err != nil {
            logrus.WithFields(logrus.Fields{
                    "err": err,
                    }).Fatal("Fail to create GearparsedImagesPath:/home/.gears/parsedImages/.")
        }
    }
    logrus.WithFields(logrus.Fields{
            }).Info("Created GearparsedImagesPath:/home/.gears/parsedImages/.")
}