import { api } from "../services/api";
import type { Adapter, Device } from "../types";
import { useApi } from "./useApi";


export function useAdapters() {
    const h = useApi(api.getAdapters)

    const link = (adapter: Adapter, device: Device) => {
        api.linkDeviceToAdapter(device, adapter).then(h.refresh)
    }

    const unlink = (adapter: Adapter, device: Device) => {
        api.unlinkDeviceToAdapter(device, adapter).then(h.refresh)
    }

    const start = (adapter: Adapter) => {
        api.startAdapter(adapter).then(h.refresh)
    }

    const restart = (adapter: Adapter) => {
        api.restartAdapter(adapter).then(h.refresh)
    }

    const stop = (adapter: Adapter) => {
        api.stopAdapter(adapter).then(h.refresh)
    }

    return {
        adapters: h.data,
        refreshAdapters: h.refresh,
        linkAdapterToDevice: link,
        unlinkAdapterToDevice: unlink,
        startAdapter: start,
        stopAdapter: stop,
        restartAdapter: restart
    }
}