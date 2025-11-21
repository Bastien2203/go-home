import { Activity, Network, Plus } from "lucide-react";
import { CreateDeviceForm } from "./components/domain/device/CreateDeviceForm";
import { useAdapters, useScanners, useDevices, useProtocols } from "./hooks/useDomain";
import { DeviceList } from "./components/domain/device/DeviceList";
import { Frame } from "./components/layouts/Frame";
import { Header } from "./components/Header";
import { FloatingButton } from "./components/atoms/FloatingButton";
import { useModal } from "./hooks/useModal";
import { Modal } from "./components/layouts/Modal";
import { StateList } from "./components/domain/StateList";



const App = () => {
  const { adapters, adaptersLoading, adaptersError } = useAdapters();
  const { scanners, scannersLoading, scannersError } = useScanners();
  const { devices, createDevice, linkAdapter, unlinkAdapter, devicesLoading, devicesError } = useDevices();
  const { protocols } = useProtocols();
  const { modalIsOpen, openModal, closeModal } = useModal();


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
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 lg:gap-8 items-start">
          <div className="lg:col-span-4 xl:col-span-3 lg:sticky lg:top-6 space-y-6">
            <Frame icon={Activity} title="État des Scanners">
              <StateList objects={scanners} objectName="Scanner" />
            </Frame>

            <Frame icon={Activity} title="État des Adapters">
              <StateList objects={adapters} objectName="Adapter" />
            </Frame>
          </div>

          <div className="lg:col-span-8 xl:col-span-9 min-h-[500px]">
            <Frame icon={Network} title="Appareils connectés">
              <DeviceList
                adapters={adapters}
                devices={devices}
                onLink={linkAdapter}
                onUnlink={unlinkAdapter}
              />
            </Frame>
          </div>
        </div>
      </main>

      <Modal
        isOpen={modalIsOpen}
        onClose={closeModal}
        title="Add new device"
        size="xl"
      >
        <CreateDeviceForm
          onSubmit={(d) => {
            createDevice(d).then(() => {
              closeModal()
            })
          }}
          adapters={adapters}
          protocols={protocols}
        />
      </Modal>

      <FloatingButton onClick={openModal}>
        <Plus />
      </FloatingButton>
    </div>
    </>
  )
}

export default App;

