import { Cpu, Server, Radio, Database, Activity } from "lucide-react";
import { CreateDeviceForm } from "./components/CreateDeviceForm";
import { useAdapters, useScanners, useDevices, useProtocols } from "./hooks/useDomain";
import { DeviceList } from "./components/DeviceList";

// Composant Badge amélioré
const StatusBadge = ({ icon, label, count, active = false }: any) => (
  <div className={`flex items-center gap-2 px-3 py-1.5 rounded-full border text-sm transition-colors ${
    active 
      ? 'bg-green-50 border-green-200 text-green-700 shadow-sm' 
      : 'bg-white border-gray-200 text-gray-600 shadow-sm'
  }`}>
    <span className={active ? "text-green-600" : "text-gray-400"}>{icon}</span>
    <span className="font-medium">{label}</span>
    <span className={`ml-1 px-1.5 py-0.5 rounded text-xs font-bold ${
        active ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-700'
    }`}>
      {count}
    </span>
  </div>
)

const App = () => {
  const { adapters } = useAdapters();
  const { scanners } = useScanners();
  const { devices, createDevice, linkAdapter, unlinkAdapter } = useDevices();
  const { protocols } = useProtocols();

  // Calcul si au moins un scanner tourne
  const isSystemScanning = scanners.some(s => s.is_running);

  return (
    <div className="min-h-screen bg-gray-50 text-gray-800 font-sans pb-10">
      
      {/* Header avec ombre légère et meilleure structure */}
      <header className="sticky top-0 z-50 bg-white/80 backdrop-blur-md border-b border-gray-200 px-6 py-4">
        <div className="container mx-auto flex flex-col md:flex-row items-center justify-between gap-4">
          
          {/* Logo Area */}
          <div className="flex items-center gap-3">
            <div className="bg-gradient-to-br from-blue-600 to-blue-700 p-2.5 rounded-xl text-white shadow-lg shadow-blue-600/20">
              <Cpu size={24} strokeWidth={2.5} />
            </div>
            <div>
              <h1 className="text-xl font-bold tracking-tight text-gray-900 leading-tight">
                GoHome <span className="text-blue-600">Control</span>
              </h1>
              <p className="text-xs text-gray-500 font-medium">Kernel v1.0.0</p>
            </div>
          </div>

          {/* Status Bar */}
          <div className="flex flex-wrap justify-center gap-3">
            <StatusBadge icon={<Server size={16} />} label="Adapters" count={adapters.length} />
            <StatusBadge icon={<Radio size={16} />} label="Scanners" count={scanners.length} active={isSystemScanning} />
            <StatusBadge icon={<Database size={16} />} label="Appareils" count={devices.length} />
          </div>
        </div>
      </header>

      {/* Main Layout : Grid 12 colonnes pour plus de précision */}
      <main className="container mx-auto p-6 grid grid-cols-1 lg:grid-cols-12 gap-8">

        {/* Colonne Gauche : Sidebar (4/12) */}
        <div className="lg:col-span-4 xl:col-span-3 space-y-6">
          
          {/* Wrapper Sticky pour que le formulaire suive le scroll */}
          <div className="sticky top-28 space-y-6">
            
            {/* Formulaire de création */}
            <CreateDeviceForm 
              onSubmit={createDevice} 
              adapters={adapters} 
              protocols={protocols} 
            />

            {/* État du système (Design Carte) */}
            <div className="bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden">
              <div className="bg-gray-50 px-4 py-3 border-b border-gray-100 flex items-center gap-2">
                <Activity size={16} className="text-blue-600" />
                <h3 className="font-semibold text-sm text-gray-700">État des Scanners</h3>
              </div>
              
              <div className="p-4">
                <ul className="space-y-3">
                  {scanners.length === 0 && (
                    <li className="text-gray-400 text-sm italic text-center py-2">
                      Aucun scanner détecté
                    </li>
                  )}
                  {scanners.map(s => (
                    <li key={s.id} className="flex items-center justify-between text-sm group">
                      <span className="text-gray-600 font-medium">{s.name || "Scanner inconnu"}</span>
                      <div className={`flex items-center gap-2 px-2 py-1 rounded-md border transition-all ${
                        s.is_running 
                          ? "bg-green-50 border-green-200 text-green-700" 
                          : "bg-red-50 border-red-200 text-red-700"
                      }`}>
                        <span className={`relative flex h-2 w-2`}>
                          {s.is_running && (
                            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
                          )}
                          <span className={`relative inline-flex rounded-full h-2 w-2 ${s.is_running ? 'bg-green-500' : 'bg-red-500'}`}></span>
                        </span>
                        <span className="text-xs font-bold tracking-wide">
                          {s.is_running ? "ACTIF" : "STOP"}
                        </span>
                      </div>
                    </li>
                  ))}
                </ul>
              </div>
            </div>

          </div>
        </div>

        {/* Colonne Droite : Liste Principale (8/12) */}
        <div className="lg:col-span-8 xl:col-span-9">
            {/* Titre de section optionnel */}
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-gray-800">Appareils Connectés</h2>
            </div>

            <DeviceList 
              devices={devices} 
              onLink={linkAdapter} 
              onUnlink={unlinkAdapter}
            />
            
            {devices.length === 0 && (
               <div className="text-center py-12 bg-white rounded-xl border border-dashed border-gray-300">
                  <p className="text-gray-500">Aucun appareil configuré.</p>
               </div>
            )}
        </div>

      </main>
    </div>
  )
}

export default App;