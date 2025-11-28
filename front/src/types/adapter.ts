import type { State } from "./states";
import type { Widget } from "./widget";

export interface Adapter {
  id: string;
  name: string;
  state: State;
  widgets: Record<string, Widget>
}