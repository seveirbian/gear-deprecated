package http_utils

import (
    "os"
    "io"
    "net/http"
    // "io/ioutil"

    // "github.com/sirupsen/logrus"

)

// file is path + filename, and url must be a complete one
func Download(file string, url string) error {
    // Create the file
    out, err := os.Create(file)
    if err != nil {
        return err
    }
    defer out.Close()

    // Get the data
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Write the body to file
    _, err = io.Copy(out, resp.Body)
    if err != nil {
        return err
    }

    return nil
}