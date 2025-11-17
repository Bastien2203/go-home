
import { type Adapter } from "../types";



const AdapterList = (props: {
    adapters: Adapter[] | undefined,
    startAdapter: (adapter: Adapter) => void,
    restartAdapter: (adapter: Adapter) => void,
    stopAdapter: (adapter: Adapter) => void,
}) => {

    if (!props.adapters) return <div>Loading...</div>;


    return (
        <div className="flex gap-2">
            {
                props.adapters?.map((adapter) => (
                    <div className="p-1 rounded border flex flex-col" key={adapter.id}>
                        <span>Name: {adapter.name}</span>
                        <span>State : <i className={`${adapter.state == "started" ? "bg-green-500" : adapter.state == "need_restart" ? "bg-amber-500" :  "bg-red-500"} w-3 h-3 rounded-full inline-block`}></i></span>

                        <AdapterButton 
                            state={adapter.state} 
                            startAdapter={() => props.startAdapter(adapter)}
                            restartAdapter={() => props.restartAdapter(adapter)}
                            stopAdapter={() => props.stopAdapter(adapter)}
                            />

                    </div>
                ))
            }
            {
                props.adapters?.length == 0 && <>No adapters</>
            }
        </div>
    )
}



const AdapterButton = (
    props: {
        state: string,
        startAdapter: () => void,
        stopAdapter: () => void,
        restartAdapter: () => void,
    }
) => {

    switch (props.state) {
        case "started":
            return <button
                onClick={() => props.stopAdapter()}
                className="mt-4 px-3 py-1.5 cursor-pointer text-sm font-medium bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
            >
                Stop Adapter
            </button>
        case "stopped":
            return <button
                onClick={() => props.startAdapter()}
                className="mt-4 px-3 py-1.5 cursor-pointer text-sm font-medium bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
            >
                Start Adapter
            </button>
        case "need_restart":
            return <button
                onClick={() => props.restartAdapter()}
                className="mt-4 px-3 py-1.5 cursor-pointer text-sm font-medium bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
            >
                Restart Adapter
            </button>
    }

}

export default AdapterList