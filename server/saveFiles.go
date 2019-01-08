package server

import (
    "os"
    "io"
    "path/filepath"
    "net/http"
    "github.com/labstack/echo"
    // "github.com/sirupsen/logrus"
    "github.com/seveirbian/gear/pkg/gear"
)

var Files = map[string]string{}
var ServerDir = filepath.Join(gear.GearRootPath, "server")

func saveFiles(c echo.Context) error {
    file, err := c.FormFile("file")
    if err != nil {
        return err
    }

    if _, ok := Files[file.Filename]; ok {
        return c.String(http.StatusOK, "This file has existed in server!")
    }
 
    // Source
    src, err := file.Open()
    if err != nil {
        return err
    }
    defer src.Close()
 
    // Destination
    dst, err := os.Create(filepath.Join(ServerDir, file.Filename))
    if err != nil {
        return err
    }
    defer dst.Close()
 
    // Copy
    if _, err = io.Copy(dst, src); err != nil {
        return err
    }

    Files[file.Filename] = filepath.Join(ServerDir, file.Filename)

    return c.HTML(http.StatusOK, "This file accepted!")
}