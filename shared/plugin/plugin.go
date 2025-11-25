package plugin

import "gohome/shared/types"

type Plugin struct {
	ID    string      `json:"id"`
	Name  string      `json:"name"`
	Type  PluginType  `json:"type"`
	State types.State `json:"state"`
}

type PluginType string

const (
	PluginAdapter PluginType = "plugin_adapter"
	PluginScanner PluginType = "plugin_scanner"
)
