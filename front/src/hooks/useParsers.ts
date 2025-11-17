import { api } from "../services/api";
import { useApi } from "./useApi";


export function useParsers() {
    const h = useApi(api.getParsers)
    return {
        parsers: h.data,
        refresh: h.refresh
    }
}