import { api } from "../services/api";
import type { DeviceCreateRequest } from "../types/device";
import { useApi } from "./useApi";

export function useAdapters() {
  const h = useApi(api.getAdapters);
  const stopAdapter = async (id: string) => {
    await api.stopAdapter(id)
    h.refresh()
  }

  const startAdapter = async (id: string) => {
    await api.startAdapter(id)
    h.refresh()
  }
  
  return { adapters: h.data || [], adaptersLoading: h.loading, adaptersError: h.error, startAdapter, stopAdapter };
}

export function useScanners() {
  const h = useApi(api.getScanners);

  const startScanner = async (id: string) => {
    await api.startScanner(id)
    h.refresh()
  }

  const stopScanner = async (id: string) => {
    await api.stopScanner(id)
    h.refresh()
  }

  return { scanners: h.data || [], scannersLoading: h.loading, scannersError: h.error, startScanner, stopScanner};
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

  const deleteDevice = async (deviceId: string) => {
    await api.deleteDevice(deviceId);
    h.refresh();
  };

  return {
    devices: h.data || [],
    devicesLoading: h.loading,
    devicesError: h.error,
    refreshDevices: h.refresh,
    createDevice: create,
    linkAdapter: link,
    unlinkAdapter: unlink,
    deleteDevice: deleteDevice,
  };
}