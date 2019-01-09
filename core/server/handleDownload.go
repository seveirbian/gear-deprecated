package server

import (
    // "os"
    // "io"
    "path/filepath"
    "net/http"
    "github.com/labstack/echo"
    // "github.com/sirupsen/logrus"
    // "github.com/seveirbian/gear/pkg/gear"
)

func (s *Server) handleDownload(c echo.Context) error {
    file := c.Param("file")

    if !fileExist(filepath.Join(s.ServerDir, file)) {
        return c.String(http.StatusOK, "This file doesn't exist!")
    }

    return c.File(filepath.Join(s.ServerDir, file))
}