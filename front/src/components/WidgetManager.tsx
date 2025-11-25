import { Check, Plus } from "lucide-react";
import { WIDGET_REGISTRY } from "../types/widget";

export const WidgetManager = (props: {
    activeIds: string[],
    onToggle: (id: string) => void
}) => {
    return (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            {WIDGET_REGISTRY.map(widget => {
                const isActive = props.activeIds.includes(widget.id);
                const Icon = widget.icon;
                
                return (
                    <button
                        key={widget.id}
                        onClick={() => props.onToggle(widget.id)}
                        className={`
                            flex items-center justify-between p-4 rounded-lg border transition-all
                            ${isActive 
                                ? "border-primary-500 bg-primary-50 text-primary-900" 
                                : "border-gray-200 hover:border-gray-300 bg-white text-gray-500"}
                        `}
                    >
                        <div className="flex items-center gap-3">
                            <div className={`p-2 rounded-md ${isActive ? "bg-primary-100" : "bg-gray-100"}`}>
                                <Icon size={20} />
                            </div>
                            <span className="font-medium">{widget.name}</span>
                        </div>
                        {isActive ? <Check size={18} /> : <Plus size={18} />}
                    </button>
                )
            })}
        </div>
    )
}