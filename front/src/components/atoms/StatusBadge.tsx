

const StatusBadge = ({ icon, label, count, active = false , noCount}: any) => (
  <div className={`flex items-center gap-2 px-3 py-1.5 rounded-full border text-sm transition-colors ${
    active 
      ? 'bg-green-50 border-green-200 text-green-700 shadow-sm' 
      : 'bg-white border-gray-200 text-gray-600 shadow-sm'
  }`}>
    <span className={active ? "text-green-600" : "text-gray-400"}>{icon}</span>
    <span className="font-medium">{label}</span>
    <span className={`ml-1 px-1.5 py-0.5 rounded text-xs font-bold ${
        active ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-700'
    }`}>
      {
        noCount ? "-" : count
      }
    </span>
  </div>
)

export default StatusBadge