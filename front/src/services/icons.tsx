import { Thermometer, Droplets, Battery, Activity } from "lucide-react";


export const getIcon = (name: string) => {
  if (name.includes('temp')) return <Thermometer size={16} className="text-red-500" />;
  if (name.includes('hum')) return <Droplets size={16} className="text-blue-500" />;
  if (name.includes('batt')) return <Battery size={16} className="text-green-500" />;
  return <Activity size={16} />;
}