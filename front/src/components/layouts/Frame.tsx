import type { LucideProps } from "lucide-react";
import type { PropsWithChildren } from "react";


type Props = {
    icon: React.ForwardRefExoticComponent<Omit<LucideProps, "ref"> & React.RefAttributes<SVGSVGElement>>
    title: string;
}

export const Frame = (props: PropsWithChildren<Props>) => (
    <div className="bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden h-full w-full">
        <div className="bg-gray-50 px-4 py-3 border-b border-gray-100 flex items-center gap-2">
            <props.icon size={16} className="text-primary-600" />
            <h3 className="font-semibold text-sm text-gray-700">{props.title}</h3>
        </div>
        <div className="p-4">
            {props.children}
        </div>
    </div>
)