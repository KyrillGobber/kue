package main

import (
	"fmt"
	"strings"

	"kyrill.dev/kue/api"
)

func getRoomData(rooms *api.RoomResponse) []RoomOrZone {
	roomData := []RoomOrZone{}
	for _, room := range rooms.Data {
		var lightGroup string
		for _, service := range room.Services {
			if service.Rtype == "grouped_light" {
				lightGroup = service.Rid
				break
			}
		}

		roomData = append(roomData, RoomOrZone{
			Id:         room.ID,
			Name:       room.Metadata.Name,
			LightGroup: lightGroup,
			Type:       room.Type,
		})
	}
	return roomData
}

func getZoneData(zones *api.ZoneResponse) []RoomOrZone {
	zoneData := []RoomOrZone{}
	for _, zone := range zones.Data {
		var lightGroup string
		for _, service := range zone.Services {
			if service.Rtype == "grouped_light" {
				lightGroup = service.Rid
				break
			}
		}

		zoneData = append(zoneData, RoomOrZone{
			Id:         zone.ID,
			Name:       zone.Metadata.Name,
			LightGroup: lightGroup,
			Type:       zone.Type,
		})
	}
	return zoneData
}

func getSceneDataByRoomOrZone(roomOrZoneId string, scenes *api.SceneResponse, filter string) []Scene {
	sceneData := []Scene{}
	for _, scene := range scenes.Data {
		if scene.Group.Rid == roomOrZoneId {
            if filter == "" || strings.Contains(strings.ToLower(scene.Metadata.Name), strings.ToLower(filter)) {
				sceneData = append(sceneData, Scene{
					Id:   scene.ID,
					Name: scene.Metadata.Name,
				})
			}
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

func getRoomNames(rooms []RoomOrZone) []string {
	roomNames := []string{}
	for i, room := range rooms {
		roomNames = append(roomNames, fmt.Sprintf("[%d] %s", i, room.Name))
	}
	return roomNames
}

