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

type loadingChannels struct {
    roomDataChan       chan *api.RoomResponse
    lightgroupDataChan chan *api.LightGroupResponse
    scenesDataChan     chan *api.SceneResponse
    stopLoader         chan struct{}
}

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize ui: %v", err)
	}
	defer ui.Close()

    mainData := loadData()

	// Sections
	header := getHeader()
	tabpane := getTabs()
	footer := getFooter()

	roommenu := menu.GetItemMenu(getRoomNames(mainData.Rooms), menu.Coords{X1: 5, Y1: 6, X2: 50, Y2: 30})
	sceneMenu := menu.GetSceneMenu(getSceneNames(mainData.Scenes), menu.Coords{X1: 50, Y1: 6, X2: 100, Y2: 30})
	zoneMenu := menu.GetItemMenu(mainData.Zones, menu.Coords{X1: 5, Y1: 6, X2: 50, Y2: 30})

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
				// Get if the lightgrup is on or off
				for i := range mainData.LightGroups.Data {
					if mainData.LightGroups.Data[i].ID == roomLightgroupId {
						mainData.LightGroups.Data[i].On.On = !mainData.LightGroups.Data[i].On.On

						err := api.ToggleRoom(roomLightgroupId, mainData.LightGroups.Data[i].On.On)
						if err != nil {
							log.Fatal(err)
						}
					}
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
	// termWidth, termHeight := ui.TerminalDimensions()
	// Show loader
	loader := widgets.NewParagraph()
	loader.Text = "Loading..."
	loader.Border = false
	loader.SetRect(30, 0, 60, 3)

	// Create a channel to receive the fetched data
	roomDataChan := make(chan *api.RoomResponse)
	lightgroupDataChan := make(chan *api.LightGroupResponse)
	scenesDataChan := make(chan *api.SceneResponse)
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

	go func() {
		scenes, err := api.FetchScenes()
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
	lightgroupData := <-lightgroupDataChan
	scenesData := <-scenesDataChan
	// Stop the loader, run the actual app
	close(stopLoader)

	// Allocate Data
	rooms := getRoomData(roomData)
	scenes := getSceneDataByRoomOrZone(rooms[0].Id, scenesData)
	zones := []string{"Main+", "Some zone", "again", "other"}
    return ActiveData{Rooms: rooms, LightGroups: lightgroupData, Zones: zones, Scenes: scenes, AllScenes: *scenesData}
}

type ActiveData struct {
	Rooms  []Room
    LightGroups *api.LightGroupResponse
	Zones  []string
	Scenes []Scene
    AllScenes api.SceneResponse
}

type Room struct {
	Id              string
	Name            string
	LightGroup      string
	LightGroupOnOff string
	Type            string
}

type Scene struct {
	Id   string
	Name string
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

func getSceneDataByRoomOrZone(roomOrZoneId string, scenes *api.SceneResponse) []Scene {
	sceneData := []Scene{}
	for _, scene := range scenes.Data {
		if scene.Group.Rid == roomOrZoneId {
			sceneData = append(sceneData, Scene{
				Id:   scene.ID,
				Name: scene.Metadata.Name,
			})
		}
	}
	return sceneData
}

func getSceneNames(scenes []Scene) []string {
	sceneNames := []string{}
	for _, scene := range scenes {
		sceneNames = append(sceneNames, fmt.Sprintf("%s", scene.Name))
	}
	return sceneNames
}

func getRoomNames(rooms []Room) []string {
	roomNames := []string{}
	for i, room := range rooms {
		roomNames = append(roomNames, fmt.Sprintf("[%d] %s", i, room.Name))
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
	footer.Text = "hjkl: navigation | enter: scene select | ESC: back | t: quicktoggle on/off | q: quit"
	footer.SetRect(0, 31, 100, 34)
	footer.Border = true
	return footer
}

func getTabs() *widgets.TabPane {
	tabpane := widgets.NewTabPane("Rooms", "Zones")
	tabpane.SetRect(5, 3, 100, 6)
	tabpane.Border = true
	return tabpane
}
