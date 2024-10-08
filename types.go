package main

import "kyrill.dev/kue/api"

type ActiveData struct {
	Rooms  []RoomOrZone
	Zones  []RoomOrZone
    LightGroups *api.LightGroupResponse
	Scenes []Scene
    AllScenes api.SceneResponse
}

type RoomOrZone struct {
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

type loadingChannels struct {
    roomDataChan       chan *api.RoomResponse
    lightgroupDataChan chan *api.LightGroupResponse
    scenesDataChan     chan *api.SceneResponse
    stopLoader         chan struct{}
}
