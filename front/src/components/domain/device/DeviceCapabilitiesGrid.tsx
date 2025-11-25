import { Activity } from "lucide-react";
import { getIcon } from "../../../services/icons";
import { formatNumber } from "../../../services/numbers";
import { formatUnit } from "../../../services/units";
import type { Device } from "../../../types/device";


export const DeviceCapabilitiesGrid = ({ device }: { device: Device }) => {
  const capabilities = Object.values(device.capabilities || {});
  const hasData = capabilities.length > 0;

  if (!hasData) {
    return (
      <div className="h-full flex flex-col items-center justify-center text-gray-400 min-h-[200px]">
        <div className="bg-white p-4 rounded-full shadow-sm mb-3">
             <Activity size={24} className="opacity-40" />
        </div>
        <span className="text-sm font-medium">No data received for now</span>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-2 sm:grid-cols-3 gap-4">
      {capabilities.map((cap) => (
        <div 
            key={cap.name} 
            className="group bg-white p-4 rounded-xl border border-gray-200 shadow-sm hover:shadow-md hover:border-blue-200 transition-all duration-200 flex flex-col justify-between"
        >
          <div className="flex items-start justify-between mb-2">
            <span className="text-gray-400 group-hover:text-blue-500 transition-colors bg-gray-50 p-2 rounded-lg group-hover:bg-blue-50">
              {getIcon(cap.name)}
            </span>
          </div>
          
          <div>
              <span className="text-xs font-medium text-gray-500 uppercase tracking-wide block mb-1">
                {cap.name}
              </span>
              <div className="flex items-baseline gap-1">
                <span className="text-2xl font-bold text-gray-900 tabular-nums tracking-tight">
                    {typeof cap.value === "number" ? formatNumber(cap.value) : cap.value}
                </span>
                <span className="text-sm text-gray-400 font-medium">
                    {formatUnit(cap.unit)}
                </span>
              </div>
          </div>
        </div>
      ))}
    </div>
  );
};