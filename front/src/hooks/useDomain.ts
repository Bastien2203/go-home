import { api } from "../services/api";
import type { DeviceCreateRequest } from "../types";
import { useApi } from "./useApi";

export function useAdapters() {
  const h = useApi(api.getAdapters);
  return { adapters: h.data || [], loading: h.loading, error: h.error };
}

export function useScanners() {
  const h = useApi(api.getScanners);
  return { scanners: h.data || [], loading: h.loading };
}

export function useProtocols() {
  const h = useApi(api.getProtocols);
  return { protocols: h.data || [], loading: h.loading };
}

export function useDevices() {
  const h = useApi(api.getDevices);

  const create = async (req: DeviceCreateRequest) => {
    await api.createDevice(req);
    h.refresh();
  };

  const link = async (deviceId: string, adapterId: string) => {
    await api.linkDeviceToAdapter(deviceId, adapterId);
    h.refresh();
  };

  const unlink = async (deviceId: string, adapterId: string) => {
    await api.unlinkDeviceFromAdapter(deviceId, adapterId);
    h.refresh();
  };

  return {
    devices: h.data || [],
    loading: h.loading,
    refreshDevices: h.refresh,
    createDevice: create,
    linkAdapter: link,
    unlinkAdapter: unlink,
  };
}