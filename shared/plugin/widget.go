package plugin

type Widget struct {
	ID         string           `json:"id"`
	Type       WidgetType       `json:"type"`
	Name       string           `json:"name"`
	Config     map[string]any   `json:"config"`
	MountPoint WidgetMountPoint `json:"mount_point"`
}

type WidgetType string

type WidgetMountPoint string

const (
	DashboardWidget  WidgetMountPoint = "root_widget"
	DeviceWidget     WidgetMountPoint = "device_widget"
	CapabilityWidget WidgetMountPoint = "capability_widget"
)

const (
	TypeLineChart WidgetType = "line-chart"
	TypeKVList    WidgetType = "key-value-list"
)
