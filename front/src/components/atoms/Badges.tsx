export const BadgeStopped = () => (
    <div className={`flex items-center gap-2 px-2 py-1 rounded-md border transition-all bg-red-50 border-red-200 text-red-700 min-w-20  justify-center`}>
        <span className={`relative flex h-2 w-2`}>
            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-red-400 opacity-75"></span>
            <span className={`relative inline-flex rounded-full h-2 w-2 bg-red-500`}></span>
        </span>
        <span className="text-xs font-bold tracking-wide">
            STOP
        </span>
    </div>
)


export const BadgeRunning = () => (
    <div className={`flex items-center gap-2 px-2 py-1 rounded-md border transition-all bg-green-50 border-green-200 text-green-700 min-w-20  justify-center`}>
        <span className={`relative flex h-2 w-2`}>
            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
            <span className={`relative inline-flex rounded-full h-2 w-2 bg-green-500`}></span>
        </span>
        <span className="text-xs font-bold tracking-wide">
            ACTIF
        </span>
    </div>
)


export const BadgeRestarting = () => (
<div className={`flex items-center gap-2 px-2 py-1 rounded-md border transition-all bg-yellow-50 border-yellow-200 text-yellow-700 min-w-20  justify-center`}>
        <span className={`relative flex h-2 w-2`}>
            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-yellow-400 opacity-75"></span>
            <span className={`relative inline-flex rounded-full h-2 w-2 bg-yellow-500`}></span>
        </span>
        <span className="text-xs font-bold tracking-wide">
            RESTART
        </span>
    </div>
)