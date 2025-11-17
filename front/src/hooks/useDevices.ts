import { api } from "../services/api";
import type { Device } from "../types";
import { useApi } from "./useApi";


export function useDevices() {
    const devices = useApi(api.getDevices)
    const deviceTypes = useApi(api.deviceTypes)

    const start = (device: Device) => {
        api.startDevice(device).then(devices.refresh)
    }

    return {
        devices: devices.data,
        refreshDevices: devices.refresh,
        startDevice: start,
        deviceTypes: deviceTypes.data
    }
}