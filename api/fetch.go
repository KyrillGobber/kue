package api

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"kyrill.dev/kue/config"
)

func FetchMe[T any](url string, id *string) (*T, error) {
    // Create new request
    var fetchUrl string 
    if id != nil {
        fetchUrl = fmt.Sprintf("%s%s/%s", config.GetConfig().BridgeAddress, url, *id)
    } else {
        fetchUrl = fmt.Sprintf("%s%s", config.GetConfig().BridgeAddress, url)
    }
    req, err := http.NewRequest(http.MethodGet, fetchUrl, nil)
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
    var response T
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
