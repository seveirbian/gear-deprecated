package archive

import (
    "os"
    "io"
    "archive/tar"

    "github.com/sirupsen/logrus"
)

func Archive(files []string, archivePath string) {
    // 1. create a file, which will store the tar data
    tarFile, err := os.Create(archivePath)
    if err != nil {
        logrus.WithFields(logrus.Fields{
                "err": err,
                }).Fatal("Fail to create archive target file...")
    }
    defer tarFile.Close()

    // 2. create a tar writer on tarFile
    tw := tar.NewWriter(tarFile)
    defer tw.Close()

    for _, file := range files {
        fInfo, err := os.Lstat(file)
        if err != nil {
            logrus.WithFields(logrus.Fields{
                    "err": err,
                    }).Fatal("Fail to get file's info...")
        }

        mode := fInfo.Mode()
        var target string

        // if this file is a regular file
        if mode.IsRegular() {
            // get file header info
            hd, err := tar.FileInfoHeader(fInfo, target)
            if err != nil {
                logrus.WithFields(logrus.Fields{
                    "err": err,
                    }).Fatal("Fail to get file header...")
            }
            // write file header info
            err = tw.WriteHeader(hd)
            if err != nil {
                logrus.WithFields(logrus.Fields{
                    "err": err,
                    }).Fatal("Fail to write file header...")
            }
            // open file
            f, err := os.Open(file)
            if err != nil {
                logrus.WithFields(logrus.Fields{
                    "err": err,
                    }).Fatal("Fail to open file...")
            }
            // copy to tar file
            _, err = io.Copy(tw, f)
            if err != nil {
                logrus.WithFields(logrus.Fields{
                    "err": err,
                    }).Fatal("Fail to write tar file...")
            }
        }

        // if this file is a symlink
        if mode & os.ModeSymlink != 0 {
            target, err = os.Readlink(file)
            if err != nil {
                logrus.WithFields(logrus.Fields{
                    "err": err,
                    }).Fatal("Fail to get symlink's target...")
            }
            // get file header info
            hd, err := tar.FileInfoHeader(fInfo, target)
            if err != nil {
                logrus.WithFields(logrus.Fields{
                    "err": err,
                    }).Fatal("Fail to get file header...")
            }
            // write file header info
            err = tw.WriteHeader(hd)
            if err != nil {
                logrus.WithFields(logrus.Fields{
                    "err": err,
                    }).Fatal("Fail to write file header...")
            }
        }

        // if this file is other files
        if mode & os.ModeSymlink != 0 {
            // get file header info
            hd, err := tar.FileInfoHeader(fInfo, target)
            if err != nil {
                logrus.WithFields(logrus.Fields{
                    "err": err,
                    }).Fatal("Fail to get file header...")
            }
            // write file header info
            err = tw.WriteHeader(hd)
            if err != nil {
                logrus.WithFields(logrus.Fields{
                    "err": err,
                    }).Fatal("Fail to write file header...")
            }
        }
    }
} 