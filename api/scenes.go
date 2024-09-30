package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func SetSceneForRoom(sceneId string) (*SetSceneResponse, error){
    // Create new request
    data := []byte(`{"recall":{"action": "active"}}`)

    req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("https://192.168.1.219/clip/v2/resource/scene/%s", sceneId), bytes.NewBuffer(data))
    if err != nil {
        return nil, err
    }
    // add headers
    req.Header.Add("hue-application-key", "1zUztGOp6k4Z7K1Krz2RJHlbHEpMYkcjTbmfdrL3")
    req.Header.Add("Content-Type", "application/json")

    // add body
    
    //create client (with ssl disabled)
    client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        },
    }

    // send request
    var response SetSceneResponse
    res, err := client.Do(req)
    if err != nil {
        return nil, err
    }

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(bodyBytes, &response)
		if err != nil {
			return nil, err
		}

	} else {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	return &response, nil
}
