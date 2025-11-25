import type { JSX } from "react"
import { ScanBluetoothDevices } from "../components/widgets/ScanBluetoothDevices"
import { StateList } from "../components/widgets/StateList"
import { DeviceList } from "../components/widgets/DeviceList"
import { Activity, Bluetooth, Network, type LucideProps } from "lucide-react"

export type Widget = {
    id: string,
    name: string,
    component: (props: any) => JSX.Element,
    cols: number,
    rows: number,
    padding?: boolean,
    icon: React.ForwardRefExoticComponent<Omit<LucideProps, "ref"> & React.RefAttributes<SVGSVGElement>>,
}

export const WIDGET_REGISTRY: Widget[] = [
    {
        id: "device-list",
        name: "Device List",
        component: DeviceList,
        cols: 4,
        rows: 2,
        icon: Network,

    },
    {
        id: "bluetooth-scan",
        name: "Bluetooth Devices Around",
        component: ScanBluetoothDevices,
        cols: 3,
        rows: 2,
        icon: Bluetooth,
    },
    {
        id: "scanner-list",
        name: "Scanners States",
        component: StateList,
        cols: 1,
        rows: 1,
        icon: Activity,
        padding: true,
    },
    {
        id: "adapter-list",
        name: "Adapters States",
        component: StateList,
        cols: 1,
        rows: 1,
        icon: Activity,
        padding: true,
    },

]