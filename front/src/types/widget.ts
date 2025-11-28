

export type Widget = {
    id: string;
    type: WidgetType;
    mount_point: MountPoint;
    name: string;
    config: any
}

export type WidgetType = "line-chart"

export type MountPoint = "root_widget" | "device_widget" | "capability_widget"