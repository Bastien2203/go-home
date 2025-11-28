import { API_HOST, API_PROTOCOL } from "../services/api";
import type { Widget } from "../types/widget"
import { WIDGET_REGISTRY } from "../widgetRegistry";
import { Frame } from "./layouts/Frame";


export const DynamicWidget = (props: {
    widget: Widget;
    deviceId?: string;
    capabilityType?: string;
}) => {
    const w = WIDGET_REGISTRY[props.widget.type]
    const config = formatDataUrl(props.widget.config, props.deviceId, props.capabilityType)

    return <Frame title={props.widget.name} icon={w.icon} padding>
        <w.component
            {...config}
        />
    </Frame>
}

const formatDataUrl = (config: Record<string, any>, deviceId?: string, capabilityType?: string): Record<string, any>=> {
    if (!config["dataUrl"]) {
        return config
    }
    let dataUrl = config["dataUrl"] as String

    if (deviceId) {
        dataUrl = dataUrl.replaceAll("{deviceId}", deviceId)
    }

    if (capabilityType) {
        dataUrl = dataUrl.replaceAll("{capabilityType}", capabilityType)
    }

    return {
        ...config,
        dataUrl: `${API_PROTOCOL}://${API_HOST}${dataUrl}`
    }
}