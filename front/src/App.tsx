import { Activity, Bluetooth, Network, Plus } from "lucide-react";
import { CreateDeviceForm } from "./components/domain/device/CreateDeviceForm";
import { useAdapters, useScanners, useDevices } from "./hooks/useDomain";

import { Header } from "./components/Header";
import { FloatingButton } from "./components/atoms/FloatingButton";
import { Modal } from "./components/layouts/Modal";
import { useState } from "react";
import { DynamicGrid } from "./components/DynamicGrid";
import DeviceDetails from "./components/domain/device/DeviceDetails";
import type { Device } from "./types/device";
import { DeviceList } from "./components/widgets/DeviceList";
import { ScanBluetoothDevices } from "./components/widgets/ScanBluetoothDevices";
import { StateList } from "./components/widgets/StateList";


const App = () => {
  const { adapters, adaptersLoading, adaptersError, startAdapter, stopAdapter } = useAdapters();
  const { scanners, scannersLoading, scannersError, startScanner, stopScanner } = useScanners();
  const { devices, createDevice, linkAdapter, unlinkAdapter, devicesLoading, devicesError, deleteDevice } = useDevices();


  const [selectedBluetoothDeviceMeta, setSelectedBluetoothDeviceMeta] = useState<{ name: string, address: string }>()
  const [isDeviceCreationModalOpen, setDeviceCreationModalOpen] = useState(false);
  const [selectedDevice, setSelectedDevice] = useState<Device | null>(null);


  return (
    <>
      <div className="min-h-screen bg-gray-50/50 text-gray-800 font-sans selection:bg-primary-100 selection:text-primary-900">
        <Header
          adapters={adapters}
          scanners={scanners}
          devices={devices}
          adaptersUnavailable={adaptersLoading || adaptersError != undefined}
          devicesUnavailable={devicesLoading || devicesError != undefined}
          scannersUnavailable={scannersLoading || scannersError != undefined}
        />

        <main className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 py-8">
          <DynamicGrid
            widgets={[
              {
                id: "device-list",
                name: "Device List",
                component: () => <DeviceList devices={devices} adapters={adapters} onDeviceSelected={setSelectedDevice} onDeviceDelete={(d: Device) => deleteDevice(d.id)}/>,
                cols: 4,
                rows: 2,
                icon: Network,

              },
              {
                id: "bluetooth-scan",
                name: "Bluetooth Devices Around",
                component: () => <ScanBluetoothDevices onConnect={(name: string, address: string) => {setSelectedBluetoothDeviceMeta({ name, address }); setDeviceCreationModalOpen(true);}} />,
                cols: 3,
                rows: 2,
                icon: Bluetooth,
              },
              {
                id: "scanner-list",
                name: "Scanners States",
                component: () => <StateList objects={scanners} objectName={"Scanner"} start={startScanner} stop={stopScanner} />,
                cols: 1,
                rows: 1,
                icon: Activity,
                padding: true,
              },
              {
                id: "adapter-list",
                name: "Adapters States",
                component: () => <StateList objects={adapters} objectName={"Adapter"} start={startAdapter} stop={stopAdapter} />,
                cols: 1,
                rows: 1,
                icon: Activity,
                padding: true,
              },
            ]}
          />
        </main >

        <Modal
          isOpen={isDeviceCreationModalOpen}
          onClose={() => setDeviceCreationModalOpen(false)}
          title="Add new device"
          size="xl"
        >
          <CreateDeviceForm
            defaultData={selectedBluetoothDeviceMeta}
            onSubmit={(d) => {
              createDevice(d).then(() => {
                setDeviceCreationModalOpen(false)
              })
            }}
            adapters={adapters}
          />
        </Modal>

        <Modal
          isOpen={selectedDevice != null}
          onClose={() => setSelectedDevice(null)}
          title={selectedDevice?.name ?? "Device Details"}
          size="full"
        >
          {
            selectedDevice && <DeviceDetails
              onLink={linkAdapter}
              onUnlink={unlinkAdapter}
              device={selectedDevice}
              adapters={adapters}
            />
          }
        </Modal>

        <FloatingButton onClick={() => {
          setSelectedBluetoothDeviceMeta(undefined)
          setDeviceCreationModalOpen(true)
        }}>
          <Plus />
        </FloatingButton>
      </div >
    </>
  )
}

export default App;

