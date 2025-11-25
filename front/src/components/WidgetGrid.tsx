import { Suspense } from "react";
import type { Widget } from "../types/widget"
import { Frame } from "./layouts/Frame";
import { X } from "lucide-react";


export const WidgetGrid = (props: {
    widgets: Widget[];
    propsMap: Record<string, any>;
    onRemove?: (id: string) => void;
    isEditing?: boolean;
}) => {
    const TOTAL_COLS = 4;
    const ROW_HEIGHT = "180px";

    return (
        <div
            style={{
                display: "grid",
                gridTemplateColumns: `repeat(${TOTAL_COLS}, 1fr)`,
                gap: "1em",
                gridAutoRows: ROW_HEIGHT,
                gridAutoFlow: "dense",
            }}
        >
            {props.widgets.map((widget) => {
                const WidgetComponent = widget.component;
                const dynamicProps = props.propsMap[widget.id] || {};
                
                const style: React.CSSProperties = {
                    gridColumn: `span ${widget.cols}`,
                    gridRow: `span ${widget.rows}`,
                };
                return (
                    <div key={widget.id} className="relative group" style={style}>
                        {props.isEditing && (
                            <button 
                                onClick={() => props.onRemove?.(widget.id)}
                                className="absolute -top-2 -right-2 z-50 bg-red-500 text-white p-1 rounded-full shadow-md hover:bg-red-600 transition-transform hover:scale-110"
                            >
                                <X size={14} />
                            </button>
                        )}
                        <Frame icon={widget.icon as any} title={widget.name} padding={widget.padding} className={props.isEditing ? "border border-dashed border-primary-300" : ""}>
                            <Suspense fallback={<div className="p-4 text-gray-400 text-sm">Chargement...</div>}>
                                <WidgetComponent
                                    {...dynamicProps}
                                />
                            </Suspense>
                        </Frame>
                    </div>
                );
            })}
        </div>
    );

}