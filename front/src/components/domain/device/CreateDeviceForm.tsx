import React, { useState } from "react";
import type { Adapter } from "../../../types/adapter";
import type { DeviceCreateRequest } from "../../../types/device";

interface Props {
  onSubmit: (req: DeviceCreateRequest) => void;
  adapters: Adapter[];
  defaultData?: {
    name?: string;
    address?: string;
  }
}

export const CreateDeviceForm: React.FC<Props> = (props: Props) => {

  const [formData, setFormData] = useState({
    name: props.defaultData?.name ?? "",
    address: props.defaultData?.address ?? "",
    address_type: "ble",
    adapter_ids: [] as string[],
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.address_type) return alert("Please select an address type");
    props.onSubmit(formData);
    setFormData({ ...formData, name: "", address: "" }); // Reset partiel
  };

  const toggleAdapter = (id: string) => {
    setFormData((prev) => {
      const exists = prev.adapter_ids.includes(id);
      return {
        ...prev,
        adapter_ids: exists
          ? prev.adapter_ids.filter((a) => a !== id)
          : [...prev.adapter_ids, id],
      };
    });
  };

  return (
    <>
      <form onSubmit={handleSubmit} className="flex flex-col gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700">Name</label>
          <input
            type="text"
            required
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm p-2 border"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            placeholder="Living Room Thermometer"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700">
            Address <i className="text-gray-400 font-light">(depend of protocol)</i>
          </label>
          <input
            type="text"
            required
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm p-2 border"
            value={formData.address}
            onChange={(e) => setFormData({ ...formData, address: e.target.value })}
            placeholder="A4:C1:38:..."
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700">Address Type</label>
          <select
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm p-2 border"
            value={formData.address_type}
            onChange={(e) => setFormData({ ...formData, address_type: e.target.value })}
          >
              <option key="ble" value="ble">
                Bluetooth
              </option>
          </select>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Initial Adapters
          </label>
          <div className="flex gap-2 flex-wrap">
            {props.adapters.map((a) => (
              <button
                key={a.id}
                type="button"
                onClick={() => toggleAdapter(a.id)}
                className={`px-3 py-1 rounded-full text-sm border ${
                  formData.adapter_ids.includes(a.id)
                    ? "bg-primary-500 text-white border-primary-600"
                    : "bg-gray-100 text-gray-600 border-gray-300"
                }`}
              >
                {a.name}
              </button>
            ))}
          </div>
        </div>

        <button
          type="submit"
          className="mt-4 bg-green-600 text-white py-2 px-4 rounded hover:bg-green-700 transition"
        >
          Register Device
        </button>
      </form>
    </>
  );
};