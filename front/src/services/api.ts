import type { Device, Adapter, Parser } from "../types";



export class ApiService {
  private baseUrl: string;

  constructor(baseUrl: string = "http://localhost:8080") {
    this.getParsers = this.getParsers.bind(this);
    this.getAdapters = this.getAdapters.bind(this);
    this.getDevices = this.getDevices.bind(this);
    this.createDevice = this.createDevice.bind(this);
    this.startDevice = this.startDevice.bind(this)
    this.linkDeviceToAdapter = this.linkDeviceToAdapter.bind(this)
    this.unlinkDeviceToAdapter = this.unlinkDeviceToAdapter.bind(this)
    this.stopAdapter = this.stopAdapter.bind(this)
    this.restartAdapter = this.restartAdapter.bind(this)
    this.deviceTypes = this.deviceTypes.bind(this)
    this.baseUrl = baseUrl;
  }

  private async getJson<T>(path: string): Promise<T> {
    const res = await fetch(`${this.baseUrl}${path}`);
    if (!res.ok) throw new Error(`Failed to fetch ${path}: ${res.status}`);
    return res.json();
  }

  private async post<T>(path: string, body: any): Promise<T> {
    const res = await fetch(`${this.baseUrl}${path}`, { method: "POST", body: JSON.stringify(body) });
    if (!res.ok) throw new Error(`Failed to post ${path}: ${res.status}`);
    return res.json();
  }

  async getAdapters(): Promise<Adapter[]> {
    return this.getJson<Adapter[]>("/adapters");
  }

  async getParsers(): Promise<Parser[]> {
    return this.getJson<Parser[]>("/parsers");
  }

  async deviceTypes(): Promise<Record<string, string>> {
    return this.getJson<Record<string, string>>("/device-types");
  }

  async getDevices(): Promise<Device[]> {
    return this.getJson<Device[]>("/devices");
  }

  async createDevice(device: Device): Promise<any> {
    this.post("/devices", device)
  }

  async startDevice(device: Device): Promise<any> {
    this.post(`/devices/${device.addr}/start`, {})
  }

  async startAdapter(adapter: Adapter): Promise<any> {
    this.post(`/adapter/${adapter.id}/start`, {})
  }

  async stopAdapter(adapter: Adapter): Promise<any> {
    this.post(`/adapter/${adapter.id}/stop`, {})
  }

  async restartAdapter(adapter: Adapter): Promise<any> {
    this.post(`/adapter/${adapter.id}/restart`, {})
  }

  async linkDeviceToAdapter(device: Device, adapter: Adapter): Promise<any> {
    this.post(`/devices/${device.addr}/link`, {
      adapter_id: adapter.id
    })
  }

  async unlinkDeviceToAdapter(device: Device, adapter: Adapter): Promise<any> {
    this.post(`/devices/${device.addr}/unlink`, {
      adapter_id: adapter.id
    })
  }
}
export const api = new ApiService()