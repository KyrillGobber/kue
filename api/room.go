package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"kyrill.dev/kue/config"
)

func FetchRooms() (*RoomResponse, error) {
    // Create new request
    req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", config.GetConfig().BridgeAddress, RoomUrl), nil)
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
    var response RoomResponse
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

func ToggleRoom(lightGroupId string, state bool) error {
    // PUT /clip/v2/resource/grouped_light/3a3e2753-d496-433e-84b6-15134d520666
    data := []byte(fmt.Sprintf(`{"on": { "on": %t }}`, state))
    req, err := http.NewRequest("PUT", fmt.Sprintf("https://192.168.1.219/clip/v2/resource/grouped_light/%s", lightGroupId), bytes.NewBuffer(data))
    if err != nil {
        return err
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
    res, err := client.Do(req)
    if err != nil {
        return err
    }

	if res.StatusCode == http.StatusOK {
        return nil
	} else {
		return fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}
}
