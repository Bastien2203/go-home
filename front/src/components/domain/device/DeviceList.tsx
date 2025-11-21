import React from "react";
import { DeviceCard } from "./DeviceCard";
import type { Adapter } from "../../../types/adapter";
import type { Device } from "../../../types/device";
import { ServerOff } from "lucide-react";

interface Props {
  devices: Device[];
  onLink: (deviceId: string, adapterId: string) => void;
  onUnlink: (deviceId: string, adapterId: string) => void;
  adapters: Adapter[];
}

export const DeviceList: React.FC<Props> = (props: Props) => {

  if (props.devices.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center h-64 text-center border-2 border-dashed border-gray-200 rounded-xl bg-gray-50/50">
        <div className="bg-white p-4 rounded-full shadow-sm mb-4">
          <ServerOff size={32} className="text-gray-300" />
        </div>
        <h3 className="text-lg font-medium text-gray-900">Aucun appareil détecté</h3>
        <p className="text-sm text-gray-500 mt-1 max-w-xs">
          Commencez par ajouter un appareil via le bouton "+".
        </p>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
      {props.devices.map((device) => (
        <DeviceCard key={device.id} device={device} adapters={props.adapters} onLink={props.onLink} onUnlink={props.onUnlink} />
      ))}
    </div>
  );
};

