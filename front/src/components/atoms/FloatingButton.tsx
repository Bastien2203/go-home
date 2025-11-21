import type { PropsWithChildren } from "react"

type Props = {
    className?: string
    onClick?: () => void
}

export const FloatingButton = (props: PropsWithChildren<Props>) => (
    <button onClick={props.onClick} className={` fixed bottom-10 right-10 shadow-lg hover:shadow-xl bg-primary-600 text-white rounded-full p-1 aspect-square flex items-center justify-center ${props.className} w-16 h-16`}>
        {props.children}
    </button>
)