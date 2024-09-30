package config

type DiscoveryResponse []struct {
	ID                string `json:"id"`
	Internalipaddress string `json:"internalipaddress"`
	Port              int    `json:"port"`
}

type LinkResponse []struct {
	Error struct {
		Type        int    `json:"type"`
		Address     string `json:"address"`
		Description string `json:"description"`
	} `json:"error"`
	Success struct {
		Username  string `json:"username"`
		Clientkey string `json:"clientkey"`
	} `json:"success"`
}
