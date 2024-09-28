package main

import (
	"fmt"

	"kyrill.dev/kue/api"
)

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

