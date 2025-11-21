import { Activity, Bluetooth, Network, Plus } from "lucide-react";
import { CreateDeviceForm } from "./components/domain/device/CreateDeviceForm";
import { useAdapters, useScanners, useDevices, useProtocols } from "./hooks/useDomain";
import { DeviceList } from "./components/domain/device/DeviceList";
import { Frame } from "./components/layouts/Frame";
import { Header } from "./components/Header";
import { FloatingButton } from "./components/atoms/FloatingButton";
import { useModal } from "./hooks/useModal";
import { Modal } from "./components/layouts/Modal";
import { StateList } from "./components/domain/StateList";
import { ScanBluetoothDevices } from "./components/domain/ScanBluetoothDevices";
import { useState } from "react";



const App = () => {
  const { adapters, adaptersLoading, adaptersError } = useAdapters();
  const { scanners, scannersLoading, scannersError } = useScanners();
  const { devices, createDevice, linkAdapter, unlinkAdapter, devicesLoading, devicesError } = useDevices();
  const { protocols } = useProtocols();
  const { modalIsOpen, openModal, closeModal } = useModal();
  const [selectedBluetoothDeviceMeta, setSelectedBluetoothDeviceMeta] = useState<{ name: string, address: string }>()


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
          <div className="grid grid-cols-5 grid-rows-[200px_200px_1fr] gap-4">
            <div className="col-span-2">
              <Frame icon={Activity} title="État des Scanners" padding>
                <StateList objects={scanners} objectName="Scanner" />
              </Frame>
            </div>


            <div className="col-span-2 col-start-1 row-start-2">
              <Frame icon={Activity} title="État des Adapters" padding>
                <StateList objects={adapters} objectName="Adapter" />
              </Frame>
            </div>

            <div className="col-span-3 row-span-2 col-start-3 row-start-1">
              <Frame icon={Bluetooth} title="Bluetooth Devices Around">
                <ScanBluetoothDevices onConnect={(name: string, address: string) => {
                  setSelectedBluetoothDeviceMeta({ name, address })
                  openModal()
                }} />
              </Frame>
            </div>

            <div className="col-span-5 row-span-2 row-start-3">
              <Frame icon={Network} title="Appareils connectés" padding>
                <DeviceList
                  adapters={adapters}
                  devices={devices}
                  onLink={linkAdapter}
                  onUnlink={unlinkAdapter}
                />
              </Frame>

            </div>
          </div>
        </main >

        <Modal
          isOpen={modalIsOpen}
          onClose={closeModal}
          title="Add new device"
          size="xl"
        >
          <CreateDeviceForm
            defaultData={selectedBluetoothDeviceMeta}
            onSubmit={(d) => {
              createDevice(d).then(() => {
                closeModal()
              })
            }}
            adapters={adapters}
            protocols={protocols}
          />
        </Modal>

        <FloatingButton onClick={() => {
          setSelectedBluetoothDeviceMeta(undefined)
          openModal()
        }}>
          <Plus />
        </FloatingButton>
      </div >
    </>
  )
}

export default App;

