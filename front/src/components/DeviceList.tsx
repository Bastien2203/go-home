import type { Adapter, Device } from "../types";



const DeviceList = (props: {
    devices: Device[] | undefined,
    adapters: Adapter[] | undefined
    startDevice: (device: Device) => void
    linkAdapterToDevice: (adapter: Adapter, device: Device) => void
    unlinkAdapterToDevice: (adapter: Adapter, device: Device) => void
}) => {

    if (!props.devices) return <div className="text-gray-500">Loading...</div>;
    if (!props.adapters) return <div className="text-gray-500">Loading...</div>;

    if (props.devices.length === 0)
        return <div className="text-gray-500 italic">No devices</div>;


    return (
        <div className="flex gap-2">
            {
                props.devices?.map((d) => (
                    <div className="p-1 rounded bg-white text-black flex flex-col border" key={d.addr}>
                        <span>Addr : {d.addr}</span>
                        <span>Name : {d.name}</span>
                        <span>Type : {d.type}</span>
                        <span>Parser : {d.parser_type}</span>
                        <span>State : <i className={`${d.running ? "bg-green-500" : "bg-red-500"} w-3 h-3 rounded-full inline-block`}></i></span>

                        <div className="flex gap-2">
                            Adapters :
                            {
                                props.adapters?.map(adapter => {
                                    const checked = adapter.devices.some(item => item.addr == d.addr)
                                    return <span className="flex items-center gap-1" key={adapter.name}>
                                        <input
                                            type="checkbox"
                                            id={`checkbox-${adapter.name}`}
                                            checked={checked}
                                            onChange={() => checked ? props.unlinkAdapterToDevice(adapter, d): props.linkAdapterToDevice(adapter, d)}
                                        />
                                        <label htmlFor={`checkbox-${adapter.name}`}>{adapter.name}</label>
                                    </span>
                                })
                            }
                        </div>
                        {
                            !d.running && <button
                                onClick={() => props.startDevice(d)}
                                className="mt-4 px-3 py-1.5 cursor-pointer text-sm font-medium bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
                            >
                                Start device
                            </button>
                        }
                    </div>
                ))
            }
        </div>
    )
}


export default DeviceList