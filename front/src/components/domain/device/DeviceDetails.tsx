import { formatDatetime } from "../../../services/dates";
import type { Adapter } from "../../../types/adapter";
import type { Device } from "../../../types/device";
import { AdapterToggle } from "./AdapterToggle";
import {
    Cpu,
    Hash,
    Radio,
    CalendarDays,
    Clock,
    ChevronsLeftRightEllipsis,
    ChartLine,
    ArrowDownUp,
} from "lucide-react";
import { DeviceCapabilitiesGrid } from "./DeviceCapabilitiesGrid";
import { Frame } from "../../layouts/Frame";


interface Props {
    adapters: Adapter[];
    device: Device;
    onLink: (dId: string, aId: string) => void;
    onUnlink: (dId: string, aId: string) => void;
}

export const DeviceDetails = ({ adapters, device, onLink, onUnlink }: Props) => {
    return (
        <div className="space-y-6">
            <div className="grid grid-cols-2 grid-rows-2 gap-6">

                {/* Adapters */}
                <Frame icon={Radio} title="Adapters" padding>
                    <div className="flex gap-3 flex-wrap">
                        {adapters.map((adapter) => {
                            const isLinked = device.adapter_ids?.includes(adapter.id);
                            return (
                                <div key={adapter.id} className="flex items-center justify-between p-2 rounded-lg hover:bg-gray-50 border border-transparent hover:border-gray-100 transition-all">

                                    <AdapterToggle
                                        adapter={adapter}
                                        isLinked={!!isLinked}
                                        onToggle={() => isLinked ? onUnlink(device.id, adapter.id) : onLink(device.id, adapter.id)}
                                    />
                                </div>
                            );
                        })}
                    </div>
                </Frame>

                {/* Technical informations */}
                <Frame icon={Cpu} title="Technical Details" padding className="col-start-1 row-start-2">
                    <div className="space-y-2">
                        <DetailRow icon={<Hash size={16} />} label="Id" value={device.id} />
                        <DetailRow icon={<ChevronsLeftRightEllipsis size={16} />} label="Adress" value={device.address} />
                        <DetailRow icon={<ArrowDownUp size={16} />} label="Protocol" value={device.protocol} />
                        <DetailRow icon={<Radio size={16} />} label="Address Type" value={device.address_type} />
                        <DetailRow icon={<CalendarDays size={16} />} label="Created" value={formatDatetime(device.created_at)} />
                        <DetailRow icon={<Clock size={16} />} label="Last update" value={formatDatetime(device.last_updated)} />
                    </div>
                </Frame>

                {/* Capabilities */}
                <Frame icon={ChartLine} title="Device Data" padding className="row-span-2 col-start-2 row-start-1">
                    <DeviceCapabilitiesGrid device={device} />
                </Frame>
            </div>
        </div>
    );
};


const DetailRow = ({ icon, label, value }: { icon: React.ReactNode, label: string, value: string }) => (
    <div className="flex items-center justify-between text-sm">
        <div className="flex items-center text-gray-500 gap-1">
            {icon}
            <span>{label}</span>
        </div>
        <span className="font-mono text-gray-900 bg-gray-100 px-2 py-0.5 rounded text-xs whitespace-nowrap select-all">
            {value}
        </span>
    </div>
);

export default DeviceDetails;