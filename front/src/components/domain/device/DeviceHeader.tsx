import { Cpu, Hash } from "lucide-react";
import type { Device } from "../../../types/device";
import { CopyableValue } from "../../atoms/CopyableValue";

const getProtocolStyle = (protocol: string) => {
  switch (protocol) {
    case "zigbee": return "bg-orange-50 text-orange-600 border-orange-100";
    case "wifi": return "bg-blue-50 text-blue-600 border-blue-100";
    case "bluetooth": return "bg-indigo-50 text-indigo-600 border-indigo-100";
    default: return "bg-gray-100 text-gray-600 border-gray-200";
  }
};

export const DeviceHeader = ({ device }: { device: Device }) => {
  return (
    <div className="p-5 pb-3 border-b border-gray-100">
      <div className="flex justify-between items-start">
        <div className="flex items-center gap-3">
          <div className="p-2.5 bg-white border border-gray-100 shadow-sm text-primary-600 rounded-xl">
            <Cpu size={20} />
          </div>
          
          <div>
            <h3 className="font-bold text-gray-900 leading-tight text-sm lg:text-base">
              {device.name}
            </h3>
            <div className="flex items-center gap-2 mt-1.5">

              <span className={`px-2 py-0.5 rounded-sm text-[10px] font-bold uppercase tracking-wider border ${getProtocolStyle(device.protocol)}`}>
                {device.protocol}
              </span>
              {device.address && (
                 <span className="hidden sm:flex text-xs text-gray-400 items-center gap-0.5" title={device.address}>
                    {device.address}
                 </span>
              )}
            </div>
          </div>
        </div>

        <CopyableValue value={device.id}>
           <Hash size={10} className="text-gray-400" />
           <span className="font-mono text-gray-500 mr-1">
             {device.id.slice(0, 8)}...
           </span>
        </CopyableValue>
      </div>
    </div>
  );
};