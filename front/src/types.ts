
export type CapabilityType = "temperature" | "humidity" | "battery_level" | string;

export interface Capability {
  name: CapabilityType;
  value: any;
  type: string;
  unit?: string;
}

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


export interface Adapter {
  id: string;
  name: string;
}

export interface Scanner {
  id: string;
  name: string;
  is_running: boolean;
}

export interface Protocol {
  id: string;
  name: string;
}

export interface DeviceCreateRequest {
  name: string;
  address: string;
  protocol: string;
  adapter_ids: string[];
}