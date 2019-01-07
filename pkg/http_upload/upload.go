package http_upload

import (
    "fmt"
    "bytes"
    "mime/multipart"
    "os"
    "io"
    "net/http"
    "io/ioutil"

    "github.com/sirupsen/logrus"

)

func Upload(file string, url string) error {
    bodyBuf := &bytes.Buffer{}
    bodyWriter := multipart.NewWriter(bodyBuf)

    fileWriter, err := bodyWriter.CreateFormFile("file", file)
    if err != nil {
        fmt.Println("error writing to buffer")
        return err
    }

    fh, err := os.Open(file)
    if err != nil {
        fmt.Println("error opening file")
        return err
    }
    defer fh.Close()
    
    //iocopy
    _, err = io.Copy(fileWriter, fh)
    if err != nil {
        return err
    }

    contentType := bodyWriter.FormDataContentType()
    bodyWriter.Close()

    resp, err := http.Post(url, contentType, bodyBuf)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    resp_body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    logrus.WithFields(logrus.Fields{
            "status": resp.Status,
            "response": string(resp_body), 
            }).Info("pushed to seaweed...")
    return nil
}