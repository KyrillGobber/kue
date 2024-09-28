package config

type DiscoveryResponse []struct {
	ID                string `json:"id"`
	Internalipaddress string `json:"internalipaddress"`
	Port              int    `json:"port"`
}
