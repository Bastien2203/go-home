import { useState } from "react"
import { useTopic } from "../../hooks/useTopic"
import type { BluetoothDeviceMessage } from "../../types/topics"
import { Plus } from "lucide-react"


export const ScanBluetoothDevices = (props: {
    onConnect: (name: string, address: string) => void
}) => {
    const [bluetoothDevices, setBluetoothDevices] = useState<Record<string, string>>({})

    const onMessage = (msg: BluetoothDeviceMessage) => {
        setBluetoothDevices((prev) => ({
            ...prev,
            [msg.address]: msg.name
        }))
    }

    const { isConnected } = useTopic<BluetoothDeviceMessage>("topic_bluetooth_device", onMessage)

    const devicesList = Object.entries(bluetoothDevices)

    return (
       <table className="min-w-full divide-y divide-gray-200">
        <thead className="bg-gray-50">
            <tr>
                <th scope="col" className="sticky top-0 z-10 bg-gray-50 px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Name
                </th>
                <th scope="col" className="sticky top-0 z-10 bg-gray-50 px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Address
                </th>
                <th scope="col" className="sticky top-0 z-10 bg-gray-50 px-6 py-3">
                    <span className="sr-only">Action</span>
                </th>
            </tr>
        </thead>

        <tbody className="bg-white divide-y divide-gray-200">
            {devicesList.length > 0 ? (
                devicesList.map(([address, name]) => (
                    <tr key={`bt-${address.replace(/:/g, '_')}`} className="h-[3em]">
                        <td className="px-6 py-4 whitespace-nowrap">
                            <div className="text-sm font-medium text-gray-900">
                                {name || "Inconnu"}
                            </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                            <div className="text-sm text-gray-500 font-mono bg-gray-100 inline-block px-2 py-0.5 rounded">
                                {address}
                            </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium flex justify-end">
                            <button className="bg-primary-600 p-1 rounded cursor-pointer hover:opacity-80" title="Create device" onClick={() => props.onConnect(name, address)}>
                                <Plus className="cursor-pointer text-white" size={16}/>
                            </button>
                        </td>
                    </tr>
                ))
            ) : (
                // Gestion de l'état vide propre
                <tr>
                    <td colSpan={3} className="px-6 py-10 text-center text-gray-500">
                        {isConnected
                            ? "Scan en cours... Aucun appareil détecté pour le moment."
                            : "En attente de connexion..."}
                    </td>
                </tr>
            )}
        </tbody>
    </table>


    )
}