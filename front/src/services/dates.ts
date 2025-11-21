

export function formatDatetime(datetime: string): string {
    const today = new Date()
    const d = new Date(datetime)

    if (d.toLocaleDateString() == today.toLocaleDateString()) {
        return `Today - ${d.toLocaleTimeString()}`
    } else {
        return d.toLocaleString()
    }    
}