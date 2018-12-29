package gear

import (
    "os"

    "github.com/sirupsen/logrus"
)

var GearRootPath = filepath.Join(os.Getenv("HOME"), ".gear")
// var GearParsedImagesPath = filepath.Join(os.Getenv("HOME"), ".gear", "parsedImages")


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

    // // create ParsedImagesPath
    // _, err = os.Stat(GearParsedImagesPath)
    // if err != nil {
    //     err = os.MkdirAll(GearParsedImagesPath, os.ModePerm)
    //     if err != nil {
    //         logrus.WithFields(logrus.Fields{
    //                 "err": err,
    //                 }).Fatal("Fail to create GearparsedImagesPath:/home/.gears/parsedImages/.")
    //     }
    // }
}