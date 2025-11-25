import { useState, useEffect } from "react";
import { WIDGET_REGISTRY } from "../types/widget";

export const useWidgetLayout = () => {
  const [activeWidgetIds, setActiveWidgetIds] = useState<string[]>(() => {
    const saved = localStorage.getItem("widget-grid-layout");
    return saved ? JSON.parse(saved) : WIDGET_REGISTRY.map(w => w.id);
  });

  useEffect(() => {
    localStorage.setItem("widget-grid-layout", JSON.stringify(activeWidgetIds));
  }, [activeWidgetIds]);

  const addWidget = (id: string) => {
    if (!activeWidgetIds.includes(id)) {
      setActiveWidgetIds([...activeWidgetIds, id]);
    }
  };

  const removeWidget = (id: string) => {
    setActiveWidgetIds(activeWidgetIds.filter(wId => wId !== id));
  };

  const activeWidgets = WIDGET_REGISTRY.filter(w => activeWidgetIds.includes(w.id));

  return { activeWidgets, activeWidgetIds, addWidget, removeWidget };
};