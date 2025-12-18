import type { Adapter } from "../../types/adapter";
import type { Device } from "../../types/device";
import { ServerOff, SquareArrowOutUpRight, Trash } from "lucide-react";

interface Props {
  devices: Device[];
  adapters: Adapter[];
  onDeviceSelected: (d: Device) => void;
  onDeviceDelete: (d: Device) => void;
}

export const DeviceList = (props: Props) => {
  if (props.devices.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center h-64 text-center">
        <div className="bg-white p-4 rounded-full shadow-sm mb-4">
          <ServerOff size={32} className="text-gray-300" />
        </div>
        <h3 className="text-lg font-medium text-gray-900">No devices</h3>
        <p className="text-sm text-gray-500 mt-1 max-w-xs">
          Add one by clicking on the "+" button
        </p>
      </div>
    );
  }

  return (
    <table className="divide-y divide-gray-200 min-w-full">
      <thead className="bg-gray-50">
        <tr>
          <th scope="col" className="sticky top-0 z-10 bg-gray-50 px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
            Name
          </th>
          <th scope="col" className="sticky top-0 z-10 bg-gray-50 px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
            Address
          </th>
          <th scope="col" className="sticky top-0 z-10 bg-gray-50 px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
            Protocol
          </th>
          <th scope="col" className="sticky top-0 z-10 bg-gray-50 px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
            
          </th>
        </tr>
      </thead>
      <tbody className="bg-white divide-y divide-gray-200">
        {props.devices.map((device) => (
          <tr key={device.id} className="h-[3em]">
            <td className="px-6 py-4 whitespace-nowrap">
              {device.name}
            </td>
            <td className="px-6 py-4 whitespace-nowrap">
              <div className="text-sm text-gray-500 font-mono bg-gray-100 inline-block px-2 py-0.5 rounded">
                  {device.address}
              </div>
            </td>
            <td className="px-6 py-4 whitespace-nowrap">
              {device.address_type}
            </td>
            <td className="px-6 py-4 whitespace-nowrap flex items-center gap-2 justify-end">
              <div className="bg-primary-600 p-1 rounded cursor-pointer hover:opacity-80" title="Open" onClick={() => props.onDeviceSelected(device)}>
                <SquareArrowOutUpRight className="cursor-pointer text-white" size={16}/>
              </div>

              <div className="bg-red-700 p-1 rounded cursor-pointer hover:opacity-80" title="Delete" onClick={() => props.onDeviceDelete(device)}>
                <Trash className="text-white"  size={16}/>
              </div>
            </td>
          </tr>
        ))}
      </tbody>

    </table>
  );
};

