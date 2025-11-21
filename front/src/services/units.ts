import type { Unit } from "../types/units"

export function formatUnit(u?: Unit): string {
    switch (u) {
        case "celsius":
            return "°C"
        case "percent":
            return "%"
        case "volt":
            return "V"
        default:
            return ""
    }
}