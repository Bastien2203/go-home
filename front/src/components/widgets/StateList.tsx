import { Play, Square } from "lucide-react";
import type { State } from "../../types/states"
import { BadgeRestarting, BadgeRunning, BadgeStopped } from "../atoms/Badges"


type StateObject = {
    id: string;
    state: State;
    name: string;
}

export const StateList = (props: {
    objects: StateObject[];
    objectName: string;
    start: (id: string) => void;
    stop: (id: string) => void;
}) => {
    return <ul className="space-y-3">
        {props.objects.length === 0 && (
            <li className="text-gray-400 text-sm italic text-center py-2">
                Aucun {props.objectName} détécté
            </li>
        )}
        {props.objects.map(o => (
            <li key={o.id} className="flex items-center justify-between text-sm group">
                <span className="text-gray-600 font-medium flex items-center gap-2">
                    {o.name}
                    {
                        o.state === "stopped" ?
                            <Play size={16} className="cursor-pointer" onClick={() => props.start(o.id)}/> :
                            o.state === "restarting" ?
                                <></> :
                                <Square size={16} className="cursor-pointer" onClick={() => props.stop(o.id)}/>
                    }
                    
                </span>
                {
                    o.state === "running" ?
                        <BadgeRunning /> :
                        o.state === "restarting" ?
                            <BadgeRestarting /> :
                            <BadgeStopped />
                }
            </li>
        ))}
    </ul>
}


