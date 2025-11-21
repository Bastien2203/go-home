import type { State } from "../../types/states"
import { BadgeRestarting, BadgeRunning, BadgeStopped } from "../atoms/Badges"


type StateObject = {
    id: string;
    state: State;
    name: string;
}

export const StateList = (props: {
    objects: StateObject[];
    objectName: string
}) => {
    return <ul className="space-y-3">
        {props.objects.length === 0 && (
            <li className="text-gray-400 text-sm italic text-center py-2">
                Aucun {props.objectName} détécté
            </li>
        )}
        {props.objects.map(o => (
            <li key={o.id} className="flex items-center justify-between text-sm group">
                <span className="text-gray-600 font-medium">{o.name}</span>
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


