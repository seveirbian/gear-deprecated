package server

import (
    "os"
    "path/filepath"

    "github.com/labstack/echo"
    "github.com/sirupsen/logrus"
    "github.com/seveirbian/gear/pkg/gear"
)

type Server struct {
    // server instance
    Server *echo.Echo

    // server root dir
    GearRootPath string    // $HOME/.gear/

    // server Dir
    ServerDir string         // $HOME/.gear/server/

    // server Ip
    Ip string

    // server port
    Port string
}

func InitServer(ip string, port string) *Server{
    // 1. init server path
    serverDir := filepath.Join(gear.GearRootPath, "server")
    err := os.MkdirAll(serverDir, os.ModePerm)
    if err != nil {
        logrus.WithFields(logrus.Fields{
                "err": err,
                }).Fatal("Fail to create tmpDir:/home/.gears/tmp.")
    }

    // 2. create server instance
    s := echo.New()

    return &Server {
        Server: s, 
        GearRootPath: gear.GearRootPath, 
        ServerDir: serverDir, 
        Ip: ip, 
        Port: port, 
    }
}

func (s *Server) InitRoute() {
    s.Server.GET("/", hello)

    s.Server.POST("/files", saveFiles)
}

func (s *Server) Start() {
    s.Server.Start(s.Ip+":"+s.Port)
}







