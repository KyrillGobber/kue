package api

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"

	"kyrill.dev/kue/config"
)

func ToggleRoom(lightGroupId string, state bool) error {
	data := []byte(fmt.Sprintf(`{"on": { "on": %t }}`, state))
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s%s/%s", config.GetConfig().BridgeAddress, LightGroupUrl, lightGroupId), bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	// add headers
	req.Header.Add("hue-application-key", config.GetConfig().UserName)
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
