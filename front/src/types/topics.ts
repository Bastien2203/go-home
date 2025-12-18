
export type Topic = "topic_bluetooth_device"

export type BluetoothDeviceMessage = {
    name: string;
    address: string;
    protocols: string[];
}