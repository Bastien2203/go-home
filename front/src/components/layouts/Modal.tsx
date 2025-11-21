import { useEffect, type PropsWithChildren, type ReactNode } from "react";
import { createPortal } from "react-dom";

const SIZES = {
  sm: "max-w-sm",
  md: "max-w-md",
  lg: "max-w-lg",
  xl: "max-w-xl",
  full: "max-w-full m-4",
};

type ModalProps = {
  isOpen: boolean;
  onClose: () => void;
  title?: string; 
  footer?: ReactNode; 
  size?: keyof typeof SIZES; 
  closeOnOverlayClick?: boolean;
};

export const Modal = ({
  isOpen,
  onClose,
  title,
  footer,
  size = "md",
  closeOnOverlayClick = true,
  children,
}: PropsWithChildren<ModalProps>) => {
  
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape" && isOpen) onClose();
    };
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [isOpen, onClose]);

  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = 'unset';
    }
    return () => { document.body.style.overflow = 'unset'; };
  }, [isOpen]);

  if (!isOpen) return null;

  return createPortal(
    <div 
      className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm transition-opacity"
      onClick={closeOnOverlayClick ? onClose : undefined}
      aria-modal="true"
      role="dialog"
    >
      <div
        className={`bg-white rounded-xl shadow-2xl w-full transform transition-all scale-100 flex flex-col max-h-[90vh] ${SIZES[size]}`}
        onClick={(e) => e.stopPropagation()} 
      >
        {/* --- Header --- */}
        {title && (
          <div className="flex items-center justify-between px-6 py-4 border-b border-gray-100">
            {title && <h3 className="text-lg font-semibold text-gray-900">{title}</h3>}
            
            <button
              onClick={onClose}
              className="p-1 ml-auto text-gray-400 transition-colors rounded-md hover:text-gray-600 hover:bg-gray-100"
              aria-label="Fermer"
            >
              <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        )}

        {/* --- Body --- */}
        <div className="p-6 overflow-y-auto">
          {children}
        </div>

        {/* --- Footer --- */}
        {footer && (
          <div className="px-6 py-4 bg-gray-50 border-t border-gray-100 rounded-b-xl flex justify-end gap-3">
            {footer}
          </div>
        )}
      </div>
    </div>,
    document.body
  );
};