import type { CapabilityType, Capability } from "./capability";


export interface Device {
  id: string;
  name: string;
  address: string;
  address_type: "ble" | string;
  adapter_ids: string[];
  created_at: string;
  capabilities: Record<CapabilityType, Capability>;
  last_updated: string;
}


export interface DeviceCreateRequest {
  name: string;
  address: string;
  address_type: string;
  adapter_ids: string[];
}