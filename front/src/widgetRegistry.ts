import type { JSX } from "react";
import type { WidgetType } from "./types/widget";
import { LineChart } from "./components/widgets/LineChart";
import { LineChartIcon, type LucideProps } from "lucide-react";



export const WIDGET_REGISTRY: Record<WidgetType, {
    component: (props: any) => JSX.Element;
    icon: React.ForwardRefExoticComponent<Omit<LucideProps, "ref"> & React.RefAttributes<SVGSVGElement>>
}> = {
    "line-chart": {
        component: LineChart,
        icon: LineChartIcon
    }
}