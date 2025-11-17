import { useState, useEffect } from "react"

export function useApi<T>(provide : () => Promise<T>) {
    const [data, setData] = useState<T>()

    const refresh = () => { provide().then(p => setData(p)) }

    useEffect(refresh, [])

    return {
        data,
        refresh
    }
}