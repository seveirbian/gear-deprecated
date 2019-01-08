package server

import (
    "net/http"
    "github.com/labstack/echo"
    // "github.com/sirupsen/logrus"
    // "github.com/seveirbian/gear/pkg/gear"
)

func hello(c echo.Context) error {
    return c.String(http.StatusOK, "Hello, World!")
}