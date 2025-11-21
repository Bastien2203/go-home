import type { Unit } from "./units";

export type CapabilityType = "temperature" | "humidity" | "battery_level" | string;

export interface Capability {
  name: CapabilityType;
  value: any;
  type: string;
  unit?: Unit;
}
