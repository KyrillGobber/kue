package api

type RoomResponse struct {
	Errors []any `json:"errors"`
	Data   []struct {
		ID       string `json:"id"`
		IDV1     string `json:"id_v1"`
		Children []struct {
			Rid   string `json:"rid"`
			Rtype string `json:"rtype"`
		} `json:"children"`
		Services []struct {
			Rid   string `json:"rid"`
			Rtype string `json:"rtype"`
		} `json:"services"`
		Metadata struct {
			Name      string `json:"name"`
			Archetype string `json:"archetype"`
		} `json:"metadata"`
		Type string `json:"type"`
	} `json:"data"`
}

type LightGroupResponse struct {
	Errors []any `json:"errors"`
	Data   []struct {
		ID    string `json:"id"`
		IDV1  string `json:"id_v1"`
		Owner struct {
			Rid   string `json:"rid"`
			Rtype string `json:"rtype"`
		} `json:"owner"`
		On struct {
			On bool `json:"on"`
		} `json:"on"`
		Dimming struct {
			Brightness float64 `json:"brightness"`
		} `json:"dimming"`
		DimmingDelta          struct{} `json:"dimming_delta"`
		ColorTemperature      struct{} `json:"color_temperature"`
		ColorTemperatureDelta struct{} `json:"color_temperature_delta"`
		Color                 struct{} `json:"color"`
		Alert                 struct {
			ActionValues []string `json:"action_values"`
		} `json:"alert"`
		Signaling struct {
			SignalValues []string `json:"signal_values"`
		} `json:"signaling"`
		Dynamics struct{} `json:"dynamics"`
		Type     string   `json:"type"`
	} `json:"data"`
}

type ZoneResponse struct {
	Errors []any `json:"errors"`
	Data   []struct {
		ID       string `json:"id"`
		IDV1     string `json:"id_v1"`
		Children []struct {
			Rid   string `json:"rid"`
			Rtype string `json:"rtype"`
		} `json:"children"`
		Services []struct {
			Rid   string `json:"rid"`
			Rtype string `json:"rtype"`
		} `json:"services"`
		Metadata struct {
			Name      string `json:"name"`
			Archetype string `json:"archetype"`
		} `json:"metadata"`
		Type string `json:"type"`
	} `json:"data"`
}

type SceneResponse struct {
	Errors []any `json:"errors"`
	Data   []struct {
		ID      string `json:"id"`
		IDV1    string `json:"id_v1"`
		Actions []struct {
			Target struct {
				Rid   string `json:"rid"`
				Rtype string `json:"rtype"`
			} `json:"target"`
			Action struct {
				On struct {
					On bool `json:"on"`
				} `json:"on"`
				Dimming struct {
					Brightness float64 `json:"brightness"`
				} `json:"dimming"`
				Color struct {
					Xy struct {
						X float64 `json:"x"`
						Y float64 `json:"y"`
					} `json:"xy"`
				} `json:"color"`
			} `json:"action"`
		} `json:"actions"`
		Palette struct {
			Color []struct {
				Color struct {
					Xy struct {
						X float64 `json:"x"`
						Y float64 `json:"y"`
					} `json:"xy"`
				} `json:"color"`
				Dimming struct {
					Brightness float64 `json:"brightness"`
				} `json:"dimming"`
			} `json:"color"`
			Dimming          []any `json:"dimming"`
			ColorTemperature []struct {
				ColorTemperature struct {
					Mirek int `json:"mirek"`
				} `json:"color_temperature"`
				Dimming struct {
					Brightness float64 `json:"brightness"`
				} `json:"dimming"`
			} `json:"color_temperature"`
			Effects   []any `json:"effects"`
			EffectsV2 []any `json:"effects_v2"`
		} `json:"palette"`
		Recall   struct{} `json:"recall"`
		Metadata struct {
			Name  string `json:"name"`
			Image struct {
				Rid   string `json:"rid"`
				Rtype string `json:"rtype"`
			} `json:"image"`
			Appdata string `json:"appdata"`
		} `json:"metadata"`
		Group struct {
			Rid   string `json:"rid"`
			Rtype string `json:"rtype"`
		} `json:"group"`
		Speed       float64 `json:"speed"`
		AutoDynamic bool    `json:"auto_dynamic"`
		Status      struct {
			Active string `json:"active"`
		} `json:"status"`
		Type string `json:"type"`
	} `json:"data"`
}

type SetSceneResponse struct {
	Data []struct {
		Rid   string `json:"rid"`
		Rtype string `json:"rtype"`
	} `json:"data"`
	Errors []any `json:"errors"`
}
