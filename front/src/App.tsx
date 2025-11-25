import { Plus, Settings } from "lucide-react";
import { CreateDeviceForm } from "./components/domain/device/CreateDeviceForm";
import { useAdapters, useScanners, useDevices, useProtocols } from "./hooks/useDomain";

import { Header } from "./components/Header";
import { FloatingButton } from "./components/atoms/FloatingButton";
import { Modal } from "./components/layouts/Modal";
import { useState } from "react";
import { WidgetGrid } from "./components/WidgetGrid";
import { useWidgetLayout } from "./hooks/useWidgetLayout";
import { WidgetManager } from "./components/WidgetManager";
import DeviceDetails from "./components/domain/device/DeviceDetails";
import type { Device } from "./types/device";


const App = () => {
  const { adapters, adaptersLoading, adaptersError } = useAdapters();
  const { scanners, scannersLoading, scannersError, startScanner, stopScanner } = useScanners();
  const { devices, createDevice, linkAdapter, unlinkAdapter, devicesLoading, devicesError, deleteDevice } = useDevices();
  const { protocols } = useProtocols();
  const { activeWidgets, activeWidgetIds, addWidget, removeWidget } = useWidgetLayout();

  const [selectedBluetoothDeviceMeta, setSelectedBluetoothDeviceMeta] = useState<{ name: string, address: string }>()
  const [isManagerOpen, setManagerOpen] = useState(false);
  const [isDeviceCreationModalOpen, setDeviceCreationModalOpen] = useState(false);
  const [selectedDevice, setSelectedDevice] = useState<Device | null>(null);
  const [isEditing, setEditing] = useState(false);


  const widgetPropsMap = {
    "bluetooth-scan": {
      onConnect: (name: string, address: string) => {
        setSelectedBluetoothDeviceMeta({ name, address });
        setDeviceCreationModalOpen(true);
      }
    },
    "scanner-list": {
      objects: scanners,
      objectName: "Scanner",
      start: startScanner,
      stop: stopScanner,
      isLoading: scannersLoading
    },
    "adapter-list": {
      objects: adapters,
      objectName: "Adapter",
      start: () => { },
      stop: () => { },
      isLoading: adaptersLoading
    },
    "device-list": {
      devices: devices,
      adapters: adapters,
      onDeviceDelete: (d: Device) => { deleteDevice(d.id) },
      onDeviceSelected: (d: Device) => { setSelectedDevice(d) }
    }
  };

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

        <div className="flex justify-end px-4 sm:px-6 lg:px-8 pt-4 gap-2">
          <button
            onClick={() => setEditing(!isEditing)}
            className={`text-sm px-3 py-1 rounded border ${isEditing ? 'bg-yellow-100 border-yellow-300 text-yellow-800' : 'bg-white'}`}
          >
            {isEditing ? 'Back to view' : 'Organize'}
          </button>

          <button onClick={() => setManagerOpen(true)} className="text-gray-500 hover:text-primary-600">
            <Settings />
          </button>
        </div>

        <main className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 py-8">
          <WidgetGrid
            widgets={activeWidgets}
            propsMap={widgetPropsMap}
            isEditing={isEditing}
            onRemove={removeWidget}
          />
        </main >


        <Modal
          isOpen={isManagerOpen}
          onClose={() => setManagerOpen(false)}
          title="Manage widgets"
          size="lg"
        >
          <WidgetManager
            activeIds={activeWidgetIds}
            onToggle={(id) => {
              if (activeWidgetIds.includes(id)) removeWidget(id);
              else addWidget(id);
            }}
          />
        </Modal>

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
            protocols={protocols}
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

