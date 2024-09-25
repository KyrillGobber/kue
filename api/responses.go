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
