package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"kyrill.dev/kue/api"
	"kyrill.dev/kue/config"
	"kyrill.dev/kue/menu"
	"kyrill.dev/kue/uiElements"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize ui: %v", err)
	}
	defer ui.Close()

	// If no config file: search for bridge and save
	err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading or creating config: %v\n", err)
		return
	}

	// Load data
	mainData := loadData()

	// Sections
	header := uiElements.GetHeader()
	tabpane := uiElements.GetTabs()
	footer := uiElements.GetFooter()

	// termWidth, termHeight := ui.TerminalDimensions()
	roommenu := menu.GetItemMenu(getRoomNames(mainData.Rooms), menu.Coords{X1: 5, Y1: 6, X2: 50, Y2: 30})
	sceneMenu := menu.GetSceneMenu(getSceneNames(mainData.Scenes), menu.Coords{X1: 50, Y1: 6, X2: 100, Y2: 30})
	zoneMenu := menu.GetItemMenu(getRoomNames(mainData.Zones), menu.Coords{X1: 5, Y1: 6, X2: 50, Y2: 30})

	activeMenu := roommenu
	renderTab := func() {
		switch tabpane.ActiveTabIndex {
		case 0:
			activeMenu = roommenu
			ui.Render(roommenu)
		case 1:
			activeMenu = zoneMenu
			ui.Render(zoneMenu)
		}
	}

	ui.Render(header, tabpane, roommenu, sceneMenu, footer)

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
			if activeMenu == roommenu {
				mainData.Scenes = getSceneDataByRoomOrZone(mainData.Rooms[activeMenu.SelectedRow].Id, &mainData.AllScenes)
				sceneMenu.Rows = getSceneNames(mainData.Scenes)
			}
		case "<Down>", "j":
			activeMenu.ScrollDown()
			if activeMenu == roommenu {
				mainData.Scenes = getSceneDataByRoomOrZone(mainData.Rooms[activeMenu.SelectedRow].Id, &mainData.AllScenes)
				sceneMenu.Rows = getSceneNames(mainData.Scenes)
			}
		case "0", "1", "2", "3", "4":
			newSelected, err := strconv.Atoi(e.ID)
			if err != nil {
				log.Fatal(err)
			}
			activeMenu.SelectedRow = newSelected
			if activeMenu == roommenu {
				mainData.Scenes = getSceneDataByRoomOrZone(mainData.Rooms[activeMenu.SelectedRow].Id, &mainData.AllScenes)
				sceneMenu.Rows = getSceneNames(mainData.Scenes)
			}
		case "t":
			if activeMenu == roommenu {
				roomLightgroupId := mainData.Rooms[activeMenu.SelectedRow].LightGroup
                toggleLightgroup(roomLightgroupId, mainData)
			}
			if activeMenu == zoneMenu {
				LightgroupId := mainData.Zones[activeMenu.SelectedRow].LightGroup
				// Get if the lightgrup is on or off
				toggleLightgroup(LightgroupId, mainData)
			}
			if activeMenu == sceneMenu {
				sceneId := mainData.Scenes[activeMenu.SelectedRow].Id
				_, err := api.SetSceneForRoom(sceneId)
				if err != nil {
					log.Fatal(err)
				}
            }
		case "<Enter>":
			if activeMenu == sceneMenu {
				sceneId := mainData.Scenes[activeMenu.SelectedRow].Id
				_, err := api.SetSceneForRoom(sceneId)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				activeMenu = sceneMenu
				sceneMenu.BorderStyle = ui.NewStyle(ui.ColorYellow)
			}
		case "<Escape>":
			sceneMenu.BorderStyle = ui.NewStyle(ui.ColorWhite)
			switch tabpane.ActiveTabIndex {
			case 0:
				activeMenu = roommenu
			case 1:
				activeMenu = zoneMenu
			}

		case "g":
			activeMenu.ScrollTop()
		case "G":
			activeMenu.ScrollBottom()
		case "d":
			activeMenu.ScrollHalfPageDown()
		case "u":
			activeMenu.ScrollHalfPageUp()
		}
		ui.Render(header, tabpane, activeMenu, sceneMenu, footer)
	}
}

func loadData() ActiveData {
	// Show loader
	loader := widgets.NewParagraph()
	loader.Text = "Loading..."
	loader.Border = false
	loader.SetRect(30, 0, 60, 3)

	// Create a channel to receive the fetched data
	roomDataChan := make(chan *api.RoomResponse)
	zoneDataChan := make(chan *api.ZoneResponse)
	lightgroupDataChan := make(chan *api.LightGroupResponse)
	scenesDataChan := make(chan *api.SceneResponse)
	stopLoader := make(chan struct{})

	// Start fetching data in a goroutine
	go func() {
		rooms, err := api.FetchMe[api.RoomResponse](api.RoomUrl, nil)
		if err != nil {
			log.Panic(err)
		}
		// time.Sleep(1 * time.Second)
		roomDataChan <- rooms
	}()

	go func() {
		zones, err := api.FetchMe[api.ZoneResponse](api.ZoneUrl, nil)
		if err != nil {
			log.Panic(err)
		}
		// time.Sleep(1 * time.Second)
		zoneDataChan <- zones
	}()

	go func() {
		lightgrups, err := api.FetchMe[api.LightGroupResponse](api.LightGroupUrl, nil)
		if err != nil {
			log.Panic(err)
		}
		// time.Sleep(1 * time.Second)
		lightgroupDataChan <- lightgrups
	}()

	go func() {
		scenes, err := api.FetchMe[api.SceneResponse](api.SceneUrl, nil)
		if err != nil {
			log.Panic(err)
		}
		// time.Sleep(1 * time.Second)
		scenesDataChan <- scenes
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
	zoneData := <-zoneDataChan
	lightgroupData := <-lightgroupDataChan
	scenesData := <-scenesDataChan
	// Stop the loader, run the actual app
	close(stopLoader)

	// Allocate Data
	rooms := getRoomData(roomData)
	scenes := getSceneDataByRoomOrZone(rooms[0].Id, scenesData)
	zones := getZoneData(zoneData)
	return ActiveData{Rooms: rooms, LightGroups: lightgroupData, Zones: zones, Scenes: scenes, AllScenes: *scenesData}
}

func toggleLightgroup(LightgroupId string, mainData ActiveData) {
	for i := range mainData.LightGroups.Data {
		if mainData.LightGroups.Data[i].ID == LightgroupId {
			mainData.LightGroups.Data[i].On.On = !mainData.LightGroups.Data[i].On.On

			err := api.ToggleLightgroup(LightgroupId, mainData.LightGroups.Data[i].On.On)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
