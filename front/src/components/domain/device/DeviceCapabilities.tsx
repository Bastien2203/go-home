import { Activity } from "lucide-react";
import { getIcon } from "../../../services/icons";
import { formatNumber } from "../../../services/numbers";
import { formatUnit } from "../../../services/units";
import type { Device } from "../../../types/device";


export const DeviceCapabilities = ({ device }: { device: Device }) => {
  const capabilities = Object.values(device.capabilities || {});
  const hasData = capabilities.length > 0;

  return (
    <div className="p-5 flex-1 bg-white">
      {hasData ? (
        <div className="space-y-3">
          {capabilities.map((cap) => (
            <div key={cap.name} className="flex items-center justify-between group/row">
              <div className="flex items-center text-sm text-gray-500">
                <span className="mr-2.5 text-gray-300 group-hover/row:text-primary-500 transition-colors duration-300">
                  {getIcon(cap.name)}
                </span>
                <span className="capitalize font-medium text-gray-600">{cap.name}</span>
              </div>
              <div className="text-sm font-semibold text-gray-900 tabular-nums">
                {typeof cap.value === "number"
                  ? formatNumber(cap.value)
                  : cap.value}
                <span className="ml-1 text-xs text-gray-400 font-normal">
                  {formatUnit(cap.unit)}
                </span>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="flex flex-col items-center justify-center py-6 text-gray-400 bg-gray-50/50 rounded-lg border border-dashed border-gray-200">
          <Activity size={20} className="mb-2 opacity-40" />
          <span className="text-xs font-medium">Aucune donnée reçue</span>
        </div>
      )}
    </div>
  );
};