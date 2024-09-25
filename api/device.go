package api

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
)

func FetchDevices() ([]byte, error) {
    // Create new request
    req, err := http.NewRequest("GET", "https://192.168.1.219/clip/v2/resource/device", nil)
    if err != nil {
        return nil, err
    }

    // add headers
    req.Header.Add("hue-application-key", "1zUztGOp6k4Z7K1Krz2RJHlbHEpMYkcjTbmfdrL3")
    req.Header.Add("Content-Type", "application/json")

    //create client (with ssl disabled)
    client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        },
    }    

    // send request
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // read response
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    fmt.Println(string(body))
    return body, nil
}

