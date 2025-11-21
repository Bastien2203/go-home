import { useState, type PropsWithChildren } from "react";
import { Check, Copy } from "lucide-react";

type Props = {
  value: string;
  className?: string;
  successMessage?: string; // Optionnel : Texte à afficher quand copié
};

export const CopyableValue = ({ 
  value, 
  className = "", 
  successMessage = "Copié !",
  children 
}: PropsWithChildren<Props>) => {
  const [isCopied, setIsCopied] = useState(false);

  const handleCopy = async (e: React.MouseEvent) => {
    e.stopPropagation();
    e.preventDefault();

    try {
      await navigator.clipboard.writeText(value);
      setIsCopied(true);
      
      setTimeout(() => setIsCopied(false), 2000);
    } catch (err) {
      console.error("Failed to copy:", err);
    }
  };

  return (
    <button
      type="button"
      onClick={handleCopy}
      title="Cliquez pour copier"
      className={`
        relative group flex items-center gap-1.5 px-1.5 py-0.5 rounded text-[10px] font-mono border transition-all duration-200 cursor-pointer select-none
        ${isCopied 
          ? "bg-green-50 text-green-600 border-green-200 shadow-sm" 
          : `text-gray-400 border-transparent hover:text-gray-600  ${className}`
        }
      `}
    >
      
      {isCopied ? (
        <>
          <Check size={12} className="animate-in zoom-in duration-200" />
          <span className="font-sans font-semibold">{successMessage}</span>
        </>
      ) : (
        <>
           {children}
           
           <Copy size={10} className="opacity-0 group-hover:opacity-100 transition-opacity absolute right-1" />
        </>
      )}
    </button>
  );
};