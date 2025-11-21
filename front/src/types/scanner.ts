import type { State } from "./states";


export interface Scanner {
  id: string;
  name: string;
  state: State;
}
