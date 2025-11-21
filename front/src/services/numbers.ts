

export function formatNumber(n: number | Number): string {
    if (n === null || n === undefined || isNaN(Number(n))) {
        return "-";
    }

    const value = Number(n);
   
    const formatter = new Intl.NumberFormat(undefined, {
        maximumFractionDigits: 2, 
        minimumFractionDigits: 0, 
        useGrouping: true         
    });

    return formatter.format(value);
}