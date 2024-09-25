package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"kyrill.dev/kue/api"
	"kyrill.dev/kue/menu"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize ui: %v", err)
	}
	defer ui.Close()

	// Show loader
	loader := widgets.NewParagraph()
	loader.Text = "Loading..."
	loader.Border = false
	loader.SetRect(30, 0, 60, 3)

	// Create a channel to receive the fetched data
	roomDataChan := make(chan *api.RoomResponse)
	lightgroupDataChan := make(chan *api.LightGroupResponse)
	stopLoader := make(chan struct{})

	// Start fetching data in a goroutine
	go func() {
		rooms, err := api.FetchRooms()
		if err != nil {
			log.Panic(err)
		}
		// time.Sleep(1 * time.Second)
		roomDataChan <- rooms
	}()

	go func() {
		lightgrups, err := api.FetchLightGroups()
		if err != nil {
			log.Panic(err)
		}
		// time.Sleep(1 * time.Second)
		lightgroupDataChan <- lightgrups
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
				loader.Text = fmt.Sprintf("Loading up your lights... %s", frames[i])
				ui.Render(loader)
				time.Sleep(100 * time.Millisecond)
				i = (i + 1) % len(frames)
			}
		}
	}()

	// Wait for data to be fetched
	roomData := <-roomDataChan
	lightgroupData := <-lightgroupDataChan
	// Stop the loader, run the actual app
	close(stopLoader)

	// Allocate Data
	rooms := getRoomData(roomData)
	zones := []string{"Main+", "Some zone", "again", "other"}
	lights := []string{"light1+", "light2", "light3"}
	mainData := ActiveData{Rooms: rooms, Zones: zones, Lights: lights}

	// Sections
	header := getHeader()
	tabpane := getTabs()
	footer := getFooter()

	roommenu := menu.GetItemMenu(getRoomNames(mainData.Rooms), menu.Coords{X1: 5, Y1: 6, X2: 50, Y2: 30})
	zoneMenu := menu.GetItemMenu(mainData.Zones, menu.Coords{X1: 5, Y1: 6, X2: 50, Y2: 30})
	lightMenu := menu.GetItemMenu(mainData.Lights, menu.Coords{X1: 5, Y1: 6, X2: 50, Y2: 30})

	activeMenu := roommenu
	renderTab := func() {
		switch tabpane.ActiveTabIndex {
		case 0:
			activeMenu = roommenu
			ui.Render(roommenu)
		case 1:
			activeMenu = zoneMenu
			ui.Render(zoneMenu)
		case 2:
			activeMenu = lightMenu
			ui.Render(lightMenu)
		}
	}

	ui.Render(header, tabpane, roommenu, footer)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "h":
			tabpane.FocusLeft()
			ui.Clear()
			ui.Render(header, tabpane, footer)
			renderTab()
		case "l":
			tabpane.FocusRight()
			ui.Clear()
			ui.Render(header, tabpane, footer)
			renderTab()
		case "<Up>", "k":
			activeMenu.ScrollUp()
		case "<Down>", "j":
			activeMenu.ScrollDown()
		case "0", "1", "2", "3", "4":
			newSelected, err := strconv.Atoi(e.ID)
			if err != nil {
				log.Fatal(err)
			}
			activeMenu.SelectedRow = newSelected
		case "t":
			if activeMenu == roommenu {
				roomLightgroupId := mainData.Rooms[activeMenu.SelectedRow].LightGroup
				// Get if the lightgrup is on or off
				for i := range lightgroupData.Data {
					if lightgroupData.Data[i].ID == roomLightgroupId {
						lightgroupData.Data[i].On.On = !lightgroupData.Data[i].On.On

						err := api.ToggleRoom(roomLightgroupId, lightgroupData.Data[i].On.On)
						if err != nil {
							log.Fatal(err)
						}
					}
				}
			}
		case "<Enter>":
			ui.Close()
			return
		}
		ui.Render(header, tabpane, activeMenu, footer)
	}
}

type Room struct {
	Id              string
	Name            string
	LightGroup      string
	LightGroupOnOff string
	Type            string
}

type ActiveData struct {
	Rooms  []Room
	Zones  []string
	Lights []string
}

func getRoomData(rooms *api.RoomResponse) []Room {
	roomData := []Room{}
	for _, room := range rooms.Data {
		var lightGroup string
		for _, service := range room.Services {
			if service.Rtype == "grouped_light" {
				lightGroup = service.Rid
				break
			}
		}

		roomData = append(roomData, Room{
			Id:         room.ID,
			Name:       room.Metadata.Name,
			LightGroup: lightGroup,
			Type:       room.Type,
		})
	}
	return roomData
}

func getRoomNames(rooms []Room) []string {
	roomNames := []string{}
	for i, room := range rooms {
		roomNames = append(roomNames, fmt.Sprintf("%d %s", i, room.Name))
	}
	return roomNames
}

func getHeader() *widgets.Paragraph {
	header := widgets.NewParagraph()
	header.Text = "Kue, your CLI Hue controller"
	header.SetRect(0, 0, 100, 3)
	header.Border = true
	header.TextStyle.Fg = ui.ColorGreen
	return header
}

func getFooter() *widgets.Paragraph {
	footer := widgets.NewParagraph()
	footer.Text = "hjkl: navigation | enter: select | t: quicktoggle on/off | q: quit"
	footer.SetRect(0, 31, 100, 34)
	footer.Border = true
	return footer
}

func getTabs() *widgets.TabPane {
	tabpane := widgets.NewTabPane("Rooms", "Zones", "Lights")
	tabpane.SetRect(5, 3, 50, 6)
	tabpane.Border = true
	return tabpane
}
