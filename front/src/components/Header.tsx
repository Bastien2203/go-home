import { Cpu, Server, Radio, Database } from "lucide-react";
import StatusBadge from "./atoms/StatusBadge";
import type { Adapter } from "../types/adapter";
import type { Device } from "../types/device";
import type { Scanner } from "../types/scanner";
import { api } from "../services/api";


export const Header = (props: {
    adapters: Adapter[];
    scanners: Scanner[];
    devices: Device[];
    adaptersUnavailable: boolean;
    scannersUnavailable: boolean;
    devicesUnavailable: boolean;
}) => {
    const isSystemScanning = props.scanners.some(s => s.state === "running");

    return <header className="sticky top-0 z-50 bg-white/80 backdrop-blur-md border-b border-gray-200 px-6 py-4">
        <div className="container mx-auto flex flex-col md:flex-row items-center justify-between gap-4">

          <div className="flex items-center gap-3">
            <div className="bg-linear-to-br from-primary-600 to-primary-700 p-2.5 rounded-xl text-white shadow-lg shadow-primary-600/20">
              <Cpu size={24} strokeWidth={2.5} />
            </div>
            <div>
              <h1 className="text-xl font-bold tracking-tight text-gray-900 leading-tight">
                GoHome
              </h1>

            </div>
          </div>


          <div className="flex flex-wrap justify-center gap-3">
            <StatusBadge icon={<Server size={16} />} label="Adapters" count={props.adapters.length} noCount={props.adaptersUnavailable} />
            <StatusBadge icon={<Radio size={16} />} label="Scanners" count={props.scanners.length} active={isSystemScanning} noCount={props.scannersUnavailable} />
            <StatusBadge icon={<Database size={16} />} label="Appareils" count={props.devices.length} noCount={props.devicesUnavailable} />
          </div>
          <button className="flex" onClick={() => {
              api.logout()
              document.location.reload()
            }}>logout</button>
        </div>
      </header>
}