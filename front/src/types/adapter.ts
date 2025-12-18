import type { State } from "./states";

export interface Adapter {
  id: string;
  name: string;
  state: State;
}