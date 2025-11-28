import type { State } from "./states";
import type { Widget } from "./widget";


export interface Scanner {
  id: string;
  name: string;
  state: State;
  widgets: Record<string, Widget>
}
