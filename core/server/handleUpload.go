package server

import (
    "os"
    "io"
    "path/filepath"
    "net/http"
    "github.com/labstack/echo"
    // "github.com/sirupsen/logrus"
    // "github.com/seveirbian/gear/pkg/gear"
)

func (s *Server) handleUpload(c echo.Context) error {
    file, err := c.FormFile("file")
    if err != nil {
        return err
    }

    if !fileExist(filepath.Join(s.ServerDir, file.Filename)) {
        return c.String(http.StatusOK, "This file has existed!")
    }

    // Source
    src, err := file.Open()
    if err != nil {
        return err
    }
    defer src.Close()
 
    // Destination
    dst, err := os.Create(filepath.Join(s.ServerDir, file.Filename))
    if err != nil {
        return err
    }
    defer dst.Close()
 
    // Copy
    if _, err = io.Copy(dst, src); err != nil {
        return err
    }

    return c.String(http.StatusOK, "This file accepted!")
}

func fileExist(filePath string) bool {
    _, err := os.Stat(filePath)
    if err != nil && os.IsNotExist(err) {
        return false
    }
    return true
}