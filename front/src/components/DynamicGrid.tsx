import { Suspense, type JSX } from "react";
import { Frame } from "./layouts/Frame";
import { type LucideProps } from "lucide-react";


type GridItem = {
    id: string;
    rows: number;
    cols: number;
    padding?: boolean;
    name: string;
    icon: React.ForwardRefExoticComponent<Omit<LucideProps, "ref"> & React.RefAttributes<SVGSVGElement>>;
    component: (props: any) => JSX.Element
}

export const DynamicGrid = (props: {
    widgets: GridItem[]
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
                
                const style: React.CSSProperties = {
                    gridColumn: `span ${widget.cols}`,
                    gridRow: `span ${widget.rows}`,
                };
                return (
                    <div key={widget.id} className="relative group" style={style}>
                        <Frame icon={widget.icon as any} title={widget.name} padding={widget.padding}>
                            <Suspense fallback={<div className="p-4 text-gray-400 text-sm">Chargement...</div>}>
                                <WidgetComponent/>
                            </Suspense>
                        </Frame>
                    </div>
                );
            })}
        </div>
    );

}