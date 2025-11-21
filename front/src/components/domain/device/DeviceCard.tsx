import { useState } from "react";
import type { Adapter } from "../../../types/adapter";
import type { Device } from "../../../types/device";
import { AdapterToggle } from "./AdapterToggle";
import { DeviceCapabilities } from "./DeviceCapabilities";
import { DeviceHeader } from "./DeviceHeader";


type DeviceCardProps = {
  device: Device;
  adapters: Adapter[];
  onLink: (dId: string, aId: string) => void;
  onUnlink: (dId: string, aId: string) => void;
};

export const DeviceCard = ({
  device,
  adapters,
  onLink,
  onUnlink,
}: DeviceCardProps) => {  
  return (
    <div className="flex flex-col justify-between bg-white rounded-xl border border-gray-200 shadow-sm hover:shadow-md transition-shadow duration-300 overflow-hidden h-full">
      
      <DeviceHeader device={device} />
      
    
      <DeviceCapabilities device={device} />

      <div className="px-5 py-4 bg-gray-50/80 border-t border-gray-100">
        <div className="flex items-center justify-between mb-3">
            <h4 className="text-[10px] font-bold text-gray-400 uppercase tracking-wider">
                Connectivité
            </h4>
            <span className="text-[10px] text-gray-400 font-medium">
                {device.adapter_ids?.length || 0} actif(s)
            </span>
        </div>

        <div className="flex flex-wrap gap-2">
          {adapters.length > 0 ? (
            adapters.map((adapter) => {
              const isLinked = device.adapter_ids?.includes(adapter.id);
              return (
                <AdapterToggle
                  key={adapter.id}
                  adapter={adapter}
                  isLinked={!!isLinked}
                  onToggle={() =>
                    isLinked
                      ? onUnlink(device.id, adapter.id)
                      : onLink(device.id, adapter.id)
                  }
                />
              );
            })
          ) : (
             <div className="text-xs text-gray-400 italic w-full text-center py-1">
                Aucun adaptateur configuré
             </div>
          )}
        </div>
      </div>
    </div>
  );
};