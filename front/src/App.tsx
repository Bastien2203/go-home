import AdapterList from "./components/AdapterList"
import CreateDeviceModal from "./components/CreateDeviceModal"
import DeviceList from "./components/DeviceList"
import { useAdapters } from "./hooks/useAdapters"
import { useDevices } from "./hooks/useDevices"




const App = () => {

  const {devices, refreshDevices, startDevice, deviceTypes} = useDevices()
  const {adapters, linkAdapterToDevice,unlinkAdapterToDevice, startAdapter, stopAdapter, restartAdapter} = useAdapters()

  return (
    <div className="w-full h-screen p-4">
      <h1>Go Home</h1>


  <div className="flex flex-col gap-4">
      <section>
        <h2>Adapters</h2>
        <AdapterList adapters={adapters} startAdapter={startAdapter} stopAdapter={stopAdapter} restartAdapter={restartAdapter}/>
      </section>


      <section>
        <h2>Devices</h2>
        <DeviceList devices={devices} startDevice={startDevice} adapters={adapters} linkAdapterToDevice={linkAdapterToDevice} unlinkAdapterToDevice={unlinkAdapterToDevice}/>
      </section>

    <CreateDeviceModal refreshDeviceList={refreshDevices} deviceTypes={deviceTypes}/>


      </div>
    </div>
  )
}



export default App
