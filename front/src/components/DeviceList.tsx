import React from "react";
import { useAdapters } from "../hooks/useDomain";
import type { Device } from "../types";
import { Activity, Battery, Thermometer, Droplets, Link as LinkIcon } from "lucide-react";

interface Props {
  devices: Device[];
  onLink: (deviceId: string, adapterId: string) => void;
  onUnlink: (deviceId: string, adapterId: string) => void;
}

export const DeviceList: React.FC<Props> = ({ devices, onLink, onUnlink }) => {
  const { adapters } = useAdapters();

  if (devices.length === 0) {
    return <div className="text-gray-500 italic">No devices registered yet.</div>;
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {devices.map((device) => (
        <div key={device.id} className="bg-white rounded-lg shadow-md p-5 border border-gray-100">
          <div className="flex justify-between items-start mb-2">
            <div>
              <h3 className="font-bold text-lg">{device.name}</h3>
              <p className="text-xs text-gray-400 font-mono">{device.address}</p>
              <span className="inline-block bg-gray-100 text-gray-600 text-xs px-2 py-0.5 rounded mt-1">
                {device.protocol}
              </span>
            </div>
            <div className="text-right">
              <div className="text-xs text-gray-400">ID: {device.id.slice(0, 8)}...</div>
            </div>
          </div>

          {/* Capabilities */}
          <div className="my-4 space-y-2">
            {Object.values(device.capabilities || {}).map((cap) => (
              <div key={cap.name} className="flex items-center justify-between text-sm">
                <span className="flex items-center text-gray-600 capitalize">
                  {getIcon(cap.name)}
                  <span className="ml-2">{cap.name}</span>
                </span>
                <span className="font-semibold">
                  {typeof cap.value === 'number'
                    ? Number(cap.value).toFixed(2)
                    : cap.value} <span className="text-xs text-gray-500">{cap.unit}</span>
                </span>
              </div>
            ))}
            {Object.keys(device.capabilities || {}).length === 0 && (
              <div className="text-sm text-gray-400 flex items-center">
                <Activity size={16} className="mr-2" /> No data received yet
              </div>
            )}
          </div>

          {/* Adapters Linking */}
          <div className="border-t pt-3 mt-2">
            <h4 className="text-xs font-semibold text-gray-500 mb-2 uppercase flex items-center">
              <LinkIcon size={12} className="mr-1" /> Linked Adapters
            </h4>
            <div className="flex flex-wrap gap-2">
              {adapters.map((adapter) => {
                const isLinked = device.adapter_ids?.includes(adapter.id);
                return (
                  <button
                    key={adapter.id}
                    onClick={() =>
                      isLinked
                        ? onUnlink(device.id, adapter.id)
                        : onLink(device.id, adapter.id)
                    }
                    className={`px-2 py-1 text-xs rounded border transition-colors ${isLinked
                        ? "bg-blue-100 text-blue-700 border-blue-300 hover:bg-red-100 hover:text-red-700 hover:border-red-300"
                        : "bg-gray-50 text-gray-500 border-gray-200 hover:bg-blue-50 hover:text-blue-600"
                      }`}
                  >
                    {adapter.id} {isLinked ? "✓" : "+"}
                  </button>
                );
              })}
            </div>
          </div>
        </div>
      ))}
    </div>
  );
};

// Helper pour les icones
const getIcon = (name: string) => {
  if (name.includes('temp')) return <Thermometer size={16} className="text-red-500" />;
  if (name.includes('hum')) return <Droplets size={16} className="text-blue-500" />;
  if (name.includes('batt')) return <Battery size={16} className="text-green-500" />;
  return <Activity size={16} />;
}