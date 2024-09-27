package api

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func FetchZones() (*ZoneResponse, error) {
    // Create new request
    req, err := http.NewRequest("GET", "https://192.168.1.219/clip/v2/resource/zone", nil)
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
    var response ZoneResponse
    res, err := client.Do(req)
    if err != nil {
        return nil, err
    }

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		// Can be = intead of := because we are reassigning the eld err variable, not creating it newly
		err = json.Unmarshal(bodyBytes, &response)
		if err != nil {
			return nil, err
		}

	} else {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	return &response, nil
}
