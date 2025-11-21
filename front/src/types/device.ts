import type { CapabilityType, Capability } from "./capability";


export interface Device {
  id: string;
  name: string;
  address: string;
  address_type: "ble" | string;
  protocol: string;
  adapter_ids: string[];
  created_at: string;
  capabilities: Record<CapabilityType, Capability>;
}


export interface DeviceCreateRequest {
  name: string;
  address: string;
  protocol: string;
  adapter_ids: string[];
}