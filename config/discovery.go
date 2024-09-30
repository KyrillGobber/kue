package config

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"kyrill.dev/kue/uiElements"
)

func discoveryProcess() *ConfigType {
	discoveryData := discoverBridges()
	bridgeAddresses := []string{}
	for _, bridge := range *discoveryData {
		bridgeAddresses = append(bridgeAddresses, fmt.Sprintf("%s (%s)", bridge.ID, bridge.Internalipaddress))
	}
	bridgesList := getDiscoveryMenu(bridgeAddresses)
	ui.Render(getBridgeHeader(), getGuideHeader(), bridgesList)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			ui.Close()
			return nil
		case "<Up>", "k":
			bridgesList.ScrollUp()
		case "<Down>", "j":
			bridgesList.ScrollDown()
		case "<Enter>":
			address := (*discoveryData)[bridgesList.SelectedRow].Internalipaddress
			linkData, err := PostNewUserToBridge(address)
			if err != nil {
				uiElements.ShowMessage(err.Error())
			} else if (*linkData)[0].Error.Description != "" {
				uiElements.ShowMessage((*linkData)[0].Error.Description)
			} else {
				userName := (*linkData)[0].Success.Username
				clientKey := (*linkData)[0].Success.Clientkey
				return &ConfigType{
					BridgeAddress: address,
					UserName:      userName,
					ClientKey:     clientKey,
				}
			}
		}
		ui.Render(getBridgeHeader(), getGuideHeader(), bridgesList)
	}
}

func getBridgeHeader() *widgets.Paragraph {
	header := widgets.NewParagraph()
	header.Text = "Kue, your CLI Hue controller"
	header.SetRect(0, 0, 100, 3)
	header.Border = true
	header.TextStyle.Fg = ui.ColorGreen
	return header
}

func getGuideHeader() *widgets.Paragraph {
	header := widgets.NewParagraph()
	header.Text = "Press the linkbutton on your bridge and select it from the list"
	header.SetRect(0, 3, 100, 6)
	header.Border = true
	header.TextStyle.Fg = ui.ColorYellow
	return header
}

func getDiscoveryMenu(bridges []string) *widgets.List {
	bridgesList := widgets.NewList()
	bridgesList.Rows = bridges
	bridgesList.TextStyle = ui.NewStyle(ui.ColorWhite)
	bridgesList.SelectedRowStyle.Fg = ui.ColorYellow
	bridgesList.WrapText = false
	bridgesList.SetRect(0, 6, 100, 50)

	return bridgesList
}

func callDiscoveryEndpoint() (*DiscoveryResponse, error) {
	// Create new request
	req, err := http.NewRequest(http.MethodGet, "https://discovery.meethue.com", nil)
	if err != nil {
		return nil, err
	}

	// add headers
	req.Header.Add("Content-Type", "application/json")

	// create client (with ssl disabled)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// send request
	var response DiscoveryResponse
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

func PostNewUserToBridge(address string) (*LinkResponse, error) {
	// Create new request
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	data := []byte(fmt.Sprintf(`{"devicetype":"%s", "generateclientkey":true}`, "kue-"+hostname))
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/api", address), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	// add headers
	req.Header.Add("Content-Type", "application/json")

	// create client (with ssl disabled)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// send request
	var response LinkResponse
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

func discoverBridges() *DiscoveryResponse {
	// Show loader
	loader := widgets.NewParagraph()
	loader.Border = false
	loader.SetRect(30, 0, 60, 3)

	// Create a channel to receive the fetched data
	discoverChan := make(chan *DiscoveryResponse)
	stopLoader := make(chan struct{})

	// Start fetching data in a goroutine
	go func() {
		discoveredBridges, err := callDiscoveryEndpoint()
		if err != nil {
			log.Panic(err)
		}
		// time.Sleep(1 * time.Second)
		discoverChan <- discoveredBridges
	}()

	ui.Render(loader)
	// Animation for the loader
	go func() {
		frames := []string{"|", "/", "-", "\\"}
		i := 0
		for {
			select {
			case <-stopLoader:
				ui.Clear()
				return
			default:
				loader.Text = fmt.Sprintf("Finding your hue bridges... %s", frames[i])
				ui.Render(loader)
				time.Sleep(100 * time.Millisecond)
				i = (i + 1) % len(frames)
			}
		}
	}()

	// Wait for data to be fetched
	discoveryData := <-discoverChan
	// Stop the loader, run the actual app
	close(stopLoader)

	return discoveryData
}
