import { Link2, Link2Off, Plus } from "lucide-react";
import type { Adapter } from "../../../types/adapter";


type Props = {
  adapter: Adapter;
  isLinked: boolean;
  onToggle: () => void;
};

export const AdapterToggle = ({ adapter, isLinked, onToggle }: Props) => {
  return (
    <button
      onClick={(e) => {
        e.stopPropagation(); 
        onToggle();
      }}
      className={`
        group relative flex items-center gap-2 px-3 py-1.5 text-xs font-medium rounded-full border transition-all duration-200
        ${
          isLinked
            ? "bg-primary-50 border-primary-200 text-primary-700 hover:bg-red-50 hover:border-red-200 hover:text-red-600 hover:shadow-sm"
            : "bg-transparent border-dashed border-gray-300 text-gray-400 hover:border-primary-400 hover:text-primary-600 hover:bg-white"
        }
      `}
      title={isLinked ? "Cliquer pour délier" : "Cliquer pour lier"}
    >
      <div className="relative flex items-center justify-center w-3 h-3">
        {isLinked ? (
          <>
            <Link2
              size={12}
              className="absolute transition-opacity duration-200 group-hover:opacity-0"
            />
            <Link2Off
              size={12}
              className="absolute opacity-0 transition-opacity duration-200 group-hover:opacity-100"
            />
          </>
        ) : (
          <Plus size={12} strokeWidth={3} />
        )}
      </div>

      <span className="truncate max-w-[100px]">{adapter.name}</span>
    </button>
  );
};