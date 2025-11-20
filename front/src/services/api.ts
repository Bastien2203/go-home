import type { Adapter, Device, DeviceCreateRequest, Protocol, Scanner } from "../types";

export class ApiService {
  private baseUrl: string;

  constructor(baseUrl: string = "http://localhost:8080/api") {
    this.baseUrl = baseUrl;
    
    // Bindings
    this.getAdapters = this.getAdapters.bind(this);
    this.getScanners = this.getScanners.bind(this);
    this.getProtocols = this.getProtocols.bind(this);
    this.getDevices = this.getDevices.bind(this);
    this.createDevice = this.createDevice.bind(this);
    this.linkDeviceToAdapter = this.linkDeviceToAdapter.bind(this);
    this.unlinkDeviceFromAdapter = this.unlinkDeviceFromAdapter.bind(this);
  }

  private async getJson<T>(path: string): Promise<T> {
    const res = await fetch(`${this.baseUrl}${path}`);
    if (!res.ok) throw new Error(`Failed to fetch ${path}: ${res.status}`);
    return res.json();
  }

  private async post<T>(path: string, body: any): Promise<T> {
    const res = await fetch(`${this.baseUrl}${path}`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });
    if (!res.ok) throw new Error(`Failed to post ${path}: ${res.status}`);
    return res.json();
  }

  private async delete(path: string): Promise<void> {
    const res = await fetch(`${this.baseUrl}${path}`, { method: "DELETE" });
    if (!res.ok) throw new Error(`Failed to delete ${path}: ${res.status}`);
  }

  // --- Getters ---
  async getAdapters(): Promise<Adapter[]> {
    return this.getJson<Adapter[]>("/adapters");
  }

  async getScanners(): Promise<Scanner[]> {
    return this.getJson<Scanner[]>("/scanners");
  }

  async getProtocols(): Promise<Protocol[]> {
    return this.getJson<Protocol[]>("/protocols");
  }

  async getDevices(): Promise<Device[]> {
    return this.getJson<Device[]>("/devices");
  }

  // --- Actions ---
  async createDevice(req: DeviceCreateRequest): Promise<Device> {
    return this.post<Device>("/devices", req);
  }

  async linkDeviceToAdapter(deviceId: string, adapterId: string): Promise<void> {
    return this.post(`/devices/${deviceId}/adapters/${adapterId}`, {});
  }

  async unlinkDeviceFromAdapter(deviceId: string, adapterId: string): Promise<void> {
    return this.delete(`/devices/${deviceId}/adapters/${adapterId}`);
  }
}

export const api = new ApiService();