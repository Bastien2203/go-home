


export interface Device {
  addr: string;
  name: string;
  type: string;
  parser_type: string;
  running? :boolean 
}

export interface Adapter{
  id: string;
  state: string;
  name: string;
  devices: Device[]
}

export interface Parser{
  name: string
}