import { useEffect, useState } from 'react';
import {
    LineChart as RechartsLineChart,
    Line,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    Legend,
    ResponsiveContainer
} from 'recharts';

type DataSet = {
    x_label: string;
    y_label: string;
    points: { x: string | number; y: number }[];
}

export const LineChart = (props: {
    dataUrl: string;
}) => {
    const [data, setData] = useState<DataSet | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(false);

    useEffect(() => {
        setLoading(true);
        fetch(props.dataUrl)
            .then(r => {
                if (!r.ok) throw new Error("Network response was not ok");
                return r.json();
            })
            .then((d: DataSet) => {
                if (d.points.length == 0) {
                    setData(null)
                    setLoading(false)
                    return
                }
                let processedPoints = d.points;
                if (new Date(d.points[0].y).toString() != "Invalid Date") {
                    processedPoints = d.points
                        .map(p => ({
                            ...p,
                            x: new Date(p.x).getTime()
                        }))
                        .sort((a, b) => (a.x as number) - (b.x as number));
                }

                setData({ ...d, points: processedPoints });
                setLoading(false);
            })
            .catch(e => {
                console.error("Fetch error:", e);
                setError(true);
                setLoading(false);
            });
    }, [props.dataUrl]);

    if (loading) return <div>Loading Data...</div>;
    if (error || !data) return <div >No data</div>;

    const xIsDate = new Date(data.points[0].y).toString() != "Invalid Date"

    return (
        <div className="w-full h-[400px]">
            <ResponsiveContainer width="100%" height="100%">
                <RechartsLineChart
                    data={data.points}
                    margin={{
                        top: 5,
                        right: 30,
                        left: 20,
                        bottom: 5,
                    }}
                >
                    <CartesianGrid strokeDasharray="3 3" stroke="#eee" />

                    <XAxis 
                        dataKey="x"
                        type={xIsDate ? "number" : "category"}
                        scale={xIsDate ? "time" : "auto"}
                        domain={xIsDate ? ['dataMin', 'dataMax'] : ['auto', 'auto']}
                        tickFormatter={xIsDate ? formatDateAxis : undefined}
                        label={{ value: data.x_label, position: 'insideBottomRight', offset: -5 }} 
                        stroke="#666"
                    />

                    <YAxis
                        label={{ value: data.y_label, angle: -90, position: 'insideLeft' }}
                        stroke="#666"
                        domain={['dataMin', "auto"]}
                    />

                    <Tooltip
                        contentStyle={{ borderRadius: '8px', border: 'none', boxShadow: '0 2px 8px rgba(0,0,0,0.15)' }}
                        labelFormatter={xIsDate ? formatDateTooltip : undefined}
                    />

                    <Legend />

                    <Line
                        type="stepAfter"
                        dataKey="y"
                        stroke="#8884d8"
                        strokeWidth={3}
                        activeDot={{ r: 8 }}
                        dot={false}
                        name={data.y_label}
                    />
                </RechartsLineChart>
            </ResponsiveContainer>
        </div>
    );
}

const formatDateAxis = (tick: number) => {
    return new Date(tick).toLocaleDateString('fr-FR', {
        day: '2-digit',
        month: '2-digit'
    });
};

const formatDateTooltip = (tick: number) => {
    return new Date(tick).toLocaleString('fr-FR', {
        weekday: 'long',
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });
};